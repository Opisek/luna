package types

type PgxScannable interface {
	Scan(dest ...interface{}) (err error)
}
