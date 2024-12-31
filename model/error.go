package model

import "fmt"

// A ErrNotFound(UpdateTODO) expresses ...
type ErrNotFound struct {
	Message string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("ErrNotfound: %s", e.Message)
}
