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
	"embed"
	"io/fs"
)

var (
	globalPackage  = NewPackage()
	globalResource = NewResource()
)

// Import is global func for Package.Import.
func Import(funcOrStruct ...any) error {
	return globalPackage.Import(funcOrStruct...)
}

// Exported is global func for Package.Exported.
func Exported(name, fieldName string) bool {
	return globalPackage.Exported(name, fieldName)
}

// MakeOptFunc is global func for Package.MakeOptFunc.
func MakeOptFunc(name string, fields map[string]any) (any, error) {
	return globalPackage.MakeOptFunc(name, fields)
}

// Use is global func for Package.Use.
func Use(name string, argsOrFields ...any) (*Object, error) {
	return globalPackage.Use(name, argsOrFields...)
}

// Register register embed.FS into drama.
func Register(fss ...embed.FS) {
	for _, fs := range fss {
		globalResource.Register(fs)
	}
}

// Load loads resources file from drama with given path.
func Load(path string) (fs.File, error) {
	return globalResource.Load(path)
}
