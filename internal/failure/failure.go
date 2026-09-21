// Package failure carries transport-neutral, safe error categories.
package failure

import "errors"

type Kind string

const (
	Invalid        Kind = "invalid"
	TooLarge       Kind = "too_large"
	NoRequirements Kind = "no_requirements"
	Capacity       Kind = "capacity"
	Configuration  Kind = "configuration"
	Evaluation     Kind = "evaluation"
	Unavailable    Kind = "unavailable"
	Timeout        Kind = "timeout"
	Internal       Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
}

func (e *Error) Error() string            { return e.Message }
func New(kind Kind, message string) error { return &Error{kind, message} }
func Classify(err error) Kind {
	var e *Error
	if errors.As(err, &e) {
		return e.Kind
	}
	return Internal
}
