// Package transport determines the shape of the available transports.
package transport

import "context"

// Server the basic structure of the transport allows you to register objects.
type Server interface {

	// Run register transport.
	Run(ctx context.Context) error
}
