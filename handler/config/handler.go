package config

type Handler interface {
	GetType() string
	GetMessage() string
	String() string
	Init(s string)
}
