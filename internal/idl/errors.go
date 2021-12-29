package idl

import "fmt"

type DecodeError struct {
	Types       Tuple
	Description string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("%s %s", e.Types.String(), e.Description)
}

type FormatError struct {
	Description string
}

func (e FormatError) Error() string {
	return fmt.Sprintf("() %s", e.Description)
}
