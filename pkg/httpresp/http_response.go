package httpresp

import (
	"errors"
	"fmt"

	fastshot "github.com/opus-domini/fast-shot"
)

func ParseHTTPResponse[T any](resp fastshot.Response, result *T) error {
	if resp.Status().IsError() {
		msg, err := resp.Body().AsString()
		if err != nil {
			return fmt.Errorf("failed to read error response: %w", err)
		}
		return errors.New(msg)
	}

	err := resp.Body().AsJSON(result)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
