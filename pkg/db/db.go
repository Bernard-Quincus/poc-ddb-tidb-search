package db

type DB interface {
	Get(any) (any, error)
	Put(any) (any, error)
	Delete(any) error
	Search(any) (any, error)
	Close(any) error
	SetTableName(string) error
}
