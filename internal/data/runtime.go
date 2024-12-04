package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

// Implement json.Marshaller interface by providing its MarshalJSON method on the Runtime type
func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	// wrap jsonValue in double quotes.
	quotedJSONValue := strconv.Quote(jsonValue)

	return []byte(quotedJSONValue), nil
}
