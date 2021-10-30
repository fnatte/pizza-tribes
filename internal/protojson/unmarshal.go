package protojson

import (
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

func Marshal(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{}.Marshal(m)
}

func MarshalWithUnpopulated(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{
		EmitUnpopulated: true,
	}.Marshal(m)
}
