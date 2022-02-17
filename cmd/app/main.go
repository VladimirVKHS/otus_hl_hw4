package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"otus_dialog_go/internal/logger"
	"otus_dialog_go/internal/otusdb"
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

	httpPort, _ := os.LookupEnv("HTTP_PORT")
	fmt.Println("HTTP server started, port: " + httpPort)
	http.ListenAndServe(":"+httpPort, nil)
}
