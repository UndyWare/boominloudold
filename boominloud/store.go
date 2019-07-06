package boominloud

type mediaStore interface {
	Get(string) (interface{}, error)
	Put(string, interface{}) error
}
