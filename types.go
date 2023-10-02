package azfunc

import (
	"encoding/json"
	"strconv"
)

// RawMessage is a type based on json.RawMessage with a custom UnmarshalJSON
// method to handle escaped JSON.
type RawMessage json.RawMessage

// UnmarshalJSON satisfies json.Unmarshaler. It unquotes
// escaped JSON if it's escaped, otherwise sets
// the data as is.
func (r *RawMessage) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err == nil {
		*r = RawMessage(unquoted)
		return nil
	}
	*r = RawMessage(b)

	return nil
}

// MarshalJSON satisfies json.Marshaler. If valid JSON
// is provided it will be escaped and returned as a
// JSON string, otherwise return the data as is.
func (r RawMessage) MarshalJSON() ([]byte, error) {
	var js json.RawMessage
	if err := json.Unmarshal(r, &js); err == nil {
		return json.Marshal(string(r))
	}
	return r, nil
}
