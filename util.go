package middleware

// prepareInputChannel converts a channel of type interface{} to a
// write-only channel of type interface{}
func prepareInputChannel(c chan interface{}) chan<- interface{} {
	return c
}

// prepareStatusChannel converts a channel of type bool to a
// read-only channel of type bool.
func prepareStatusChannel(c chan bool) <-chan bool {
	return c
}
