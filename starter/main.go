package main

import (
	"context"
	"log"

	wf "github.com/ekypradhanaags/temporal-simple-workflow/workflow"

	"github.com/joho/godotenv"
	"go.temporal.io/sdk/client"
)

var (
	newCompanyName = "New Company from Temporal"
	newCompanyCity = "Blitar"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "company_signup_workflow_id",
		TaskQueue: "company-signup",
	}

	cs := wf.NewCompanySignup()
	properties := wf.CreateCompanyProperties{
		Name: newCompanyName,
		City: newCompanyCity,
	}
	param := wf.CreateCompanyParams{
		Properties: properties,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, cs.Workflow, param)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
