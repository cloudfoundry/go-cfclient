package cfclient

import "fmt"

type CloudFoundryErrors struct {
	Errors []CloudFoundryError `json:"errors"`
}

func (cfErrs CloudFoundryErrors) Error() string {
	err := ""

	for _, cfErr := range cfErrs.Errors {
		err += fmt.Sprintf("%s\n", cfErr)
	}

	return err
}

type CloudFoundryError struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (cfErr CloudFoundryError) Error() string {
	return fmt.Sprintf("cfclient: error (%d): %s", cfErr.Code, cfErr.Title)
}
