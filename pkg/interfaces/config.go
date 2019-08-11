package interfaces

type IConfig interface {
	Get(string) string
	Set(string, string)
	IClose
}
