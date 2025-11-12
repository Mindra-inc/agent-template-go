# Contributing to Go Agent Template

Thank you for your interest in contributing to the Mindra Platform Go Agent Template! 🎉

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Pull Request Process](#pull-request-process)
- [Style Guidelines](#style-guidelines)
- [Commit Messages](#commit-messages)

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code:

- Be respectful and inclusive
- Be patient and welcoming
- Be considerate and constructive
- Focus on what is best for the community
- Show empathy towards others

## How Can I Contribute?

### 🐛 Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title** - Descriptive summary of the issue
- **Steps to reproduce** - Detailed steps to reproduce the problem
- **Expected behavior** - What you expected to happen
- **Actual behavior** - What actually happened
- **Environment** - Go version, OS, architecture
- **Code samples** - Minimal reproduction code if possible

### 💡 Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. Include:

- **Clear description** - Explain the enhancement
- **Use case** - Why would this be useful?
- **Examples** - Show how it would work
- **Alternatives** - Other solutions you've considered

### 🔧 Pull Requests

1. **Fork the repository** and create your branch from `develop`
2. **Make your changes** with clear, focused commits
3. **Add tests** if applicable (optional for simple changes)
4. **Update documentation** as needed
5. **Ensure code compiles** - Run `go build`
6. **Format code** - Run `gofmt -s -w .`
7. **Submit pull request** with clear description

## Development Setup

### Prerequisites

- Go >= 1.21
- Optional: Docker for containerization

### Setup Steps

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/agent-template-go.git
cd agent-template-go

# Set environment variables
export ANTHROPIC_API_KEY="your-api-key"
export PORT=8002

# Run agent
go run main.go

# Or build and run
go build -o agent main.go
./agent
```

### Testing

```bash
# Test health endpoint
curl http://localhost:8002/health

# Test info endpoint
curl http://localhost:8002/info

# Test execute endpoint
curl -X POST http://localhost:8002/execute \
  -H "Content-Type: application/json" \
  -d '{"input": {"prompt": "test"}, "metadata": {"requestId": "test"}}'
```

### Cross-Compilation

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o agent-linux main.go

# Build for macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o agent-macos-intel main.go

# Build for macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o agent-macos-arm main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o agent.exe main.go
```

## Pull Request Process

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation changes
- `refactor/description` - Code refactoring
- `chore/description` - Maintenance tasks

### PR Guidelines

1. **Update CHANGELOG.md** - Add entry under `[Unreleased]`
2. **Update README.md** - If adding/changing features
3. **Follow style guidelines** - Use `gofmt` and `go vet`
4. **Write descriptive commits** - Follow commit message conventions
5. **Reference issues** - Link related issues in PR description
6. **Request review** - Wait for maintainer review

## Style Guidelines

### Go

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt -s -w .` for formatting
- Use `go vet` for static analysis
- Use meaningful variable names
- Comment exported functions
- Keep functions small and focused

### Code Style

```go
package main

import (
    "encoding/json"
    "net/http"
)

// AgentInput represents the input for agent execution
type AgentInput struct {
    Prompt  string                 `json:"prompt"`
    Context map[string]interface{} `json:"context,omitempty"`
}

// ExecuteAgent processes the input and returns a result
func ExecuteAgent(input AgentInput) (map[string]interface{}, error) {
    // Validate input
    if input.Prompt == "" {
        return nil, fmt.Errorf("prompt is required")
    }

    // Process
    result := map[string]interface{}{
        "result": input.Prompt,
    }

    return result, nil
}
```

### Formatting

```bash
# Format code
gofmt -s -w .

# Run static analysis
go vet ./...

# Run tests (if any)
go test ./...
```

## Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `style:` - Code style (formatting)
- `refactor:` - Code change without fixing bug or adding feature
- `perf:` - Performance improvement
- `test:` - Adding tests
- `chore:` - Maintenance tasks

### Examples

```bash
feat(api): add support for streaming responses
fix(auth): handle API key validation correctly
docs(readme): update Go version requirements
refactor(client): simplify HTTP request handling
perf(json): optimize JSON parsing performance
```

## Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR** (1.0.0) - Breaking changes
- **MINOR** (0.1.0) - New features (backwards compatible)
- **PATCH** (0.0.1) - Bug fixes (backwards compatible)

## Testing Guidelines

```go
// Example test structure
func TestExecuteAgent(t *testing.T) {
    input := AgentInput{
        Prompt: "test prompt",
    }

    result, err := ExecuteAgent(input)
    if err != nil {
        t.Fatalf("ExecuteAgent failed: %v", err)
    }

    if result == nil {
        t.Error("Expected non-nil result")
    }
}
```

## Questions?

- 💬 Open a [Discussion](https://github.com/Mindra-inc/agent-template-go/discussions)
- 🐛 Report a [Bug](https://github.com/Mindra-inc/agent-template-go/issues/new?template=bug_report.md)
- 💡 Request a [Feature](https://github.com/Mindra-inc/agent-template-go/issues/new?template=feature_request.md)
- 📧 Email: support@mindra.co

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- Release notes
- Project README

Thank you for contributing! 🙏
