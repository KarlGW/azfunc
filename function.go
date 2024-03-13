package azfunc

// function is an internal structure that represents a function
// in a FunctionApp.
type function struct {
	name     string
	trigger  triggerable
	bindings []bindable
}

// FunctionOption sets options to the function.
type FunctionOption func(f *function)

// Trigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func Trigger(name string, fn TriggerFunc) FunctionOption {
	return func(f *function) {
		f.trigger = trigger{
			fn:   fn,
			name: name,
		}
	}
}

// HTTPTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func HTTPTrigger(fn HTTPTriggerFunc) FunctionOption {
	return func(f *function) {
		f.trigger = httpTrigger{
			fn:   fn,
			name: "req",
		}
	}
}

// TimerTrigger takes the provided function and sets it as
// the function to be run by the trigger.
func TimerTrigger(fn TimerTriggerFunc) FunctionOption {
	return func(f *function) {
		f.trigger = timerTrigger{
			fn:   fn,
			name: "timer",
		}
	}
}

// QueueTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func QueueTrigger(name string, fn QueueTriggerFunc) FunctionOption {
	return func(f *function) {
		f.trigger = queueTrigger{
			fn:   fn,
			name: name,
		}
	}
}

// ServiceBusTrigger takes the provided name and function and sets it as
// the function to be run by the trigger.
func ServiceBusTrigger(name string, fn ServiceBusTriggerFunc) FunctionOption {
	return func(f *function) {
		f.trigger = serviceBusTrigger{
			fn:   fn,
			name: name,
		}
	}
}

// Binding sets the provided binding to the function.
func Binding(binding bindable) FunctionOption {
	return func(f *function) {
		if f.bindings == nil {
			f.bindings = []bindable{binding}
			return
		}
		f.bindings = append(f.bindings, binding)
	}
}
