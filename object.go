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
	"fmt"
	"reflect"
)

type Object struct {
	objValue reflect.Value
}

func newObject(obj reflect.Value) *Object {
	return &Object{
		objValue: obj,
	}
}

// Itf returns value of current Object.
func (o *Object) Itf() interface{} {
	return o.objValue.Interface()
}

// Call invokes func on this object with or without arguments.
func (o *Object) Call(name string, args ...interface{}) (error, []interface{}) {
	in := make([]reflect.Value, 0)

	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}

	method := o.objValue.MethodByName(name)
	if !method.IsValid() || method.IsNil() {
		elemName := o.objValue.Type().Elem().Name()
		return fmt.Errorf("func '%s' not found on '%v'", name, elemName), nil
	}
	out := method.Call(in)
	if len(out) == 0 {
		return nil, nil
	}

	itfs := make([]interface{}, 0)
	for _, v := range out {
		itfs = append(itfs, v.Interface())
	}

	return nil, itfs
}
