package agent

import (
	"fmt"

	"go-ollama-agent/internal/ollama"
)

// Agent represents the code generation agent that interfaces with Ollama
type Agent struct {
	client       *ollama.Client
	model        string
	systemPrompt string
}

// NewAgent creates a new agent with the specified Ollama client and model
func NewAgent(client *ollama.Client, model string) *Agent {
	return &Agent{
		client:       client,
		model:        model,
		systemPrompt: `You are an expert software engineer AI assistant. Write clean, efficient, and well-documented code. Address the user's request directly. Include comments for complex sections. Follow best practices. Use the requested programming language or the most appropriate one. Only respond with code, no explanations or additional text.`,
	}
}

// GenerateCode takes a user prompt and returns generated code from Ollama
func (a *Agent) GenerateCode(userPrompt string) (chan ollama.StreamingCompletionResponse, chan error) {
	prompt := fmt.Sprintf("%s\n\nUser request: %s", a.systemPrompt, userPrompt)

	// Create the completion request
	request := ollama.CompletionRequest{
		Model:  a.model,
		Prompt: prompt,
		Stream: true,
	}

	// Send the streaming request to Ollama
	responseChannel, errChannel := a.client.CreateCompletionStream(request)

	return responseChannel, errChannel
}

func (a *Agent) GenerateCodeWithHistory(userPrompt string, chatHistory []ollama.ConversationItem) (chan ollama.ChatStreamResponse, chan error) {
	prompt := ollama.ConversationItem{
		Role:    "user",
		Content: userPrompt,
	}

	// Create the completion request
	request := ollama.ChatRequest{
		Model:    a.model,
		Messages: append(chatHistory, prompt),
		Stream:   true,
	}

	// Send the streaming request to Ollama
	responseChannel, errChannel := a.client.CreateChatStream(request)

	return responseChannel, errChannel
}
