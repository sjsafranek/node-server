package models

type Database interface {
	Has(string) bool
	Create(string) (Bucket, error)
	Delete(string) error
	Get(string) (Bucket, error)
}

type Bucket interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Delete(string) error
}
