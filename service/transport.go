package service

type Transport interface {
	Serve() error
}
