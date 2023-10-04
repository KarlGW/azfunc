package azfunc

import (
	"encoding/base64"
	"encoding/json"
	"strconv"
	"unicode/utf8"
)

// Payload is a type based on []byte with a custom UnmarshalJSON
// method to handle escaped JSON.
type Payload []byte

// UnmarshalJSON satisfies json.Unmarshaler. It unquotes
// escaped JSON if it's escaped, otherwise sets
// the data as is.
func (p *Payload) UnmarshalJSON(b []byte) error {
	unquoted, err := strconv.Unquote(string(b))
	if err == nil {
		innerUnquoted, innerErr := strconv.Unquote(unquoted)
		if innerErr == nil {
			*p = Payload(innerUnquoted)
			return nil
		}
		trimmed := trimDoubleQuotes(b)
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(trimmed)))
		if n, err := base64.StdEncoding.Decode(decoded, trimmed); err == nil {
			*p = Payload(decoded[:n])
			return nil
		}
		*p = Payload(unquoted)
		return nil
	}
	*p = Payload(b)
	return nil
}

// MarshalJSON satisfies json.Marshaler. If valid JSON
// is provided it will be escaped and returned as a
// JSON string, otherwise return the data as is.
func (p Payload) MarshalJSON() ([]byte, error) {
	var js any
	if err := json.Unmarshal(p, &js); err == nil {
		return json.Marshal(string(p))
	}
	if len(p) > 0 && utf8.Valid(p) && p[0] != '"' && p[len(p)-1] != '"' {
		return json.Marshal(string(p))
	}
	return json.Marshal(base64.StdEncoding.EncodeToString(p))
}

// trimDoubleQuotes removes leading and trailing double quotes.
func trimDoubleQuotes(b []byte) []byte {
	if len(b) > 0 && b[0] == '"' && b[len(b)-1] == '"' {
		return b[1 : len(b)-1]
	}
	return b
}
