package interfaces

type IDatebase interface {
	Select(string) [][]string
	Update(string, [][]string)
	Insert(string, [][]string)
	Delete(string)
	IClose
}
