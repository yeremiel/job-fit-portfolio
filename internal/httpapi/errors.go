package httpapi

import (
	"context"
	"errors"
	"job-fit/internal/failure"
)

// Machine-readable codes are centralized here; messages are not client branching keys.
const (
	codeInvalid       = "REQUEST_INVALID"
	codeMalformed     = "REQUEST_MALFORMED_JSON"
	codeNotFound      = "RESOURCE_NOT_FOUND"
	codeMethod        = "REQUEST_METHOD_NOT_ALLOWED"
	codeTooLarge      = "REQUEST_INPUT_TOO_LARGE"
	codeMedia         = "REQUEST_UNSUPPORTED_MEDIA_TYPE"
	codeValidation    = "REQUEST_VALIDATION_FAILED"
	codeConfiguration = "INTERNAL_CONFIGURATION_ERROR"
	codeInternal      = "INTERNAL_ERROR"
	codeEvaluation    = "INTERNAL_EVALUATION_FAILED"
	codeUnavailable   = "INTERNAL_DEPENDENCY_UNAVAILABLE"
	codeTimeout       = "INTERNAL_EVALUATION_TIMEOUT"
)

type fieldError struct {
	Field   string `json:"field"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
type apiError struct {
	status        int
	code, message string
	fields        []fieldError
}

func problem(status int) *apiError {
	codes := map[int]string{400: codeInvalid, 404: codeNotFound, 405: codeMethod, 413: codeTooLarge, 415: codeMedia, 422: codeValidation, 500: codeInternal, 502: codeEvaluation, 503: codeUnavailable, 504: codeTimeout}
	messages := map[int]string{400: "Request structure is invalid.", 404: "Resource not found.", 405: "Method not allowed.", 413: "Input exceeds the allowed size.", 415: "Use application/json with UTF-8 and no content encoding.", 422: "Request validation failed.", 500: "An internal error occurred.", 502: "Evaluation could not be completed.", 503: "Evaluation dependency is unavailable.", 504: "Evaluation timed out."}
	return &apiError{status: status, code: codes[status], message: messages[status]}
}
func malformed() *apiError {
	e := problem(400)
	e.code = codeMalformed
	e.message = "Request must contain one valid UTF-8 JSON object."
	return e
}
func validation(field, reason, message string) *apiError {
	e := problem(422)
	e.fields = []fieldError{{field, reason, message}}
	return e
}
func mapError(err error) *apiError {
	if errors.Is(err, context.DeadlineExceeded) {
		return problem(504)
	}
	switch failure.Classify(err) {
	case failure.Invalid:
		return validation("job.description", "invalid_format", "Description is invalid.")
	case failure.TooLarge:
		return problem(413)
	case failure.NoRequirements:
		return validation("job.description", "no_capability_requirements", "No in-scope capability requirements were identified.")
	case failure.Capacity:
		return validation("job.description", "evaluation_capacity_exceeded", "Description exceeds evaluation capacity.")
	case failure.Configuration:
		e := problem(500)
		e.code = codeConfiguration
		e.message = "Evaluation service is not configured correctly."
		return e
	case failure.Evaluation:
		return problem(502)
	case failure.Unavailable:
		return problem(503)
	case failure.Timeout:
		return problem(504)
	default:
		return problem(500)
	}
}
