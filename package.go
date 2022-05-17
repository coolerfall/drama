// Copyright (c) 2022 Vincent Cheung (coolingfall@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package drama

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"

	"github.com/mitchellh/mapstructure"
)

type Package struct {
	registry map[string]any
}

// NewPackage creates a new instance for Package.
func NewPackage() *Package {
	return &Package{
		registry: make(map[string]any, 0),
	}
}

// Import imports new funcs or struct types into registry.
func (p *Package) Import(funcOrStruct ...any) error {
	if len(funcOrStruct) == 0 {
		return errors.New("nothing to import")
	}

	for _, fs := range funcOrStruct {
		if fs == nil {
			return errors.New("cannot import nil")
		}

		var fullName string
		fsType := reflect.TypeOf(fs)
		if fsType.Kind() == reflect.Func {
			if fsType.NumOut() != 1 {
				return errors.New("only func with one out parameter supported")
			}

			outType := fsType.Out(0)
			outKind := outType.Kind()
			if outKind != reflect.Interface && outKind != reflect.Pointer &&
				outType.Elem().Kind() != reflect.Struct {
				return errors.New("only func with struct pointer or " +
					"interface out parameter supported")
			}

			fptr := reflect.ValueOf(fs).Pointer()
			fullName = runtime.FuncForPC(fptr).Name()
		} else if fsType.Kind() == reflect.Pointer && fsType.Elem().Kind() == reflect.Struct {
			realType := fsType.Elem()
			fullName = fmt.Sprintf("%s.%s", realType.PkgPath(), realType.Name())
		} else {
			return errors.New("only func or struct type supported")
		}

		p.registry[fullName] = fs
	}

	return nil
}

// HasExportedField checks if struct has exported field with given package path.
func (p *Package) HasExportedField(name, fieldName string) bool {
	opt, ok := p.registry[name]
	if !ok {
		return false
	}

	optType := reflect.TypeOf(opt)
	field, ok := optType.Elem().FieldByName(fieldName)
	if !ok {
		return false
	}

	return field.IsExported()
}

// MakeOptFunc makes a option func variable package with name and fields.
func (p *Package) MakeOptFunc(name string, fields map[string]any) (any, error) {
	opt, ok := p.registry[name]
	if !ok {
		return nil, fmt.Errorf("no type found for: %v", name)
	}

	optType := reflect.TypeOf(opt)

	for k := range fields {
		field, ok := optType.Elem().FieldByName(k)
		if !ok {
			return nil, fmt.Errorf("cannot find field '%s' in %s", k, name)
		}

		if !field.IsExported() {
			return nil, fmt.Errorf("field '%s' is not exported", k)
		}
	}

	optFuncType := reflect.FuncOf([]reflect.Type{optType}, []reflect.Type{}, false)
	optProxyFunc := reflect.MakeFunc(optFuncType, func(args []reflect.Value) []reflect.Value {
		if len(args) != 1 {
			return nil
		}

		// map all fields to the option
		_ = mapstructure.WeakDecode(fields, args[0].Interface())

		return nil
	})
	optFunc := reflect.New(optFuncType)
	optFunc.Elem().Set(optProxyFunc)

	return optFunc.Elem().Interface(), nil
}

// Use makes a new Object with given package path and arguments.
func (p *Package) Use(name string, argsOrFields ...any) (*Object, error) {
	fnOrStruct, ok := p.registry[name]
	if !ok {
		return nil, fmt.Errorf("no func found for: %v", name)
	}

	tp := reflect.TypeOf(fnOrStruct)
	kind := tp.Kind()
	if kind == reflect.Func {
		in := make([]reflect.Value, 0)
		for _, arg := range argsOrFields {
			in = append(in, reflect.ValueOf(arg))
		}
		return newObject(reflect.ValueOf(fnOrStruct).Call(in)[0]), nil
	} else {
		st := reflect.New(tp.Elem())

		if len(argsOrFields) == 0 {
			return newObject(st), nil
		}

		fields, ok := argsOrFields[0].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("args should be a map when using struct")
		}

		// check if field can be set
		for k := range fields {
			field := st.Elem().FieldByName(k)
			if !field.CanSet() {
				return nil, fmt.Errorf("cannot find exported field '%v' in %s", k, name)
			}
		}

		// map all fields to struct
		if err := mapstructure.WeakDecode(fields, st.Interface()); err != nil {
			return nil, err
		}

		return newObject(st), nil
	}
}
