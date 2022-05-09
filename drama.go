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
)

type Drama struct {
	registry map[string]interface{}
}

func NewDrama() *Drama {
	return &Drama{
		registry: make(map[string]interface{}, 0),
	}
}

// Import imports new funcs or option types into registry.
func (d *Drama) Import(funcOrStruct ...interface{}) error {
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
				return errors.New("only one out parameter support")
			}

			fptr := reflect.ValueOf(fs).Pointer()
			fullName = runtime.FuncForPC(fptr).Name()
		} else if fsType.Kind() == reflect.Ptr {
			realType := fsType.Elem()
			fullName = fmt.Sprintf("%s.%s", realType.PkgPath(), realType.Name())
		}

		d.registry[fullName] = fs
	}

	return nil
}

// MakeOptFunc makes a option func variable with name and fields.
func (d *Drama) MakeOptFunc(name string, fields map[string]interface{}) (interface{}, error) {
	opt, ok := d.registry[name]
	if !ok {
		return nil, fmt.Errorf("no type found for: %v", name)
	}

	optType := reflect.TypeOf(opt)
	optFuncType := reflect.FuncOf([]reflect.Type{optType}, []reflect.Type{}, false)
	optProxyFunc := reflect.MakeFunc(optFuncType, func(args []reflect.Value) []reflect.Value {
		if len(args) != 1 {
			return nil
		}

		opt := args[0]
		for k, v := range fields {
			field := opt.Elem().FieldByName(k)

			if !field.CanSet() {
				continue
			}
			field.Set(reflect.ValueOf(v))
		}
		return nil
	})
	optFunc := reflect.New(optFuncType)
	optFunc.Elem().Set(optProxyFunc)

	return optFunc.Elem().Interface(), nil
}

// Use makes a new Object with given func name and arguments.
func (d *Drama) Use(name string, args ...interface{}) (*Object, error) {
	fn, ok := d.registry[name]
	if !ok {
		return nil, fmt.Errorf("no func found for: %v", name)
	}

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return nil, fmt.Errorf("%v is not a func", name)
	}

	in := make([]reflect.Value, 0)
	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}

	return newObject(reflect.ValueOf(fn).Call(in)[0]), nil
}
