package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ============================================================================
// Request/Response Models (Agent Protocol Contract)
// ============================================================================

type AgentInput struct {
	Prompt  string                 `json:"prompt"`
	Context map[string]interface{} `json:"context,omitempty"`
}

type AgentMetadata struct {
	RequestID string `json:"requestId"`
	UserID    string `json:"userId"`
	Timeout   int    `json:"timeout"`
}

type ExecuteRequest struct {
	Input    AgentInput    `json:"input"`
	Metadata AgentMetadata `json:"metadata"`
}

type TokenUsage struct {
	Input  int `json:"input"`
	Output int `json:"output"`
}

type ResultMetadata struct {
	Cost       float64     `json:"cost"`
	Duration   int64       `json:"duration"` // milliseconds
	Model      string      `json:"model,omitempty"`
	TokensUsed *TokenUsage `json:"tokensUsed,omitempty"`
}

type ExecuteResponse struct {
	Result   interface{}    `json:"result"`
	Metadata ResultMetadata `json:"metadata"`
	Error    *string        `json:"error,omitempty"`
}

type InfoResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Version      string              `json:"version"`
	Capabilities []string            `json:"capabilities"`
	Pricing      map[string]interface{} `json:"pricing"`
}

type HealthResponse struct {
	Status    string  `json:"status"`
	Timestamp float64 `json:"timestamp"`
}

type RootResponse struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Status  string `json:"status"`
}

// ============================================================================
// Claude API Client
// ============================================================================

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	System    string          `json:"system,omitempty"`
	Messages  []ClaudeMessage `json:"messages"`
}

type ClaudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ClaudeUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type ClaudeResponse struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Role    string          `json:"role"`
	Content []ClaudeContent `json:"content"`
	Model   string          `json:"model"`
	Usage   ClaudeUsage     `json:"usage"`
}

type ClaudeClient struct {
	apiKey             string
	model              string
	inputCostPerToken  float64
	outputCostPerToken float64
}

func NewClaudeClient() (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	return &ClaudeClient{
		apiKey:             apiKey,
		model:              "claude-sonnet-4-5-20250929",
		inputCostPerToken:  3.0 / 1_000_000,  // $3 per 1M tokens
		outputCostPerToken: 15.0 / 1_000_000, // $15 per 1M tokens
	}, nil
}

func (c *ClaudeClient) CalculateCost(inputTokens, outputTokens int) float64 {
	inputCost := float64(inputTokens) * c.inputCostPerToken
	outputCost := float64(outputTokens) * c.outputCostPerToken
	// Round to 4 decimal places
	return float64(int((inputCost+outputCost)*10000)) / 10000
}

func (c *ClaudeClient) CreateMessage(system, userMessage string) (*ClaudeResponse, error) {
	reqBody := ClaudeRequest{
		Model:     c.model,
		MaxTokens: 4000,
		System:    system,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: userMessage,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(body))
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &claudeResp, nil
}

// ============================================================================
// Agent Implementation
// ============================================================================

type GoAgent struct {
	client *ClaudeClient
}

func NewGoAgent() (*GoAgent, error) {
	client, err := NewClaudeClient()
	if err != nil {
		return nil, err
	}

	return &GoAgent{
		client: client,
	}, nil
}

func (a *GoAgent) Execute(prompt string, context map[string]interface{}) (map[string]interface{}, error) {
	startTime := time.Now()

	// Build system prompt
	systemPrompt := a.buildSystemPrompt()

	// Build user message
	userMessage := a.buildUserMessage(prompt, context)

	// Call Claude API
	response, err := a.client.CreateMessage(systemPrompt, userMessage)
	if err != nil {
		return nil, fmt.Errorf("Claude API call failed: %w", err)
	}

	// Extract result
	var resultText string
	if len(response.Content) > 0 {
		resultText = response.Content[0].Text
	}

	// Parse result (JSON or plain text)
	result := a.parseResult(resultText)

	// Calculate metrics
	duration := time.Since(startTime).Milliseconds()
	cost := a.client.CalculateCost(response.Usage.InputTokens, response.Usage.OutputTokens)

	return map[string]interface{}{
		"result":   result,
		"cost":     cost,
		"duration": duration,
		"model":    a.client.model,
		"tokensUsed": map[string]int{
			"input":  response.Usage.InputTokens,
			"output": response.Usage.OutputTokens,
		},
	}, nil
}

func (a *GoAgent) buildSystemPrompt() string {
	return `You are a helpful AI assistant specialized in data analysis and insights.

Your task is to analyze the provided information and generate actionable insights.

Respond with valid JSON in this format:
{
  "analysis": "Your detailed analysis",
  "insights": ["insight 1", "insight 2", "insight 3"],
  "recommendations": ["recommendation 1", "recommendation 2"]
}`
}

func (a *GoAgent) buildUserMessage(prompt string, context map[string]interface{}) string {
	message := fmt.Sprintf("User Query: %s\n\n", prompt)

	if context != nil && len(context) > 0 {
		message += "Additional Context:\n"
		for key, value := range context {
			message += fmt.Sprintf("- %s: %v\n", key, value)
		}
	}

	return message
}

func (a *GoAgent) parseResult(text string) interface{} {
	// Try to parse as JSON
	var jsonResult map[string]interface{}

	// Remove markdown code blocks if present
	cleanText := text
	if bytes.Contains([]byte(text), []byte("```json")) {
		parts := bytes.Split([]byte(text), []byte("```json"))
		if len(parts) > 1 {
			parts2 := bytes.Split(parts[1], []byte("```"))
			if len(parts2) > 0 {
				cleanText = string(bytes.TrimSpace(parts2[0]))
			}
		}
	} else if bytes.Contains([]byte(text), []byte("```")) {
		parts := bytes.Split([]byte(text), []byte("```"))
		if len(parts) > 2 {
			cleanText = string(bytes.TrimSpace(parts[1]))
		}
	}

	if err := json.Unmarshal([]byte(cleanText), &jsonResult); err == nil {
		return jsonResult
	}

	// Return as plain text if not JSON
	return map[string]string{"text": text}
}

// ============================================================================
// HTTP Handlers
// ============================================================================

var agent *GoAgent

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RootResponse{
		Name:    "Mindra Go Agent",
		Version: "1.0.0",
		Status:  "running",
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:    "healthy",
		Timestamp: float64(time.Now().Unix()),
	})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(InfoResponse{
		ID:          "agent-go-template",
		Name:        "Go Template Agent",
		Description: "Template for building Go-based agents",
		Version:     "1.0.0",
		Capabilities: []string{
			"Data analysis",
			"Insight generation",
			"Custom processing",
		},
		Pricing: map[string]interface{}{
			"estimatedCost": 0.45,
			"currency":      "USD",
		},
	})
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Execute agent
	resultData, err := agent.Execute(req.Input.Prompt, req.Input.Context)
	if err != nil {
		errMsg := err.Error()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ExecuteResponse{
			Result: nil,
			Metadata: ResultMetadata{
				Cost:     0,
				Duration: 0,
			},
			Error: &errMsg,
		})
		return
	}

	// Build response
	tokensUsed := resultData["tokensUsed"].(map[string]int)
	response := ExecuteResponse{
		Result: resultData["result"],
		Metadata: ResultMetadata{
			Cost:     resultData["cost"].(float64),
			Duration: resultData["duration"].(int64),
			Model:    resultData["model"].(string),
			TokensUsed: &TokenUsage{
				Input:  tokensUsed["input"],
				Output: tokensUsed["output"],
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ============================================================================
// Main
// ============================================================================

func main() {
	var err error
	agent, err = NewGoAgent()
	if err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/execute", executeHandler)

	fmt.Printf(`
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    🐹 Mindra Go Agent Starting
    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

    Port:     %s
    Model:    claude-sonnet-4-5-20250929

    Endpoints:
      • GET  /           - Root
      • GET  /health     - Health check
      • GET  /info       - Agent metadata
      • POST /execute    - Execute agent

    Test with:
      curl http://localhost:%s/health

    ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
    `, port, port)

	portNum, _ := strconv.Atoi(port)
	addr := fmt.Sprintf(":%d", portNum)
	log.Printf("Server listening on %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
