# mcp-server-demo

使用 mcp-go 开发的简单服务器示例。

## 目录结构

```
mcp-server-demo/
├── cmd/                    # 主要的应用程序入口
│   └── server/            # 服务器启动相关代码
├── deploy/                # 部署相关配置和脚本
├── internal/
│   ├── config/           # 配置管理
│   ├── handler/          # HTTP 处理器
│   ├── model/            # 数据模型
│   ├── service/          # 业务逻辑服务
│   └── utils/            # 工具函数
├── pkg/
│   └── log/              # 日志包
└── go.mod                # Go 模块依赖文件
```

## 开始使用

```bash
go mod tidy

go run cmd/server
```

## 许可证

[MIT](https://opensource.org/licenses/MIT) 