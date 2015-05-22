package config

type Config struct {
    Sources   []Source
    Listeners map[string][]Listener
}
