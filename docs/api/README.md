# ChromLLM API Documentation

This directory contains the API documentation for the ChromLLM backend service.

## API Reference

See the OpenAPI specification in `/specs/001-specify-scripts-bash/contracts/api-spec.yaml` for detailed API documentation.

## Endpoints

### Authentication
- `POST /v1/sessions` - Create new session
- `GET /v1/sessions/{sessionId}` - Get session details
- `DELETE /v1/sessions/{sessionId}` - Delete session

### Conversations
- `POST /v1/conversations` - Create new conversation
- `GET /v1/conversations/{conversationId}` - Get conversation details
- `GET /v1/conversations/{conversationId}/messages` - Get conversation messages
- `POST /v1/conversations/{conversationId}/messages` - Send message

### Streaming
- `GET /v1/stream` - WebSocket endpoint for real-time LLM streaming

### Health
- `GET /v1/health` - Health check endpoint