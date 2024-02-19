package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gossie/cash-book/cashbook"
)

const (
	DEFAULT_PORT      = "8080"
	DEFAULT_BASE_PATH = "/"
)

//go:embed assets/*
var assets embed.FS

//go:embed templates/*
var htmlTemplates embed.FS

func main() {
	mongoClient := cashbook.ConnectToDatabase(context.Background())
	defer cashbook.DisconnectFromDatabase(context.Background(), mongoClient)

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = DEFAULT_BASE_PATH
	}

	server := cashbook.NewServer(mongoClient, assets, htmlTemplates, basePath)

	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	log.Default().Println("start server on port", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), server))
}
