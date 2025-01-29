# deepseek-go
The Go client for the DeepSeek API

## Usage:
```go
package main

import (
	"fmt"
	"log"

	deepseek "github.com/kashifkhan0771/deepseek-go"
)

func main() {
	// Replace "your-api-token" with your actual DeepSeek API token
	token := "your-api-token"

	// Create a new deepseek client
	client, err := deepseek.NewClient(token)
	if err != nil {
		log.Fatalf("Failed to create DeepSeek config: %v", err)
	}

	// Define the chat completion request payload
	payload := deepseek.ChatCompletionRequest{
		Model: "deepseek-chat", // Replace with the desired model
		Messages: []deepseek.Message{
			{
				Role:    "user",
				Content: "Hello, how can I assist you today?",
			},
		},
		MaxTokens: 50, // Adjust as needed
	}

	// Call the ChatCompletion method
	response, err := client.ChatCompletion(payload)
	if err != nil {
		log.Fatalf("Failed to get chat completion: %v", err)
	}

	// Print the response from the API
	fmt.Println("Chat Completion Response:")
	fmt.Println(response.Choices[0].Message.Content)
}
