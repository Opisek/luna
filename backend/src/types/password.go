package types

type PasswordEntry struct {
	Hash       []byte
	Salt       []byte
	Algorithm  string
	Parameters map[string]int
}
