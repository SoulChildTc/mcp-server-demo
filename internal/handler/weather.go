package handler

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/soulchildtc/mcp-server-weather/internal/service"
)

type WeatherHandler struct {
	weatherService *service.WeatherService
}

func NewWeatherHandler(weatherService *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: weatherService,
	}
}

func (h *WeatherHandler) GetCurrentWeather(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	location, ok := request.Params.Arguments["location"].(string)
	if !ok || location == "" {
		return nil, fmt.Errorf("location 参数不能为空")
	}
	return h.weatherService.GetCurrentWeather(ctx, location)
}

func (h *WeatherHandler) GetGeoInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	city, ok := request.Params.Arguments["city"].(string)
	if !ok || city == "" {
		return nil, fmt.Errorf("city 参数不能为空")
	}
	return h.weatherService.GetGeoInfo(ctx, city)
}

func (h *WeatherHandler) GetSmartWeather(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, ok := request.Params.Arguments["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query 参数不能为空")
	}
	return h.weatherService.GetSmartWeather(ctx, query)
}
