package prompt

import (
	"fmt"
	"strconv"
)

func Get[T string | int](text string, v *T) error {
	fmt.Print(text)

	var input string
	fmt.Scanln(&input)

	switch v := any(v).(type) {
	case *string:
		*v = input
	case *int:
		parsed, err := strconv.Atoi(input)
		if err != nil {
			return err
		}
		*v = parsed
	}
	return nil
}
