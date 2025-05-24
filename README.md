# go-ollama-agent

## Overview
Exercise in learning how to create a coding agent. Currently supports chat with history in the terminal.

## Getting Started

### Prerequisites
- Go 1.24 or later
- Access to the Ollama service
- or Docker to use the devcontainer

### Installation
1. Clone the repository:
   ```
   git clone https://github.com/alpacamybags118/go-ollama-agent.git
   cd go-ollama-agent
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Running the Service
Start the ollama container
```
make ollama-up
```

Build the service
```
make build
```

Then run the binary
```
./build/agent
```

Supports the following flags:

- `url` - URL to your ollama instance
- `model` - What model to run against (ensure you have pulled the model in Ollama)
- `speed` - Simulated typing "speed" of the streamed response

## Future Features / Ideas

- code syntax highlighting
- "Code Review" agent option

The end state I'd like to get to is a system where you can queue coding "jobs", where solution files are generated and reviewed, then provided back to the user in some form.