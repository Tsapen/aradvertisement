package postgres

import (
	"github.com/Tsapen/aradvertisement/internal/ara"
	"github.com/lib/pq"
)

const (
	dataException            = "22"
	uniqueConstrainViolation = "23"
)

func translateError(err error) error {
	if errType, ok := err.(*pq.Error); ok {
		var class = errType.Code.Class()

		if class == dataException ||
			class == uniqueConstrainViolation {
			return ara.ErrBadParameters
		}
	}

	return err
}
