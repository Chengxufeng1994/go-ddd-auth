package assert

import "fmt"

func AssertArgumentNotEqual[T comparable](a T, b T) error {
	if a == b {
		return fmt.Errorf("arguments cannot be equal")
	}

	return nil
}
