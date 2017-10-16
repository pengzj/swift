package common

import "fmt"

type Errors struct {
	Code string
	Message string
}

func (e Errors) Error() {
	return fmt.Printf("code: %v: %v", e.Code, e.Message)
}

