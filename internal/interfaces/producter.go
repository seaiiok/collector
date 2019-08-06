package interfaces

type IProducter interface {
	Run()
	Stop()
	Output() chan string
	DoneAFile(string)
}
