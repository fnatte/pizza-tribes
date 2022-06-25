package persist

import (
	"strings"
)

type Operation struct {
	Op    string
	Path  string
	Value interface{}
	From  string
}

type Patch struct {
	Ops []Operation
}

func NewPatch() Patch {
	return Patch{
		Ops: []Operation{},
	}
}

func (p *Patch) Add(path string, value interface{}) {
	p.Ops = append(p.Ops, Operation{
		Op:    "add",
		Path:  path,
		Value: value,
	})
}

func (p *Patch) Replace(path string, value interface{}) {
	p.Ops = append(p.Ops, Operation{
		Op:    "replace",
		Path:  path,
		Value: value,
	})
}

func JsonPointerToJsonPath(s string) string {
	var b strings.Builder
	b.WriteRune('$')

	parts := strings.Split(s, "/")
	for _, part := range parts {
		if len(part) > 0 {
			b.WriteRune('.')
			b.WriteString(part)
		}
	}

	return b.String()
}
