package validator

import "fmt"

// 2つの文字列が一致
func IsEqualStrings(value *string, value2 *string) error {
	if value == nil || value2 == nil || *value == "" || *value2 == "" {
		return nil
	}

	if *value != *value2 {
		return fmt.Errorf("%v, %v are not equal", *value, *value2)
	}

	return nil
}
