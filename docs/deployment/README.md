# ChromLLM Deployment Documentation

This directory contains deployment documentation for the ChromLLM system.

## Backend Deployment

### Docker Deployment

```bash
# Build the Docker image
docker build -t chromllm-backend ./backend

# Run the container
docker run -p 8080:8080 \
  -e OPENAI_API_KEY=your_key \
  -e REDIS_URL=redis://redis:6379 \
  chromllm-backend
```

### Kubernetes Deployment

See `k8s.yaml` for Kubernetes deployment configuration.

### Environment Variables

- `OPENAI_API_KEY` - OpenAI API key for LLM integration
- `REDIS_URL` - Redis connection URL
- `JWT_SECRET` - Secret for JWT token signing
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development/production)

## Extension Deployment

### Chrome Web Store

1. Build the extension: `npm run build`
2. Package the extension: Create ZIP of `dist/` directory
3. Upload to Chrome Web Store Developer Dashboard
4. Submit for review and publish

### Development Deployment

Load the extension unpacked in developer mode for testing.

## Monitoring

- Health check endpoint: `/v1/health`
- Metrics endpoint: `/v1/metrics` (if configured)
- Structured logging with JSON format