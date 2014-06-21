package config

import (
	"log"
	"os/user"
	"strings"
)

// C is the active configuration instance
var C ConfigSource

// Config represents the program configuration options
type Config struct {
	Host        string        `json:"host"`
	MediaFolder string        `json:"mediaFolder"`
	Sqlite      *SqliteConfig `json:"sqlite"`
}

// Media returns the media folder from config, but with special
// characters such as '~' replaced
func (c Config) Media() string {
	// Get current user
	user, err := user.Current()
	if err != nil {
		log.Println(err)
		return c.MediaFolder
	}

	// Return path with strings replaced, trailing slash removed
	return strings.TrimRight(strings.Replace(c.MediaFolder, "~", user.HomeDir, -1), "/\\")
}

// SqliteConfig represents configuration for an sqlite backend
type SqliteConfig struct {
	File string `json:"file"`
}

// ConfigSource represents the configuration source for the program
type ConfigSource interface {
	Load() (*Config, error)
	Use(string) error
}
