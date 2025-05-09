package state

import (
	"gator/internal/config"
)

// State represents the application state and holds a pointer to the configuration
type State struct {
	// Config is a pointer to the application configuration
	Config *config.Config
}
