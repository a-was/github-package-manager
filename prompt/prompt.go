package prompt

import (
	"fmt"
	"strconv"
	"strings"
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

func PrintTable(table [][]string) {
	// get number of columns from the first table row
	columnLengths := make([]int, len(table[0]))
	for _, line := range table {
		for i, val := range line {
			if len(val) > columnLengths[i] {
				columnLengths[i] = len(val)
			}
		}
	}

	var lineLength int
	for _, c := range columnLengths {
		lineLength += c + 3 // +3 for 3 additional characters before and after each field: "| %s "
	}
	lineLength += 1 // +1 for the last "|" in the line

	for i, line := range table {
		if i == 0 { // table header
			fmt.Printf("+%s+\n", strings.Repeat("-", lineLength-2)) // lineLength-2 because of "+" as first and last character
		}
		for j, val := range line {
			// change table header to upper
			if i == 0 {
				val = strings.ToUpper(val)
			}
			fmt.Printf("| %-*s ", columnLengths[j], val)
			if j == len(line)-1 {
				fmt.Printf("|\n")
			}
		}
		if i == 0 || i == len(table)-1 { // table header or last line
			fmt.Printf("+%s+\n", strings.Repeat("-", lineLength-2)) // lineLength-2 because of "+" as first and last character
		}
	}
}
