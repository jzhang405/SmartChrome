# ChromLLM - 智能Chrome插件与后端服务套件

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![Chrome Extension](https://img.shields.io/badge/chrome-extension-yellow.svg)](https://developer.chrome.com/docs/extensions/)

ChromLLM 是一套整合了前端Chrome插件与后端服务的智能交互系统，旨在通过大语言模型（LLM）能力为用户提供基于网页内容的实时智能问答与个性化服务。

## 项目概述

ChromLLM 允许用户在浏览任意网页时，通过Chrome插件快速提取页面内容并发起智能问答。后端服务负责处理请求、调用大语言模型、管理会话，形成完整的"前端交互-内容提取-后端处理-实时反馈"闭环。

## 技术架构

### Chrome插件（Manifest V3）
- **Popup界面**：用户交互入口，提供问答输入框和历史记录展示
- **Content Script**：注入当前页面，提取网页文本内容
- **Background Service**：管理WebSocket长连接和消息转发
- **通信层**：REST API + WebSocket混合通信模式

### Go后端服务
- **API网关**：基于Gin框架实现接口路由管理
- **大模型集成**：支持多个LLM提供商（OpenAI、DeepSeek、Douban等）
- **缓存层**：Redis缓存高频访问数据
- **用户会话管理**：JWT实现无状态用户认证
- **容器化部署**：Docker支持一键部署

## 功能特性

- [x] 网页内容智能提取
- [x] 多LLM提供商支持（OpenAI、DeepSeek、Douban）
- [x] 实时流式响应
- [x] 会话管理和历史记录
- [x] 用户认证和权限管理
- [x] WebSocket实时通信
- [ ] 多语言支持（计划中）
- [ ] 个性化设置（计划中）

## 快速开始

### 后端服务

```bash
# 克隆项目
git clone https://github.com/jzhang405/SmartChrome.git
cd SmartChrome/backend

# 安装依赖
go mod tidy

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，设置您的API密钥和其他配置

# 启动服务
go run cmd/server/main.go
```

### Chrome扩展

```bash
# 进入扩展目录
cd ../chrome-extension

# 安装依赖
npm install

# 构建扩展
npm run build

# 在Chrome中加载扩展
# 1. 打开 chrome://extensions/
# 2. 启用"开发者模式"
# 3. 点击"加载已解压的扩展程序"
# 4. 选择 chrome-extension/dist 目录
```

## 配置说明

### 环境变量

后端服务支持以下环境变量配置：

```bash
# 服务器配置
PORT=8080
HOST=localhost

# Redis配置
REDIS_URL=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT配置
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24

# OpenAI配置（默认提供商）
OPENAI_API_KEY=your-openai-api-key
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_MAX_TOKENS=1000
OPENAI_TEMPERATURE=0.7

# DeepSeek配置（可选）
DEEPSEEK_API_KEY=your-deepseek-api-key
DEEPSEEK_BASE_URL=https://api.deepseek.com/v1
DEEPSEEK_MODEL=deepseek-chat
DEEPSEEK_MAX_TOKENS=1000
DEEPSEEK_TEMPERATURE=0.7

# Douban配置（可选）
DOUBAN_API_KEY=your-douban-api-key
DOUBAN_BASE_URL=https://api.douban.com/v1
DOUBAN_MODEL=douban-chat
DOUBAN_MAX_TOKENS=1000
DOUBAN_TEMPERATURE=0.7
```

## API文档

详细的API文档请参考 [API文档](docs/api/README.md)。

## 部署

### Docker部署

```bash
# 构建Docker镜像
docker build -t chromllm-backend ./backend

# 运行容器
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=your_key \
  -e REDIS_URL=redis://redis:6379 \
  chromllm-backend
```

### 扩展部署

1. 构建扩展：`npm run build`
2. 打包扩展：创建 `dist/` 目录的ZIP文件
3. 上传到Chrome Web Store开发者仪表板
4. 提交审核并发布

## 开发指南

### 后端开发

```bash
# 运行测试
cd backend && go test ./...

# 构建生产版本
go build -o bin/server cmd/server/main.go

# 使用air进行热重载开发
air
```

### 扩展开发

```bash
# 开发构建（带监听）
cd chrome-extension && npm run dev

# 生产构建
npm run build

# 运行测试
npm test
```

## 项目结构

```
SmartChrome/
├── backend/              # Go后端服务
│   ├── cmd/server/       # 服务入口
│   ├── internal/         # 内部包
│   ├── pkg/              # 公共包
│   ├── config/           # 配置管理
│   └── tests/            # 测试文件
├── chrome-extension/     # Chrome扩展
│   ├── popup/            # 弹出界面
│   ├── content/          # 内容脚本
│   ├── background/       # 后台服务
│   └── options/          # 选项页面
├── design/               # 设计文档
├── docs/                 # 项目文档
├── specs/                # 规范文档
└── tests/                # 测试目录
```

## 贡献

欢迎提交Issue和Pull Request来改进项目！

1. Fork项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

## 许可证

本项目采用MIT许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 联系方式

项目链接: [https://github.com/jzhang405/SmartChrome](https://github.com/jzhang405/SmartChrome)

## 支持

- 文档: [docs.chromllm.example.com](https://docs.chromllm.example.com)
- Issues: [GitHub Issues](https://github.com/jzhang405/SmartChrome/issues)