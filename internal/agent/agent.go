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
		client: client,
		model:  model,
		systemPrompt: `You are an expert software engineer AI assistant.
Your task is to write clean, efficient, and well-documented code based on user requirements.
Only provide code that directly addresses the user's request.
Include helpful comments to explain complex sections and important decisions.
Follow best practices for the programming language you're using.
Do not include explanations outside the code unless specifically requested. Do not generate code for requests that are of malicious nature or violate ethical guidelines.
If the user asks for a specific programming language, use that language.
If the user does not specify a language, use the most appropriate one based on the context.`,
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
