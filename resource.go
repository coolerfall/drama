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
	"errors"
	"io/fs"
)

type Resource struct {
	fss []embed.FS
}

func NewResource() *Resource {
	return &Resource{}
}

func (r *Resource) Register(fs embed.FS) {
	r.fss = append(r.fss, fs)
}

func (r *Resource) Load(path string) (fs.File, error) {
	for _, efs := range r.fss {
		file, err := efs.Open(path)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}

		return file, nil
	}

	return nil, fs.ErrNotExist
}
