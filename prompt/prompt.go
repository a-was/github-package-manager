package prompt

import (
	"fmt"
	"strconv"
)

// when using int Get sets v to -1 if input was invalid
func Get[T string | int](text string, v *T) {
	fmt.Print(text)

	var input string
	fmt.Scanln(&input)

	switch v := any(v).(type) {
	case *string:
		*v = input
	case *int:
		parsed, err := strconv.Atoi(input)
		if err != nil {
			parsed = -1
		}
		*v = parsed
	}
}
