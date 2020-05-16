package ara

import "fmt"

// Error implements error interface.
type Error string

func (err Error) Error() string {
	return string(err)
}

const (
	// ErrBadParameters is bad parameters error.
	ErrBadParameters Error = "bad parameters"
	// ErrNotAccess is access lack error.
	ErrNotAccess Error = "not access"
)

// BadParametersError implements error interface.
type BadParametersError struct {
	Err error
}

func (err BadParametersError) Error() string {
	return err.Err.Error()
}

// ErrPair contains deferred and returned error.
type ErrPair struct {
	Def error
	Ret error
}

func (errPair ErrPair) Error() string {
	return fmt.Sprintf("returned: %s; deferred: %s", errPair.Def, errPair.Ret)
}

// HandleErrPair contains deferred and returned errors.
func HandleErrPair(def, ret error) error {
	if ret == nil {
		return def
	}

	if def == nil {
		return ret
	}

	return ErrPair{
		Def: def,
		Ret: ret,
	}
}
