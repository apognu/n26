package api

import (
	"encoding/json"
	"io"
)

type JSON struct {
	object interface{}
}

func NewJSON(o interface{}) *JSON {
	return &JSON{object: o}
}

func (j JSON) Decode(r io.Reader) (interface{}, error) {
	err := json.NewDecoder(r).Decode(j.object)
	return j.object, err
}
