package server

type Option interface {
	Apply(arg ...interface{})

	ServiceName() string

	ProtocolName() string
}

type Options func(Option)
