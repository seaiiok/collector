package interfaces

type ICache interface {
	IGet
	ISet
	IGetMap
	IClose
}
