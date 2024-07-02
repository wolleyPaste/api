package main

import (
	"fmt"
	"wolley-api/src/common"
	"wolley-api/src/db"
	"wolley-api/src/queue"
	"wolley-api/src/views"

	"github.com/charmbracelet/log"

	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
)

const iniFileName string = "settings.ini"

func loadIni() map[string]string {
	inidata, err := ini.Load(iniFileName)
	if err != nil {
		log.Fatalf("failed to load ini file: %v", err)
	}

	settings := make(map[string]string)

	// Load API parameters
	section := inidata.Section("api")
	settings["apiHost"] = section.Key("host").String()
	settings["apiPort"] = section.Key("port").String()

	// Load RabbitMQ parameters
	section = inidata.Section("rabbitmq")
	settings["queue"] = section.Key("queue").String()

	settings["rabbitmqUsername"] = section.Key("username").String()
	settings["rabbitmqPassword"] = section.Key("password").String()

	settings["rabbitmqHost"] = section.Key("host").String()
	settings["rabbitmqPort"] = section.Key("port").String()

	// Load PostgreSQL parameters
	section = inidata.Section("postgresql")
	settings["psqlDBName"] = section.Key("db").String()

	settings["psqlUsername"] = section.Key("username").String()
	settings["psqlPassword"] = section.Key("password").String()

	settings["psqlHost"] = section.Key("host").String()
	settings["psqlPort"] = section.Key("port").String()

	return settings
}

// Load ini file
var settings = loadIni()

func main() {
	rabbitmq := common.RabbitMQ{
		Queue: settings["queue"],

		Host: settings["rabbitmqHost"],
		Port: settings["rabbitmqPort"],

		Username: settings["rabbitmqUsername"],
		Password: settings["rabbitmqPassword"],
	}
	queue.InitRabbitMQ(rabbitmq)
	defer queue.CleanupRabbitMQ()

	// Connect to PostgreSQL
	postgresql := common.PostgreSQL{
		DBName: settings["psqlDBName"],

		Host: settings["psqlHost"],
		Port: settings["psqlPort"],

		Username: settings["psqlUsername"],
		Password: settings["psqlPassword"],
	}
	db.InitPostgres(postgresql)
	defer db.ClosePostgres()

	r := gin.Default()
	r.POST("/add", views.Add)
	r.GET("/get", views.GetPaste)

	r.Run(fmt.Sprintf("%s:%s", settings["apiHost"], settings["apiPort"])) // listen and serve on 0.0.0.0:8080
}
