package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"forkbomb.eu/yammya"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, yammya.TaskQueue, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(yammya.Validation)
	w.RegisterActivity(yammya.ValidateYaml)
	w.RegisterActivity(yammya.SendEmail)
	w.RegisterActivity(yammya.MakeAPICall)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}