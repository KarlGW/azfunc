package data

// JSONMarshaler is an interface that types can implement to provide a JSON representation of themselves.
type JSONMarshaler interface {
	JSON() []byte
}
