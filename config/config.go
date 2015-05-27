// Various stucts that are used to configure sources and listeners.
package config

type Config struct {
	Sources   []Source
	Listeners map[string][]Listener
}
