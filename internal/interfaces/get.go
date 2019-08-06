package interfaces

type IGet interface {
	Get(string) interface{}
}

type IGetMap interface {
	GetMap() map[string]interface{}
}
