package yammya

type AppInput struct {
	Url   string
	Email string
}

type ValidationInput struct {
	Url string
}

type ValidationOutput struct {
	IsValid bool
	Message string
}

type SendEmailInput struct {
	Email   string
	Message string
}

type APICallInput struct {
	Url string
}

const (
	APICallURL = "https://api.github.com/users/me/"
	TaskQueue = "yammya"
)
