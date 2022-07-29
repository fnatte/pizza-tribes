package models

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func (p *PatchMask) AppendPath(path string) {
	for _, p := range p.Paths {
		if p == path {
			return
		}
	}

	p.Paths = append(p.Paths, path)
}

func GetValueByPath(m proto.Message, path string) (interface{}, error) {
	segments := strings.Split(path, ".")

	fd := m.ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name(segments[0]))
	if fd == nil {
		return nil, fmt.Errorf("could not find field '%s' in message", segments[0])
	}
	fv := m.ProtoReflect().Get(fd)
	wasList := false

	for _, segment := range segments[1:] {
		if fd.IsList() && !wasList {
			n, err := strconv.Atoi(segment)
			if err != nil {
				return nil, fmt.Errorf("can not use '%s' as key in list", segment)
			}
			fv = fv.List().Get(n)
			wasList = true
		} else if fd.IsMap() {
			fv = fv.Map().Get(protoreflect.ValueOf(segment).MapKey())
			fd = fd.MapValue()
			wasList = false
		} else if fd.Kind() == protoreflect.MessageKind {
			mf := fv.Message()
			fd = mf.Descriptor().Fields().ByName(protoreflect.Name(segment))
			fv = mf.Get(fd)
			wasList = false
		}

		if fd == nil {
			return nil, fmt.Errorf("could not find '%s' in message as part of path '%s'", segment, path)
		}
	}

	if fd.IsList() {
		mt, err := protoregistry.GlobalTypes.FindMessageByName(fd.Message().FullName())
		if err != nil {
			return nil, fmt.Errorf("could not find message type: %w", err)
		}

		list := fv.List()

		val := mt.New().Interface()
		typ := reflect.SliceOf(reflect.TypeOf(val))
		slice := reflect.MakeSlice(typ, 0, list.Len())

		for i := 0; i < list.Len(); i++ {
			slice = reflect.Append(slice, reflect.ValueOf(list.Get(i).Message().Interface()))
		}

		return slice.Interface(), nil
	}

	if fd.IsMap() {
		mt, err := protoregistry.GlobalTypes.FindMessageByName(fd.MapValue().Message().FullName())
		if err != nil {
			return nil, fmt.Errorf("could not find message type: %w", err)
		}

		val := mt.New().Interface()
		typ := reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(val))
		m := reflect.MakeMap(typ)

		fv.Map().Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
			m.SetMapIndex(reflect.ValueOf(k.String()), reflect.ValueOf(v.Message().Interface()))

			return true
		})

		return m.Interface(), nil
	}

	if fd.Kind() == protoreflect.MessageKind {
		if fv.IsValid() {
			return fv.Message().Interface(), nil
		} else {
			return nil, nil
		}
	}

	return fv.Interface(), nil
}

func rangeFields(path string, f func(field string) bool) bool {
	for {
		var field string
		if i := strings.IndexByte(path, '.'); i >= 0 {
			field, path = path[:i], path[i:]
		} else {
			field, path = path, ""
		}

		if !f(field) {
			return false
		}

		if len(path) == 0 {
			return true
		}
		path = strings.TrimPrefix(path, ".")
	}
}
