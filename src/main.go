package main

import (
	"encoding/json"
	"fmt"
	c "github.com/pbalan/ontrack/src/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func initLogger() {
	// ########## Init Viper
	var vpr = viper.New()
	// Set the file name of the configurations file
	vpr.SetConfigName("config")
	// Set the path to look for the configurations file
	vpr.AddConfigPath(".")
	// Enable VIPER to read Environment Variables
	vpr.AutomaticEnv()
	vpr.SetConfigType("yml")

	var configuration c.Configurations

	if err := vpr.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := vpr.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	if configuration.Server.Debug {
		// Only log the warning severity or above.
		log.SetLevel(log.DebugLevel)

		// Reading variables using the model
		log.Info("Reading variables using the model..")
		database, _ := json.Marshal(configuration.Database)
		server, _ := json.Marshal(configuration.Server)
		log.Info("\nDatabase Config \n", string(database)+"\n", "\nServer Config \n", string(server))
	}
}

func main() {
	initLogger()
}
