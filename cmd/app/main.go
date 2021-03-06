package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"otus_dialog_go/internal/logger"
	"otus_dialog_go/internal/otusdb"
	"otus_dialog_go/internal/otusprometheus"
	"otus_dialog_go/internal/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Config not found: .env!")
	}
	logger.Init()
	otusdb.InitDb()
	defer otusdb.CloseDb()

	// Chi routes
	http.Handle("/", routes.RegisterRouter())
	http.Handle("/metrics", otusprometheus.Handler())

	httpPort, _ := os.LookupEnv("HTTP_PORT")
	fmt.Println("HTTP server started, port: " + httpPort)
	err := http.ListenAndServe(":"+httpPort, nil)
	if err != nil {
		logger.Error(err.Error())
	}
}
