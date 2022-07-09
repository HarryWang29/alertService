package store

type Store interface {
	Get(code, key string) (interface{}, error)
	GetString(code, key string) (string, error)
	GetLast(code string) (string, interface{})
	Set(code, key string, value interface{}) error
	Delete(code, key string)
}
