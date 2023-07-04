package domain

type IMessanger interface {
	Send(msg string) error
}
