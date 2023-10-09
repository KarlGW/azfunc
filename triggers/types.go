package triggers

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"unicode/utf8"
)

// RawData is a type based on []byte with a custom UnmarshalJSON
// method to handle escaped JSON.
type RawData []byte

// UnmarshalJSON satisfies json.Unmarshaler. It unquotes
// escaped JSON if it's escaped, otherwise sets
// the data as is.
func (r *RawData) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err == nil {
		innerUnquoted, innerErr := strconv.Unquote(unquoted)
		if innerErr == nil {
			*r = RawData(innerUnquoted)
			return nil
		}
		trimmed := trimDoubleQuotes(b)
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(trimmed)))
		if n, err := base64.StdEncoding.Decode(decoded, trimmed); err == nil {
			*r = RawData(decoded[:n])
			return nil
		}
		*r = RawData(unquoted)
		return nil
	}
	*r = RawData(b)
	return nil
}

// MarshalJSON satisfies json.Marshaler. If valid JSON
// is provided it will be escaped and returned as a
// JSON string, otherwise return the data as is.
func (r RawData) MarshalJSON() ([]byte, error) {
	var js any
	if err := json.Unmarshal(r, &js); err == nil {
		return json.Marshal(string(r))
	}
	if len(r) > 0 && utf8.Valid(r) && r[0] != '"' && r[len(r)-1] != '"' {
		return json.Marshal(string(r))
	}
	return json.Marshal(base64.StdEncoding.EncodeToString(r))
}

// trimDoubleQuotes removes leading and trailing double quotes.
func trimDoubleQuotes(b []byte) []byte {
	if len(b) > 0 && b[0] == '"' && b[len(b)-1] == '"' {
		return b[1 : len(b)-1]
	}
	return b
}
