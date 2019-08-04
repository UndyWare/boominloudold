package server

type mediaStore interface {
	Get(string) (interface{}, error)
	Put(string, interface{}) error
}
