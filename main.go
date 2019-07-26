package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	env "github.com/mchmarny/gcputil/env"
	prj "github.com/mchmarny/gcputil/project"
)

var (
	logger  = log.New(os.Stdout, "[EPS] ", 0)
	project = prj.GetIDOrFail()
	topic   = env.MustGetEnvVar("TOPIC", "cloud-builds")
	que     *queue
)

func main() {

	q, err := newQueue(context.Background(), project, topic)
	if err != nil {
		logger.Fatalf("Error creating pubsub client: %v", err)
	}
	que = q

	http.HandleFunc("/", requestHandler)
	port := fmt.Sprintf(":%s", env.MustGetEnvVar("PORT", "8080"))
	if err := http.ListenAndServe(port, nil); err != nil {
		logger.Fatal(err)
	}

}
