package crud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/sefaphlvn/bigbang/pkg/errstr"
	"github.com/sefaphlvn/bigbang/pkg/helper"
	"github.com/sefaphlvn/bigbang/pkg/models"
)

func extractValidationErrors(err error) []string {
	if err == nil {
		return nil
	}

	errMsg := err.Error()
	errors := strings.Split(errMsg, ";")
	var result []string
	for _, e := range errors {
		e = strings.TrimSpace(e)
		if e != "" {
			result = append(result, e)
		}
	}

	return result
}

func Validate(gtype models.GTypes, resource interface{}) ([]string, bool, error) {
	msg := gtype.ProtoMessage()
	if msg == nil {
		return nil, true, fmt.Errorf("no message found for GType %v", gtype)
	}

	switch reflect.TypeOf(resource).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(resource)
		var allErrors []string

		for i := 0; i < s.Len(); i++ {
			elem := s.Index(i).Interface()
			if err := validateSingleResource(gtype, elem); err != nil {
				allErrors = append(allErrors, extractValidationErrors(err)...)
			}
		}
		if len(allErrors) > 0 {
			return allErrors, true, errstr.ErrValidationFailed
		}
	default:
		if err := validateSingleResource(gtype, resource); err != nil {
			return extractValidationErrors(err), true, errstr.ErrValidationFailed
		}
	}

	return nil, false, nil
}

func validateSingleResource(gtype models.GTypes, resource interface{}) error {
	msg := gtype.ProtoMessage()
	resourceBytes, err := json.Marshal(resource)
	if err != nil {
		return fmt.Errorf("failed to marshal resource: %w", err)
	}

	if err := helper.Unmarshaler.Unmarshal(resourceBytes, msg); err != nil {
		return fmt.Errorf("failed to unmarshal resource: %w", err)
	}

	validatable, ok := msg.(interface{ ValidateAll() error })
	if !ok {
		return fmt.Errorf("GType %T does not implement ValidateAll()", msg)
	}

	return validatable.ValidateAll()
}
