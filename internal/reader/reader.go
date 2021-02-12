package reader

type Reader interface {
	Read(path string) ([]string, error)
}
