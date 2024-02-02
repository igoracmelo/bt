package bt

import "fmt"

const (
	CodeNotFound = iota + 1
)

type Error struct {
	Code    int
	Details string
}

func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", CodeMessage[e.Code], e.Details)
}

var CodeMessage = map[int]string{
	CodeNotFound: "not found",
}
