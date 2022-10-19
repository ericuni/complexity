package internal

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_line2item(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	tds := []struct {
		name string
		line string
		want *Item
	}{
		{
			name: "function",
			line: "5 srv_material createMaterial service/srv_material/save_material.go:153:1",
			want: &Item{
				Fun:        "createMaterial",
				Complexity: 5,
				File:       "service/srv_material/save_material.go",
				Pkg:        "srv_material",
				Line:       153,
				Column:     1,
			},
		},
		{
			name: "method",
			line: "6 dbops (*_MaterialRepoStruct).Update dal/dbops/material.generated.go:42:1",
			want: &Item{
				Fun:        "(*_MaterialRepoStruct).Update",
				Complexity: 6,
				File:       "dal/dbops/material.generated.go",
				Pkg:        "dbops",
				Line:       42,
				Column:     1,
			},
		},
	}

	for _, td := range tds {
		t.Run(td.name, func(t *testing.T) {
			item, err := line2item(ctx, td.line)
			assert.Nil(err)
			assert.Equal(td.want, item)
		})
	}
}

func Test_skipOrNot(t *testing.T) {
	assert := assert.New(t)

	tds := []struct {
		name string
		item *Item
		want bool
	}{
		{
			name: "normal",
			item: &Item{File: "service/srv_material/save_material.go"},
			want: false,
		},
		{
			name: "generated",
			item: &Item{File: "dal/ddl/material.generated.go"},
			want: true,
		},
		{
			name: "gomock",
			item: &Item{File: "mock/dal/dbops/material.go"},
			want: true,
		},
	}

	for _, td := range tds {
		t.Run(td.name, func(t *testing.T) {
			skip := skipOrNot(td.item)
			assert.Equal(td.want, skip)
		})
	}
}

func TestParseComplexity(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	// prepare test data
	var buffer bytes.Buffer
	buffer.WriteString("5 srv_material createMaterial service/srv_material/save_material.go:153:1\n")
	buffer.WriteString("4 ddl (*Material).ToMaterialData dal/ddl/material.go:46:1")
	name := "test_parse.txt"
	err := os.WriteFile(name, buffer.Bytes(), 0644)
	assert.Nil(err)

	items, err := ParseComplexity(ctx, name)
	assert.Nil(err, "%+v", err)

	expected := []*Item{
		{
			Fun:        "createMaterial",
			Complexity: 5,
			File:       "service/srv_material/save_material.go",
			Pkg:        "srv_material",
			Line:       153,
			Column:     1,
		},
		{
			Fun:        "(*Material).ToMaterialData",
			Complexity: 4,
			File:       "dal/ddl/material.go",
			Pkg:        "ddl",
			Line:       46,
			Column:     1,
		},
	}

	assert.Equal(expected, items)

	err = os.Remove(name)
	assert.Nil(err)
}

func TestMerge(t *testing.T) {
	assert := assert.New(t)

	base := []*Item{
		{
			Fun:        "foo",
			File:       "pkg1/f1.go",
			Complexity: 6,
		},
		{
			Fun:        "bar",
			File:       "pkg2/f2.go",
			Complexity: 10,
		},
	}

	current := []*Item{
		{
			Fun:        "foo", // move from f1.go to f2.go in the same pkg
			File:       "pkg1/f2.go",
			Complexity: 8,
		},
		{
			Fun:        "new", // add a new function
			File:       "pkg1/f1.go",
			Complexity: 7,
		},
	}

	merged := Merge(base, current)

	want := []Pair{
		{
			Current: &Item{
				Fun:        "foo",
				File:       "pkg1/f2.go",
				Complexity: 8,
			},
			Base: &Item{
				Fun:        "foo",
				File:       "pkg1/f1.go",
				Complexity: 6,
			},
		},
		{
			Current: &Item{
				Fun:        "new",
				File:       "pkg1/f1.go",
				Complexity: 7,
			},
		},
	}

	assert.Equal(want, merged)
}
