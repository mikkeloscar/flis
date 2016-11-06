package context

import "fmt"

// Context is a map of objects, e.g. config and backend, to be passed around in
// different parts of the application.
type Context map[string]interface{}

// MustGet returns a context object by key name. It panics if the key doesn't
// exist.
func (c Context) MustGet(key string) interface{} {
	if v, ok := c[key]; ok {
		return v
	}

	panic(fmt.Sprintf("Key '%s' does not exist", key))
}

// Set adds an object to the context.
func (c Context) Set(key string, value interface{}) {
	c[key] = value
}
