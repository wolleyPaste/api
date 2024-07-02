package common

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/log"
)

type PostgreSQL struct {
	DBName string

	Host string
	Port string

	Username string
	Password string
}

type Paste struct {
	ID string

	Title string
	Text  string

	Password string

	ExpirationDate time.Time
	CreationDate   time.Time
}

type RabbitMQ struct {
	Queue string

	Host string
	Port string

	Username string
	Password string
}

// This is a custom base64 with the "/" removed because you can't use it in urls
const customBase64 string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+"

func GenerateRandomBase64(idLength int) string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = customBase64[rand.Intn(len(customBase64))]
	}
	return string(b)
}

func StringToTime(strTime string) *time.Time {
	if expTime, err := time.Parse(time.RFC3339, strTime); err == nil {
		return &expTime
	} else {
		log.Errorf("error parsing expiration date: %v", err)
		return nil
	}
}
