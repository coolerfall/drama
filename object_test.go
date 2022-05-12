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
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ObjectStruct struct {
	Name string
}

func (os *ObjectStruct) SetName(name string) {
	os.Name = name
}

func TestObjectAssign(t *testing.T) {
	os := &ObjectStruct{}
	obj := newObject(reflect.ValueOf(os))
	err := obj.Assign("Name", "obj")
	assert.Nil(t, err)
	assert.Equal(t, "obj", os.Name)
}

func TestObjectCall(t *testing.T) {
	os := &ObjectStruct{}
	obj := newObject(reflect.ValueOf(os))
	err, _ := obj.Call("SetName", "dynamic")
	assert.Nil(t, err)
	assert.Equal(t, "dynamic", os.Name)
}
