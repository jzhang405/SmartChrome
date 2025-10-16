# Quickstart Guide: ChromLLM

**Feature**: Smart Chrome Extension with LLM Integration
**Date**: 2025-10-16
**Branch**: 001-specify-scripts-bash

## Overview

ChromLLM is a Chrome extension that enables intelligent question-answering about webpage content using large language models. This guide walks through the complete setup and usage workflow.

## Prerequisites

- Chrome browser (version 88+)
- Node.js 18+ (for development)
- Go 1.21+ (for backend development)
- Redis 7+ (for caching)
- OpenAI API key

## Quick Setup

### 1. Backend Setup

```bash
# Clone the repository
git clone https://github.com/your-org/chromllm.git
cd chromllm/backend

# Install dependencies
go mod download

# Configure environment
cp .env.example .env
# Edit .env with your OpenAI API key and Redis URL

# Start the backend
go run cmd/server/main.go
```

### 2. Chrome Extension Setup

```bash
cd ../chrome-extension

# Install dependencies
npm install

# Build the extension
npm run build

# Load in Chrome
# 1. Open chrome://extensions/
# 2. Enable "Developer mode"
# 3. Click "Load unpacked"
# 4. Select the chrome-extension/dist directory
```

### 3. Configuration

1. Click the ChromLLM extension icon in Chrome
2. Go to Settings
3. Configure backend URL (default: http://localhost:8080)
4. Save settings

## Usage Workflow

### Basic Usage

1. **Browse to any webpage** with text content
2. **Click the ChromLLM extension icon** in the toolbar
3. **Extract content** by clicking "Extract Page Content"
4. **Ask questions** about the content in the chat interface
5. **Receive real-time responses** from the LLM

### Example Scenarios

#### Scenario 1: Content Summary
```
User: "What is this page about?"
LLM: "This page discusses the implementation of microservices architecture..."
```

#### Scenario 2: Technical Documentation
```
User: "How do I implement this API endpoint?"
LLM: "Based on the documentation, you need to..."
```

#### Scenario 3: Learning Assistance
```
User: "Explain this concept in simpler terms"
LLM: "Think of it like this..."
```

## API Integration Examples

### Create Session
```bash
curl -X POST http://localhost:8080/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{}'
```

### Start Conversation
```bash
curl -X POST http://localhost:8080/v1/conversations \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "url": "https://example.com/article",
    "title": "Example Article",
    "webpage_content": {
      "extracted_text": "Article content here...",
      "content_hash": "abc123"
    }
  }'
```

### Send Message
```bash
curl -X POST http://localhost:8080/v1/conversations/CONVERSATION_ID/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "content": "What is the main point of this article?",
    "type": "user_question"
  }'
```

### WebSocket Streaming
```javascript
const ws = new WebSocket(
  `ws://localhost:8080/v1/stream?conversationId=CONVERSATION_ID&messageId=MESSAGE_ID`
);

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Streaming response:', data.content);
};
```

## Development Workflow

### Backend Development
```bash
# Run tests
go test ./...

# Build for production
go build -o bin/server cmd/server/main.go

# Run with hot reload (using air)
air
```

### Extension Development
```bash
# Development build with watch
npm run dev

# Production build
npm run build

# Run tests
npm test
```

### Testing Integration

1. **Unit Tests**: Verify individual components
2. **Integration Tests**: Test extension-backend communication
3. **Contract Tests**: Validate API specifications
4. **End-to-End Tests**: Complete user workflow

## Deployment

### Backend Deployment

**Docker Deployment**:
```bash
# Build image
docker build -t chromllm-backend .

# Run container
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=your_key \
  -e REDIS_URL=redis://redis:6379 \
  chromllm-backend
```

**Kubernetes Deployment**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chromllm-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: chromllm-backend
  template:
    metadata:
      labels:
        app: chromllm-backend
    spec:
      containers:
      - name: backend
        image: chromllm-backend:latest
        ports:
        - containerPort: 8080
        env:
        - name: OPENAI_API_KEY
          valueFrom:
            secretKeyRef:
              name: chromllm-secrets
              key: openai-api-key
```

### Extension Deployment

1. **Build for production**: `npm run build`
2. **Package extension**: Create ZIP of `dist/` directory
3. **Upload to Chrome Web Store**: Developer dashboard
4. **Publish**: Submit for review and publish

## Monitoring and Debugging

### Backend Monitoring

- **Health checks**: `/health` endpoint
- **Metrics**: Prometheus metrics endpoint
- **Logs**: Structured JSON logging
- **Tracing**: Distributed tracing with OpenTelemetry

### Extension Debugging

1. **Open DevTools**: Right-click extension popup → Inspect
2. **View console logs**: Extension background and content scripts
3. **Network inspection**: Monitor API calls
4. **Storage inspection**: Check localStorage and session data

### Common Issues

**Extension not loading**:
- Check manifest.json syntax
- Verify permissions in Chrome
- Clear extension cache

**API connection issues**:
- Verify backend URL in settings
- Check CORS configuration
- Validate JWT tokens

**Content extraction failures**:
- Check page permissions
- Verify content script injection
- Test on different websites

## Performance Optimization

### Backend Optimization

- **Caching**: Redis for session and response caching
- **Connection pooling**: Database and HTTP client pooling
- **Compression**: Gzip compression for responses
- **CDN**: Static asset delivery via CDN

### Extension Optimization

- **Lazy loading**: Load resources on demand
- **Caching**: localStorage for user preferences
- **Minification**: Production build optimization
- **Tree shaking**: Remove unused code

## Security Considerations

- **HTTPS**: Always use encrypted connections
- **CORS**: Configure proper cross-origin policies
- **Input validation**: Sanitize all user inputs
- **Rate limiting**: Prevent API abuse
- **Token expiration**: Short-lived JWT tokens

## Next Steps

1. **Explore advanced features**: Conversation history, multi-language support
2. **Customize models**: Configure different LLM providers
3. **Integrate with tools**: Export conversations, share results
4. **Scale deployment**: Load balancing, auto-scaling

## Support

- **Documentation**: [docs.chromllm.example.com](https://docs.chromllm.example.com)
- **Issues**: [GitHub Issues](https://github.com/your-org/chromllm/issues)
- **Community**: [Discord/Slack channel]
- **Email**: support@chromllm.example.com