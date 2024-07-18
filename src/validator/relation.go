package validator

import (
	"fmt"
	"time"
)

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

// 2つの日付の時刻比較
func IsBeforeTime(value time.Time, value2 time.Time) error {
	if value.Year() < 1900 || value2.Year() < 1900 {
		return nil
	}

	if value2.Before(value) {
		return fmt.Errorf("the back and forth between the two times is incorrect")
	}

	return nil
}
