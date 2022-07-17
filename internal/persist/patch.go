package persist

import (
	"strings"

	"github.com/fnatte/pizza-tribes/internal/models"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

/*
type Operation struct {
	Op    string
	Path  string
	Value interface{}
	From  string
}
*/
type Operation = models.JsonPatchOp

type Patch struct {
	// Ops []Operation
	Ops []*models.JsonPatchOp
}

func NewPatch() Patch {
	return Patch{
		// Ops: []Operation{},
		Ops: []*models.JsonPatchOp{},
	}
}

func (p *Patch) Add(path string, value interface{}) {
	/*
	p.Ops = append(p.Ops, Operation{
		Op:    "add",
		Path:  path,
		Value: value,
	})
	*/
	val, err := structpb.NewValue(value)
	if err != nil {
		panic(err);
	}
	p.Ops = append(p.Ops, &models.JsonPatchOp{
		Op:    "add",
		Path:  path,
		Value: val,
	})
}

func (p *Patch) Replace(path string, value interface{}) {
	/*
	p.Ops = append(p.Ops, Operation{
		Op:    "replace",
		Path:  path,
		Value: value,
	})
	*/

	if v, ok := value.(proto.Message); ok {
		// value = protov
	}

	val, err := structpb.NewValue(value)
	if err != nil {
		panic(err);
	}
	p.Ops = append(p.Ops, &models.JsonPatchOp{
		Op:    "replace",
		Path:  path,
		Value: val,
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
