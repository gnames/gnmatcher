package binary

import (
	"bytes"
	"encoding/gob"
)

func Encode(input interface{}) ([]byte, error) {
	var respBytes bytes.Buffer
	enc := gob.NewEncoder(&respBytes)
	if err := enc.Encode(input); err != nil {
		return nil, err
	}
	return respBytes.Bytes(), nil
}

func Decode(input []byte, output interface{}) {
	b := bytes.NewBuffer(input)
	dec := gob.NewDecoder(b)
	dec.Decode(output)
}
