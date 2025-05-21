package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"go-ollama-agent/internal/agent"
	"go-ollama-agent/internal/ollama"
)

func main() {
	// Define command-line flags
	ollamaURL := flag.String("url", "http://localhost:11434", "Ollama API URL")
	model := flag.String("model", "gemma3:1b", "Model to use for code generation")
	typingSpeed := flag.Int("speed", 10, "Typing speed in milliseconds per character")
	flag.Parse()

	// Create Ollama client and agent
	client := ollama.NewClient(*ollamaURL)
	codeAgent := agent.NewAgent(client, *model)

	fmt.Println("Go Ollama Code Agent")
	fmt.Println("-------------------")
	fmt.Println("Enter your code generation prompt (type 'exit' to quit):")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		if userInput == "exit" {
			break
		}

		// Skip empty inputs
		if strings.TrimSpace(userInput) == "" {
			continue
		}

		// Generate code based on user input
		responseChannel, errChannel := codeAgent.GenerateCode(userInput)

		// Display the generated code
		fmt.Println("\nGenerated Code:")
		fmt.Println("---------------")

		// For character-by-character printing
		done := false
		for !done {
			select {
			case chunk, ok := <-responseChannel:
				if !ok {
					// Channel closed, streaming complete
					done = true
					break
				}

				// Print the response character by character with a delay
				for _, char := range chunk.Response {
					fmt.Print(string(char))
					// Flush stdout to ensure immediate display
					os.Stdout.Sync()
					// Sleep to create typing effect
					time.Sleep(time.Duration(*typingSpeed) * time.Millisecond)
				}

				// If this is the last chunk, we're done
				if chunk.Done {
					done = true
				}

			case err, ok := <-errChannel:
				if !ok {
					// Error channel closed
					done = true
					break
				}
				fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
				done = true
			}
		}

		fmt.Println("\n---------------")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
