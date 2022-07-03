package main

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	c "github.com/pbalan/ontrack/src/config"
	"github.com/pbalan/ontrack/src/graph"
	"github.com/pbalan/ontrack/src/graph/generated"
	"github.com/pbalan/ontrack/src/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
)

var Db *gorm.DB
var configuration c.Configurations
var vpr = viper.New()

func initEnv() {
	// Set the file name of the configurations file
	vpr.SetConfigName("config")
	// Set the path to look for the configurations file
	vpr.AddConfigPath(".")
	// Enable VIPER to read Environment Variables
	vpr.AutomaticEnv()
	vpr.SetConfigType("yml")

	if err := vpr.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := vpr.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}
}

func initLogger() {
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

func initDb() *gorm.DB {
	Db = connectDB()
	return Db
}

func connectDB() *gorm.DB {
	var err error
	dsn := configuration.Database.DBUser + ":" + string(configuration.Database.DBPassword) +
		"@tcp" + "(" + configuration.Database.DBHost + ":" +
		strconv.FormatInt(int64(configuration.Database.DBPort), 10) + ")/" +
		configuration.Database.DBName + "?" + "parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Error connecting to database : error=%v", err)
		log.Println(dsn)
		return nil
	}

	return db
}

func autoMigrateSchema(db *gorm.DB) {
	db.AutoMigrate(&models.User{}, &models.TokenDetail{})
}

func startServer(db *gorm.DB) {
	port := os.Getenv("PORT")
	if port == "" {
		port = strconv.FormatInt(int64(configuration.Server.Port), 10)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	initEnv()
	initLogger()
	db := initDb()
	autoMigrateSchema(db)
	startServer(db)
}
