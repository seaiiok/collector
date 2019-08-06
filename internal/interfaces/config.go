package interfaces

type IConfig interface {
	GetString(string) string
	IGet
	IClose
}
