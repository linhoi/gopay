package paypal

import (
	"fmt"
)

type ErrorResp map[string]interface{}

// Err error caused by http request and response
func (e ErrorResp) Err(err error) error {
	if err != nil {
		return err
	}

	if e != nil {
		return fmt.Errorf("%v", e)
	}

	return nil
}
