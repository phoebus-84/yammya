package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	app "forkbomb.eu/yammya"
)

func main() {
	// Create the client object just once per process
	c, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer c.Close()

	input := app.AppInput{
		Email: "ennio.fex@gmail.com",
		Url:  "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/ValidYAML.yaml",
	}

	options := client.StartWorkflowOptions{
		ID:        "validate-id",
		TaskQueue: app.TaskQueue,
	}

	log.Printf("Starting Workflow with ID: %s\n, input: %v", options.ID, input)

	we, err := c.ExecuteWorkflow(context.Background(), options, app.Validation, input)
	if err != nil {
		log.Fatalln("Unable to start the Workflow:", err)
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result string

	err = we.Get(context.Background(), &result)

	if err != nil {
		log.Fatalln("Unable to get Workflow result:", err)
	}

	log.Println(result)
}

