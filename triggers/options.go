package triggers

// Options contains options for functions and methods related
// to triggers.
type Options struct {
	Name string
}

// Option is a function that sets options on Options.
type Option func(*Options)

// WithName sets the trigger name to get data from. The name should
// match the incoming trigger (binding) name in function.json.
func WithName(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}
