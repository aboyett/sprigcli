/*
Copyright 2016 The Kubernetes Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chartutil

import (
	"testing"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	path, data string
}{
	{"ship/captain.txt", "The Captain"},
	{"ship/stowaway.txt", "Legatt"},
	{"story/name.txt", "The Secret Sharer"},
	{"story/author.txt", "Joseph Conrad"},
	{"multiline/test.txt", "bar\nfoo"},
}

func getTestFiles() []*any.Any {
	a := []*any.Any{}
	for _, c := range cases {
		a = append(a, &any.Any{TypeUrl: c.path, Value: []byte(c.data)})
	}
	return a
}

func TestNewFiles(t *testing.T) {
	files := NewFiles(getTestFiles())
	if len(files) != len(cases) {
		t.Errorf("Expected len() = %d, got %d", len(cases), len(files))
	}

	for i, f := range cases {
		if got := string(files.GetBytes(f.path)); got != f.data {
			t.Errorf("%d: expected %q, got %q", i, f.data, got)
		}
		if got := files.Get(f.path); got != f.data {
			t.Errorf("%d: expected %q, got %q", i, f.data, got)
		}
	}
}

func TestFileGlob(t *testing.T) {
	as := assert.New(t)

	f := NewFiles(getTestFiles())

	matched := f.Glob("story/**")

	as.Len(matched, 2, "Should be two files in glob story/**")
	as.Equal("Joseph Conrad", matched.Get("story/author.txt"))
}

func TestToConfig(t *testing.T) {
	as := assert.New(t)

	f := NewFiles(getTestFiles())
	out := f.Glob("**/captain.txt").AsConfig()
	as.Equal("captain.txt: The Captain\n", out)

	out = f.Glob("ship/**").AsConfig()
	as.Equal("captain.txt: The Captain\nstowaway.txt: Legatt\n", out)
}

func TestToSecret(t *testing.T) {
	as := assert.New(t)

	f := NewFiles(getTestFiles())

	out := f.Glob("ship/**").AsSecrets()
	as.Equal("captain.txt: VGhlIENhcHRhaW4=\nstowaway.txt: TGVnYXR0\n", out)
}

func TestLines(t *testing.T) {
	as := assert.New(t)

	f := NewFiles(getTestFiles())

	out := f.Lines("multiline/test.txt")
	as.Len(out, 2)

	as.Equal("bar", out[0])
}

func TestToYaml(t *testing.T) {
	expect := "foo: bar\n"
	v := struct {
		Foo string `json:"foo"`
	}{
		Foo: "bar",
	}

	if got := ToYaml(v); got != expect {
		t.Errorf("Expected %q, got %q", expect, got)
	}
}
