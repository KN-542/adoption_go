package validator

import (
	"errors"
	"fmt"
)

// []stringの各要素が互いに異なるか？
type UniqueValidator struct{}

func (v UniqueValidator) Validate(value interface{}) error {
	if value == nil {
		return nil
	}

	s, ok := value.([]string)
	if !ok {
		return errors.New("invalid data type: expected []string")
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

// uintの最小値
type MinUintValidator struct {
	Min uint
}

func (v MinUintValidator) Validate(value interface{}) error {
	val, ok := value.(uint)
	if !ok {
		return errors.New("invalid type to uint")
	}
	if val < v.Min {
		return errors.New("value is less than the minimum allowed")
	}
	return nil
}

// uintの最大値
type MaxUintValidator struct {
	Max uint
}

func (v MaxUintValidator) Validate(value interface{}) error {
	val, ok := value.(uint)
	if !ok {
		return errors.New("invalid type to uint")
	}
	if val > v.Max {
		return errors.New("value exceeds the maximum allowed")
	}
	return nil
}
