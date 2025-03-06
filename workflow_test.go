package yammya

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_SuccessfulWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := AppInput{
		Email: "Ennio.fex@gmail.com",
		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/ValidYAML.yaml",
	}

	testDetailsII := ValidationInput{
		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/ValidYAML.yaml",
	}

	env.OnActivity(ValidateYaml, mock.Anything, testDetailsII).Return(ValidationOutput{IsValid: true, Message: ""}, nil)
	env.OnActivity(MakeAPICall, mock.Anything, APICallInput{Url: APICallURL}).Return(nil)

	env.ExecuteWorkflow(Validation, testDetails)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

}

func Test_InvalidYAMLWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := AppInput{
		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/InvalidYAML.yaml",
		Email: "user@example.com",
	}

	testDetailsII := ValidationInput{
		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/InvalidYAML.yaml",

	}

	env.OnActivity(ValidateYaml, mock.Anything, testDetailsII).Return(ValidationOutput{IsValid: false, Message: ""}, nil)
	env.OnActivity(SendEmail, mock.Anything, SendEmailInput{Email: testDetails.Email, Message: ""}).Return(nil)

	env.ExecuteWorkflow(Validation, testDetails)

	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}

func Test_UnsuccessfulWorkflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	testDetails := AppInput{
		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/ValidYAML.yaml",
	}

	testDetailsII := ValidationInput{
		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/ValidYAML.yaml",
	}

	env.OnActivity(ValidateYaml, mock.Anything, testDetailsII).Return(true, nil)
	env.OnActivity(MakeAPICall, mock.Anything, APICallInput{Url: APICallURL}).Return(errors.New("API call failed"))

	env.ExecuteWorkflow(Validation, testDetails)

	require.True(t, env.IsWorkflowCompleted())
	require.Error(t, env.GetWorkflowError())
}

// func Test_InvalidYAMLActivity(t *testing.T) {
// 	testSuite := &testsuite.WorkflowTestSuite{}
// 	env := testSuite.NewTestActivityEnvironment()

// 	testDetails := ValidationInput{
// 		Url: "https://raw.githubusercontent.com/phoebus-84/kaa/refs/heads/main/cmd/fixtures/InvalidYAML.yaml",
// 	}

// 	env.ExecuteActivity(ValidateYaml, testDetails)
// }

