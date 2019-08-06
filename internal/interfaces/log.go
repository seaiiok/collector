package interfaces

type ILog interface {
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	IClose
}
