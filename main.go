package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func main() {
	r := gin.Default()
	r.POST("/chat", handleChat)
	r.Run(":3000")
}

func handleChat(c *gin.Context) {
	var request struct {
		Message string `json:"message"`
	}
	c.BindJSON(&request)

	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	// Use GPT-3.5
	resp, err := client.CreateCompletion(ctx, openai.CompletionRequest{
		Model:  openai.GPT3Dot5Turbo, // GPT-3.5 model
		Prompt: request.Message,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"response": resp.Choices[0].Text})
}
