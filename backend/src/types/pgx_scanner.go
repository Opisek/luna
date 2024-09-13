package types

type PgxScanner interface {
	Scan(dest ...interface{}) (err error)
}
