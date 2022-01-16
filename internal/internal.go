package internal

type Logger interface {
	Error(e error)
	Fatal(e error)
	Info(t string)
}
