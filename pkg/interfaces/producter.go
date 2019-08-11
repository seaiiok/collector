package interfaces

type IProducter interface {
	Run()
	Stop()
	Set(string, []byte)
	Get(string) []byte
	OutQueue() chan string
}
