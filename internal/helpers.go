package internal

import "google.golang.org/protobuf/encoding/protojson"

func NewInt64(i int64) *int64 { return &i }
func NewString(s string) *string { return &s }

var protojsonu = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}
