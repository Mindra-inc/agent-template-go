# Mindra Go Agent Template

A Go-based agent template for the Mindra Platform that implements the standard HTTP Agent Protocol.

## Features

- ✅ **Standard HTTP Protocol** - Compatible with TypeScript orchestrator
- ✅ **Claude Sonnet 4.5 Integration** - Latest Claude AI model
- ✅ **Zero External Dependencies** - Uses only Go standard library
- ✅ **Fast & Efficient** - Compiled binary with minimal overhead
- ✅ **Cost Tracking** - Automatic token usage and cost calculation
- ✅ **JSON Support** - Automatic JSON parsing and validation
- ✅ **Health Checks** - Built-in health and info endpoints

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Anthropic API Key

### Installation

```bash
# Navigate to template directory
cd services/agents/go-template

# Set your API key
export ANTHROPIC_API_KEY="your-api-key-here"

# Run the agent
go run main.go
```

### Using with Custom Port

```bash
PORT=8003 go run main.go
```

### Building Binary

```bash
# Build for current platform
go build -o agent main.go

# Run the compiled binary
./agent
```

### Cross-compilation Examples

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o agent-linux main.go

# Build for macOS
GOOS=darwin GOARCH=arm64 go build -o agent-macos main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o agent.exe main.go
```

## API Endpoints

### GET /
Root endpoint with basic information.

**Response:**
```json
{
  "name": "Mindra Go Agent",
  "version": "1.0.0",
  "status": "running"
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": 1704654321.123
}
```

### GET /info
Agent metadata and capabilities.

**Response:**
```json
{
  "id": "agent-go-template",
  "name": "Go Template Agent",
  "description": "Template for building Go-based agents",
  "version": "1.0.0",
  "capabilities": [
    "Data analysis",
    "Insight generation",
    "Custom processing"
  ],
  "pricing": {
    "estimatedCost": 0.45,
    "currency": "USD"
  }
}
```

### POST /execute
Execute agent with given input.

**Request:**
```json
{
  "input": {
    "prompt": "Analyze this data and provide insights",
    "context": {
      "data": "sample data",
      "options": {}
    }
  },
  "metadata": {
    "requestId": "req_123",
    "userId": "user_456",
    "timeout": 30000
  }
}
```

**Response:**
```json
{
  "result": {
    "analysis": "Detailed analysis...",
    "insights": ["insight 1", "insight 2"],
    "recommendations": ["rec 1", "rec 2"]
  },
  "metadata": {
    "cost": 0.0234,
    "duration": 1523,
    "model": "claude-sonnet-4-5-20250929",
    "tokensUsed": {
      "input": 150,
      "output": 300
    }
  }
}
```

## Testing

### Test Health Endpoint
```bash
curl http://localhost:8002/health
```

### Test Execute Endpoint
```bash
curl -X POST http://localhost:8002/execute \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "prompt": "What are the benefits of Go for building microservices?"
    },
    "metadata": {
      "requestId": "test_001",
      "userId": "test_user",
      "timeout": 30000
    }
  }'
```

## Customization

### Modify Agent Behavior

Edit the `buildSystemPrompt()` method in `main.go`:

```go
func (a *GoAgent) buildSystemPrompt() string {
    return `Your custom system prompt here...`
}
```

### Add Custom Processing

Extend the `Execute()` method:

```go
func (a *GoAgent) Execute(prompt string, context map[string]interface{}) (map[string]interface{}, error) {
    // Add preprocessing
    preprocessed := a.preprocess(prompt, context)

    // Call Claude
    response, err := a.client.CreateMessage(systemPrompt, preprocessed)

    // Add postprocessing
    result := a.postprocess(response)

    return result, nil
}
```

### Change AI Model

Update the model in `NewClaudeClient()`:

```go
return &ClaudeClient{
    apiKey: apiKey,
    model:  "claude-opus-4-20250514", // or any other model
    // ...
}
```

## Environment Variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `ANTHROPIC_API_KEY` | Yes | - | Your Anthropic API key |
| `PORT` | No | `8002` | Port to run the server on |

## Deployment

### Using systemd (Linux)

Create `/etc/systemd/system/mindra-go-agent.service`:

```ini
[Unit]
Description=Mindra Go Agent
After=network.target

[Service]
Type=simple
User=mindra
WorkingDirectory=/opt/mindra/agents/go-template
Environment=ANTHROPIC_API_KEY=your-key-here
Environment=PORT=8002
ExecStart=/opt/mindra/agents/go-template/agent
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable mindra-go-agent
sudo systemctl start mindra-go-agent
```

### Using Docker

Create `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod ./
COPY main.go ./
RUN go build -o agent main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/agent .
ENV PORT=8002
EXPOSE 8002
CMD ["./agent"]
```

Build and run:
```bash
docker build -t mindra-go-agent .
docker run -p 8002:8002 -e ANTHROPIC_API_KEY=your-key mindra-go-agent
```

## Performance

Go agents offer excellent performance characteristics:

- **Startup Time**: < 10ms (compiled binary)
- **Memory Usage**: ~10-15MB base memory
- **Concurrent Requests**: Handles thousands of concurrent connections
- **Response Time**: ~1-3s for typical Claude API calls
- **CPU Efficient**: Minimal overhead compared to interpreted languages

## Troubleshooting

### Agent won't start

Check that `ANTHROPIC_API_KEY` is set:
```bash
echo $ANTHROPIC_API_KEY
```

### API errors

Enable detailed logging:
```bash
# Add to main()
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

### Port already in use

Change the port:
```bash
PORT=8003 go run main.go
```

## License

Part of the Mindra Platform.
