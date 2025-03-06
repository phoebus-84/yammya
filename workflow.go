package yammya

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func Validation(ctx workflow.Context, input AppInput) error {

	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    500,
	}

	options := workflow.ActivityOptions{
		TaskQueue: TaskQueue,
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retrypolicy,
	}

	// Apply the options.
	ctx = workflow.WithActivityOptions(ctx, options)

	var output ValidationOutput

	ValidationErr := workflow.ExecuteActivity(ctx, ValidateYaml, input).Get(ctx, &output)

	if ValidationErr != nil {
		return ValidationErr
	}

	if !output.IsValid {
		sendEmailInput := SendEmailInput{
			Email:   input.Email,
			Message: output.Message,
		}
		err := workflow.ExecuteActivity(ctx, SendEmail, sendEmailInput).Get(ctx, nil)

		if err != nil {
			return err
		}
		return nil
	} else {
		input := APICallInput{
			Url: APICallURL,
		}
		err := workflow.ExecuteActivity(ctx, MakeAPICall, input).Get(ctx, nil)
		if err != nil {
			return err
		}
		return nil
	}
}


