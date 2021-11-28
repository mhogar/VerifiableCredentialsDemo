package handlers

const (
	TypeNoError       = iota
	TypeClientError   = iota
	TypeInternalError = iota
)

type CustomError struct {
	Type    int
	Message string
}

func NoError() CustomError {
	return CustomError{
		Type: TypeNoError,
	}
}

func ClientError(message string) CustomError {
	return CustomError{
		Type:    TypeClientError,
		Message: message,
	}
}

func InternalError() CustomError {
	return CustomError{
		Type:    TypeInternalError,
		Message: "An internal error occurred.",
	}
}
