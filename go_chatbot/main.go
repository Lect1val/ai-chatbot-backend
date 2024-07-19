package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"cloud.google.com/go/dialogflow/apiv2/dialogflowpb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("Starting the server...")
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	r := gin.Default()

	// Apply CORS middleware to the router
	r.Use(CORSMiddleware())

	r.POST("/dialogflow/session/", dialogflowSessionHandler)

	r.Run(":8000")
}

func dialogflowSessionHandler(c *gin.Context) {
	var requestData struct {
		SessionID string `json:"session_id"`
		Text      string `json:"text"`
	}
	if err := c.BindJSON(&requestData); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	projectID := os.Getenv("DIALOGFLOW_PROJECT_ID")
	sessionID := requestData.SessionID
	if sessionID == "" {
		sessionID = "default_session"
	}
	text := requestData.Text
	if text == "" {
		text = "Hello"
	}

	sessionPath := "projects/" + projectID + "/agent/sessions/" + sessionID
	ctx := c.Request.Context()

	sessionClient, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsJSON([]byte(os.Getenv("GOOGLE_CREDENTIALS_JSON"))))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error2": err.Error()})
		return
	}
	defer sessionClient.Close()

	textInput := &dialogflowpb.TextInput{Text: text, LanguageCode: "en-US"}
	queryInput := &dialogflowpb.QueryInput{Input: &dialogflowpb.QueryInput_Text{Text: textInput}}
	request := &dialogflowpb.DetectIntentRequest{
		Session:    sessionPath,
		QueryInput: queryInput,
	}

	response, err := sessionClient.DetectIntent(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
		return
	}

	queryResult := response.GetQueryResult()
	intentName := queryResult.GetIntent().GetDisplayName()

	switch intentName {
	case "session search":
		sessionName := queryResult.Parameters.Fields["session_name"].GetStringValue()
		if sessionName != "" {
			c.JSON(http.StatusOK, gin.H{"fulfillmentText": "Information for session: " + sessionName})
		} else {
			c.JSON(http.StatusOK, gin.H{"fulfillmentText": "Session name not provided."})
		}
		return
	case "Aster arcade URL":
		c.JSON(http.StatusOK, gin.H{
			"fulfillmentText": "This is the URL of Aster Arcade: [Aster-arcade](https://aster.arisetech.dev/aster-arcade/)",
		})
		return
	case "Default Fallback Intent", "":
		client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		resp, err := client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model: "gpt-3.5-turbo",
				Messages: []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleUser, Content: "Please provide general information or engage in a casual conversation about: " + text},
				},
				MaxTokens:   150,
				Temperature: 0.7,
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error4": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"fulfillmentText": resp.Choices[0].Message.Content})
		return
	default:
		c.JSON(http.StatusOK, gin.H{"fulfillmentText": queryResult.GetFulfillmentText()})
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		// If it's an OPTIONS request, we should return HTTP 200 OK
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		} else {
			c.Next()
		}
	}
}
