package assert

import (
	"fmt"
	"strings"
)

func AssertArgumentNotEmpty(arg string) error {
	if strings.Trim(arg, " ") == "" {
		return fmt.Errorf("argument cannot be empty")
	}

	return nil
}
