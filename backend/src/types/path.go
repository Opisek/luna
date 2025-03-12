package types

import (
	"encoding/json"
)

type Path string

func (p *Path) MarshalJSON() ([]byte, error) {
	if p == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(p.String())
}

func (p *Path) UnmarshalJSON(data []byte) error {
	var path string
	if err := json.Unmarshal(data, &path); err != nil {
		return err
	}
	*p = Path(path)
	return nil
}

func (p *Path) String() string {
	return string(*p)
}

func NewPath(path string) (*Path, error) {
	return (*Path)(&path), nil
}
