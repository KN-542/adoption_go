package validator

import (
	"errors"
	"fmt"
)

// []stringの各要素が互いに異なるか？
type UniqueValidator struct{}

func (v UniqueValidator) Validate(value interface{}) error {
	s, ok := value.([]string)
	if !ok {
		return errors.New("invalid data type")
	}
	seen := make(map[string]struct{})
	for _, item := range s {
		if _, exists := seen[item]; exists {
			return fmt.Errorf("duplicate value: %s", item)
		}
		seen[item] = struct{}{}
	}
	return nil
}

// uintであるかの検証
type IsUintValidator struct{}

func (v IsUintValidator) Validate(value interface{}) error {
	_, ok := value.(uint)
	if !ok {
		return errors.New("invalid type to uint")
	}

	return nil
}
