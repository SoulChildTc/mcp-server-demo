package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/soulchildtc/mcp-server-weather/internal/config"
	"github.com/soulchildtc/mcp-server-weather/internal/model"
)

type WeatherService struct {
	cfg        *config.Config
	httpClient *http.Client
}

func NewWeatherService(cfg *config.Config) *WeatherService {
	return &WeatherService{
		cfg:        cfg,
		httpClient: &http.Client{},
	}
}

func (s *WeatherService) callQweatherAPI(ctx context.Context, apiPath string, params url.Values) ([]byte, error) {
	params.Set("key", s.cfg.QWeatherAPIKey)
	fullURL := fmt.Sprintf("%s%s?%s", s.cfg.QWeatherBaseURL, apiPath, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建 API 请求失败 (%s): %w", apiPath, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送 API 请求失败 (%s): %w", apiPath, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取 API 响应失败 (%s): %w", apiPath, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 请求失败 (%s)，状态码: %d, 响应: %s", apiPath, resp.StatusCode, string(body))
	}

	return body, nil
}

func (s *WeatherService) GetCurrentWeather(ctx context.Context, location string) (*mcp.CallToolResult, error) {
	params := url.Values{}
	params.Set("location", location)
	body, err := s.callQweatherAPI(ctx, "/v7/weather/now", params)
	if err != nil {
		return nil, err
	}

	var weatherResp model.WeatherResponse
	if err = json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("解析天气响应 JSON 失败: %w, 原始响应: %s", err, string(body))
	}

	if weatherResp.Code != "200" {
		return nil, fmt.Errorf("天气 API 业务错误，代码: %s, 原始响应: %s", weatherResp.Code, string(body))
	}

	now := weatherResp.Now
	resultText := fmt.Sprintf("观测时间：%s\n天气：%s，气温：%s℃ (体感 %s℃)\n风：%s %s级 (%s km/h)\n湿度：%s%%，降水：%s mm，气压：%s hPa\n能见度：%s km",
		now.ObsTime,
		now.Text, now.Temp, now.FeelsLike,
		now.WindDir, now.WindScale, now.WindSpeed,
		now.Humidity, now.Precip, now.Pressure,
		now.Vis)

	return mcp.NewToolResultText(resultText), nil
}

func (s *WeatherService) GetGeoInfo(ctx context.Context, city string) (*mcp.CallToolResult, error) {
	params := url.Values{}
	params.Set("location", city)
	body, err := s.callQweatherAPI(ctx, "/geo/v2/city/lookup", params)
	if err != nil {
		return nil, err
	}

	var geoResp model.GeoResponse
	if err = json.Unmarshal(body, &geoResp); err != nil {
		return nil, fmt.Errorf("解析地理位置响应 JSON 失败: %w, 原始响应: %s", err, string(body))
	}

	if geoResp.Code != "200" {
		return nil, fmt.Errorf("地理位置 API 业务错误，代码: %s, 原始响应: %s", geoResp.Code, string(body))
	}

	if len(geoResp.Location) == 0 {
		return nil, fmt.Errorf("未找到指定城市 '%s' 的地理位置信息", city)
	}

	loc := geoResp.Location[0]
	resultText := fmt.Sprintf("城市：%s (%s, %s), 经度: %s, 纬度: %s, ID: %s",
		loc.Name, loc.Adm2, loc.Adm1, loc.Lon, loc.Lat, loc.ID)

	return mcp.NewToolResultText(resultText), nil
}

func (s *WeatherService) GetSmartWeather(ctx context.Context, query string) (*mcp.CallToolResult, error) {
	var locationCoords string
	var err error

	if strings.Contains(query, ",") {
		locationCoords = query
	} else {
		params := url.Values{}
		params.Set("location", query)
		geoBody, err := s.callQweatherAPI(ctx, "/geo/v2/city/lookup", params)
		if err != nil {
			return nil, fmt.Errorf("智能查询：获取 '%s' 的地理位置失败: %w", query, err)
		}

		var geoResp model.GeoResponse
		if err = json.Unmarshal(geoBody, &geoResp); err != nil {
			return nil, fmt.Errorf("智能查询：解析 '%s' 的地理位置响应 JSON 失败: %w, 原始响应: %s", query, err, string(geoBody))
		}

		if geoResp.Code != "200" {
			return nil, fmt.Errorf("智能查询：地理位置 API 业务错误 (城市: %s)，代码: %s, 原始响应: %s", query, geoResp.Code, string(geoBody))
		}

		if len(geoResp.Location) == 0 {
			return nil, fmt.Errorf("智能查询：未找到城市 '%s' 的地理位置信息", query)
		}

		loc := geoResp.Location[0]
		locationCoords = fmt.Sprintf("%s,%s", loc.Lon, loc.Lat)
	}

	weatherParams := url.Values{}
	weatherParams.Set("location", locationCoords)
	body, err := s.callQweatherAPI(ctx, "/v7/weather/now", weatherParams)
	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, err
		}
		return nil, fmt.Errorf("智能查询：获取坐标 '%s' 的天气失败: %w", locationCoords, err)
	}

	var weatherResp model.WeatherResponse
	if err = json.Unmarshal(body, &weatherResp); err != nil {
		return nil, fmt.Errorf("智能查询：解析天气响应 JSON 失败: %w, 原始响应: %s", err, string(body))
	}

	if weatherResp.Code != "200" {
		return nil, fmt.Errorf("智能查询：天气 API 业务错误 (坐标: %s)，代码: %s, 原始响应: %s", locationCoords, weatherResp.Code, string(body))
	}

	now := weatherResp.Now
	resultText := fmt.Sprintf("查询地点：%s\n观测时间：%s\n天气：%s，气温：%s℃ (体感 %s℃)\n风：%s %s级 (%s km/h)\n湿度：%s%%，降水：%s mm，气压：%s hPa\n能见度：%s km",
		query,
		now.ObsTime,
		now.Text, now.Temp, now.FeelsLike,
		now.WindDir, now.WindScale, now.WindSpeed,
		now.Humidity, now.Precip, now.Pressure,
		now.Vis)

	return mcp.NewToolResultText(resultText), nil
}
