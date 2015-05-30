package config

type Listener interface {
	GetType() string
	GetMessage() string
	String() string
	Init(s string)
}
