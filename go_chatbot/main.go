package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/webhook", func(c *gin.Context) {
		handleWebhook(c)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func handleWebhook(c *gin.Context) {
	var req struct {
		Message string `json:"message"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	textResponse, err := sendMessageToDialogflow(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process the message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": textResponse})
}

func sendMessageToDialogflow(message string) (string, error) {
	fmt.Println(0)
	ctx := context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	fmt.Println(1)

	projectID := os.Getenv("DIALOGFLOW_PROJECT_ID")
	if projectID == "" {
		log.Fatal("Environment variable DIALOGFLOW_PROJECT_ID not set")
	}

	sessionID := os.Getenv("DIALOGFLOW_SESSION_ID")
	if sessionID == "" {
		log.Fatal("Environment variable DIALOGFLOW_SESSION_ID not set")
	}

	fmt.Println(2)

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: message, LanguageCode: "en-US"}
	queryInput := dialogflowpb.QueryInput{Input: &dialogflowpb.QueryInput_Text{Text: &textInput}}

	fmt.Println(3)

	response, err := sessionClient.DetectIntent(ctx, &dialogflowpb.DetectIntentRequest{
		Session:    sessionPath,
		QueryInput: &queryInput,
	})
	if err != nil {
		return "", err
	}

	return response.GetQueryResult().GetFulfillmentText(), nil
}
