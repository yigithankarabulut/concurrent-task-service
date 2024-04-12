package main

import (
	"log"
	"os"

	"github.com/yigithankarabulut/ConcurrentTaskService/apiserver"
)

// @Title           Task API
// @Version         1.0
// @Description     This is a basic server for managing tasks concurrently. It provides endpoints for creating, updating, deleting, and listing tasks. The server also supports JWT authentication for secure access to the API.
// @Host            localhost:8080
// @BasePath        /task

// @securityDefinitions.apikey BearerAuth
// @In header
// @Name Authorization
// @Description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345". Endpoint created for token generating

// @ExternalDocs.description  OpenAPI
// @ExternalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := apiserver.New(
		apiserver.WithLogLevel(os.Getenv("LOG_LEVEL")),
	); err != nil {
		log.Fatal(err)
	}
}
