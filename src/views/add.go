package views

import (
	"context"
	"encoding/json"
	"sync"
	"time"
	"wolley-api/src/common"
	"wolley-api/src/db"
	"wolley-api/src/queue"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func CreatePasteID() string {
	randomBase64 := common.GenerateRandomBase64(8)
	if db.CheckPasteExistence(randomBase64) {
		return CreatePasteID()
	}

	return randomBase64
}

func fillOptionalField(c *gin.Context, pasteObject *map[string]interface{}, fieldName string, wg *sync.WaitGroup) {
	defer wg.Done()
	if optionalField := c.Request.Header.Get(fieldName); optionalField != "" {
		(*pasteObject)[fieldName] = optionalField
	}
}

func Add(c *gin.Context) {
	paste := map[string]interface{}{
		"id":   CreatePasteID(),
		"text": c.Request.Header.Get("text"),
	}

	// Fill optional fields (if exists)
	var optionalFields = []string{"title", "password", "expiration"}
	var wg sync.WaitGroup
	for _, optionalField := range optionalFields {
		go fillOptionalField(c, &paste, optionalField, &wg)
	}

	// Dump the paste into a json
	body, err := json.Marshal(paste)
	if err != nil {
		log.Errorf("Failed to marshal JSON: %v", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	// Publish the message
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = queue.RabbitChan.PublishWithContext(ctx,
		"",       // exchange
		"wolley", // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		log.Errorf("failed to publish a message: %v", err)
		c.JSON(500, gin.H{"error": "Failed to publish message"})
		return
	}

	c.JSON(201, gin.H{
		"message": "sucessfully created paste",
		"id":      paste["id"],
	})
}
