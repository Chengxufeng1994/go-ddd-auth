package assert

import (
	"fmt"
	"strings"
)

func AssertArgumentEmpty(arg string) error {
	if strings.Trim(arg, " ") == "" {
		return nil
	}

	return fmt.Errorf("argument can be empty")
}
