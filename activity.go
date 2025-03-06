package yammya

import (
	"bytes"
	"context"
	"errors"

	"github.com/phoebus-84/kaa/cmd"
	"github.com/pluja/pocketbase"
	"go.temporal.io/sdk/activity"
	gomail "gopkg.in/mail.v2"
)

const (
	schemaPath = "schemas/schema.yaml"
)

func ValidateYaml(ctx context.Context, input ValidationInput) (ValidationOutput, error) {
	logger := activity.GetLogger(ctx)
	schema, err := cmd.LoadYAMLSchema(schemaPath)
	falsyOutput := ValidationOutput{IsValid: false, Message: ""}
	if err != nil {
		logger.Error("Failed to load schema", "Error", err)
		return falsyOutput, err
	}
	f, err := cmd.LoadYAMLFromURL(input.Url)
	if err != nil {
		logger.Error("Failed to load YAML file", "Error", err)
		return falsyOutput, err
	}
	buffer := bytes.Buffer{}
	errValidation := cmd.ValidateYAML(f, schema, &buffer)
	logger.Debug("Validation result", "Error", errValidation)
	if errValidation != nil {
		return ValidationOutput{IsValid: false, Message: buffer.String()}, nil
	}
	return ValidationOutput{IsValid: true, Message: ""}, nil
}

func SendEmail(ctx context.Context, input SendEmailInput) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "testmail@dyne.org")
	message.SetHeader("To", input.Email)
	message.SetHeader("Subject", "YAML Validation")
	message.SetBody("text/html", input.Message)
	dialer := gomail.NewDialer("mail.dyne.org", 465, "testmail@dyne.org", "odQxmE4LPakFDwrfNuPqbCLYZCg34A")
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}

func MakeAPICall(ctx context.Context, input APICallInput) error {
	client := pocketbase.NewClient("http://localhost:8090")
	logger := activity.GetLogger(ctx)
	record, err := client.One("test", "25mu78f4uy7o848")
	if err != nil {
		logger.Error("Failed to make API call", "Error", err)
		return err
	}
	logger.Debug("API call response", "Response", record)
	logger.Debug("Quanto: ", record["quanto"])
	quanto, ok := record["quanto"].(float64)
	if !ok {
		logger.Error("Failed to convert 'quanto' to int")
		return errors.New("type assertion to int failed")
	}
	value := quanto + 1
	errUpdate := client.Update("test", "25mu78f4uy7o848", map[string]any{"quanto": value})
	if errUpdate != nil {
		logger.Error("Failed to update record", "Error", errUpdate)
		return errUpdate
	}
	return nil
}
