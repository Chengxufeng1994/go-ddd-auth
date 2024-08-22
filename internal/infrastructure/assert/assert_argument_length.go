package assert

import "fmt"

func AssertArgumentLength(arguments string, min, max int) error {
	if len(arguments) < min || len(arguments) > max {
		return fmt.Errorf("argument length must be between %d and %d", min, max)
	}
	return nil
}
