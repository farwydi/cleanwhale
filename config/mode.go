package config

// Mode define configuration current app work mode.
type Mode string

const (
	// ModeLocal local environment, running on the developer's machine for example.
	ModeLocal Mode = "local"

	// ModeRelease production environment.
	ModeRelease Mode = "release"

	// ModeTest for testing environment.
	ModeTest Mode = "test"
)
