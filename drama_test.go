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
	"testing"

	"github.com/stretchr/testify/assert"
)

type Checker struct {
	name    string
	version int
}

type CheckerOption struct {
	Version int
}

func NewChecker(opts ...func(*CheckerOption)) *Checker {
	opt := &CheckerOption{
		Version: 1,
	}

	for _, f := range opts {
		f(opt)
	}

	return &Checker{
		version: opt.Version,
	}
}

func (c *Checker) WithName(name string) {
	c.name = name
}

func TestDynamicFunc(t *testing.T) {
	d := NewDrama()
	_ = d.Import(NewChecker)

	fn, err := d.Use("github.com/coolerfall/drama.NewChecker")
	assert.Nil(t, err)

	err, _ = fn.Call("WithName", "dynamic")
	assert.Nil(t, err)

	c, ok := fn.Itf().(*Checker)
	assert.Equal(t, true, ok)
	assert.Equal(t, "dynamic", c.name)
}

func TestOptFunc(t *testing.T) {
	d := NewDrama()
	_ = d.Import(NewChecker, (*CheckerOption)(nil))

	var cf = map[string]interface{}{"Version": 6}
	optFn, err := d.MakeOptFunc("github.com/coolerfall/drama.CheckerOption", cf)
	assert.Nil(t, err)
	fn, err := d.Use("github.com/coolerfall/drama.NewChecker", optFn)
	assert.Nil(t, err)
	c, ok := fn.Itf().(*Checker)
	assert.Equal(t, true, ok)
	assert.Equal(t, 6, c.version)
}
