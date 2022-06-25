package protojson

import (
	"bytes"
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// This package provides default options over
// google.golang.org/protobuf/encoding/protojson.

type UnmarshalOptions = protojson.UnmarshalOptions
type MarshalOptions = protojson.MarshalOptions

func Unmarshal(b []byte, m proto.Message) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(b, m)
}

func UnmarshalArray(b []byte, fn func(buf json.RawMessage) error) error {
	d := json.NewDecoder(bytes.NewReader(b))

	_, err := d.Token()
	if err != nil {
		return err
	}
	for d.More() {
		var buf json.RawMessage
		if err := d.Decode(&buf); err != nil {
			return err
		}

		err := fn(buf)
		if err != nil {
			return err
		}
	}

	return nil
}

func Marshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(m)
}

func MarshalWithUnpopulated(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(m)
}
