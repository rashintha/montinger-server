package exceptions

import "errors"

var (
	FailedToPushDeviceData        = errors.New("failed to push device data")
	InvalidAPIKey                 = errors.New("invalid api key")
	InvalidatedToken              = errors.New("invalid user token")
	InvalidInput                  = errors.New("invalid input")
	InvalidToken                  = errors.New("invalid token")
	InvalidUsernameOrPassword     = errors.New("username or password invalid")
	RegistrationFailed            = errors.New("registration failed")
	RequestBodyValidationFailed   = errors.New("request body validation failed")
	RequestParamsValidationFailed = errors.New("request params validation failed")
	ResourceNotFound              = errors.New("resource not found")
	UserAlreadyExists             = errors.New("user already exists in our database")
	UserNotFoundInOurDB           = errors.New("user does not exists in our database")
)
