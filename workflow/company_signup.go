package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.temporal.io/sdk/log"

	"github.com/go-resty/resty/v2"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

type CreateCompanyParams struct {
	Properties CreateCompanyProperties `json:"properties"`
}

type CreateCompanyProperties struct {
	Name  string `json:"name"`
	City  string `json:"city"`
	State string `json:"state"`
	Phone string `json:"phone"`
}

type CompanySignup struct {
	token string
}

func NewCompanySignup() *CompanySignup {
	return &CompanySignup{
		token: os.Getenv("HUBSPOT_TOKEN"),
	}
}

func (c *CompanySignup) Workflow(ctx workflow.Context, params CreateCompanyParams) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("CompanySignup workflow started", "name", params.Properties.Name)

	var result string
	err := workflow.ExecuteActivity(ctx, c.Activity, params).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("CompanySignup workflow completed.", "result", result)

	return result, nil
}

func (c *CompanySignup) Activity(ctx context.Context, params CreateCompanyParams) error {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", params.Properties.Name)
	err := c.CreateCompanyToHubSpot(ctx, params, logger)
	if err != nil {
		logger.Info(fmt.Sprintf("Error CreateCompanyToHubSpot: %v", err))
		return err
	}
	return nil
}

func (c *CompanySignup) CreateCompanyToHubSpot(ctx context.Context,
	params CreateCompanyParams,
	logger log.Logger,
) error {
	client := resty.New()

	body, err := json.Marshal(params)
	if err != nil {
		logger.Info(fmt.Sprintf("Error marshalling CreateCompanyParams: %v", err))
		return err
	}

	// Call the HubSpot API to create a company
	resp, err := client.R().
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.token)).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("https://api.hubapi.com/crm/v3/objects/companies")

	if err != nil {
		logger.Info(fmt.Sprintf("Failed to create company: %v", err))
		return err
	}
	logger.Info(fmt.Sprintf("HubSpot response: %s", resp.String()))
	return nil
}
