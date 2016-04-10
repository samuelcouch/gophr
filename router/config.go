package main

// TODO (Sandile): make this module a real one

var _config *Config

func getConfig() *Config {
	if _config == nil {
		_config = &Config{
			dev:    true,
			domain: "gophr.dev",
		}
	}

	return _config
}

// Config keeps track of environment related configuration variables that
// affect server behavior and execution.
type Config struct {
	dev    bool
	domain string
}