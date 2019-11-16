package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/gcputil/env"
	"github.com/mchmarny/gcputil/project"
)

var (
	logger    = log.New(os.Stdout, "", 0)
	projectID = project.GetIDOrFail()
)

func main() {
	http.HandleFunc("/", requestHandler)
	port := fmt.Sprintf(":%s", env.MustGetEnvVar("PORT", "8080"))
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal(err)
	}
}
