package main

import (
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/soulchildtc/mcp-server-weather/internal/config"
	"github.com/soulchildtc/mcp-server-weather/internal/handler"
	"github.com/soulchildtc/mcp-server-weather/internal/service"
	"github.com/soulchildtc/mcp-server-weather/pkg/log"
)

func main() {
	// 初始化配置
	cfg := config.NewConfig()

	// 初始化服务
	weatherService := service.NewWeatherService(cfg)

	// 初始化处理器
	weatherHandler := handler.NewWeatherHandler(weatherService)

	// 创建 MCP 服务器
	mcpServer := server.NewMCPServer(
		"weather-server",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
	)

	/////////////////////////////////////////////////////////////////////////////////////
	// 声明工具
	weatherTool := mcp.NewTool("get_current_weather",
		mcp.WithDescription("获取指定经纬度的当前天气信息"),
		mcp.WithString("location",
			mcp.Description("位置坐标，格式为 经度,纬度 例如 116.41,39.92"),
		),
	)

	geoTool := mcp.NewTool("get_geo",
		mcp.WithDescription("获取城市名称对应的位置坐标信息"),
		mcp.WithString("city",
			mcp.Description("城市名称, 例如 北京"),
		),
	)

	smartWeatherTool := mcp.NewTool("get_smart_weather",
		mcp.WithDescription("智能查询天气，自动识别城市名称或经纬度坐标"),
		mcp.WithString("query",
			mcp.Description("查询条件，可以是城市名称（如 北京）或经纬度坐标（如 116.41,39.92）"),
		),
	)

	// 注册工具
	mcpServer.AddTool(weatherTool, weatherHandler.GetCurrentWeather)
	mcpServer.AddTool(geoTool, weatherHandler.GetGeoInfo)
	mcpServer.AddTool(smartWeatherTool, weatherHandler.GetSmartWeather)
	/////////////////////////////////////////////////////////////////////////////////////

	// 创建 SSE 服务器

	sseServer := server.NewSSEServer(mcpServer,
		server.WithBaseURL("http://:8080"),
		server.WithBasePath("/weather"),
	)

	// 创建自定义 HTTP 处理器，包含健康检查端点和 SSE 服务器
	customHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理健康检查请求
		if r.URL.Path == "/health/live" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"UP"}`))
			return
		}

		if r.URL.Path == "/health/ready" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"READY"}`))
			return
		}

		// 其他请求交给 SSE 服务器处理
		sseServer.ServeHTTP(w, r)
	})

	// 启动服务器
	log.Info("SSE server listening on :8080")
	if err := http.ListenAndServe(":8080", customHandler); err != nil {
		log.Fatal("服务器错误: %v", err)
	}
}
