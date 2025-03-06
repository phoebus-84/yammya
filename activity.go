package yammya

import (
	"bytes"
	"context"

	"github.com/phoebus-84/kaa/cmd"
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
	return nil
}
