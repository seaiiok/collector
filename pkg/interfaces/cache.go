package interfaces

type ICache interface {
	Get(string) []byte
	Set(string, []byte)
	Delete(string)
	GetMap() (map[string][]byte, error)
	IClose
}

