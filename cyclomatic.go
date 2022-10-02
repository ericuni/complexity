package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/ericuni/errorx"
	"github.com/spf13/cast"
)

// 19 srv_material completeSaveMaterialReq service/srv_material/save_material.go:84:1
type Element struct {
	Pkg        string // service/srv_material
	Fun        string // completeSaveMaterialReq
	Complexity int    // 19
	Location   string // service/srv_material/save_material.go:84:1
}

func parseCyclomatic(ctx context.Context, path string) ([]*Element, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.Trace(err)
	}

	lines := strings.Split(string(bs), "\n")

	res := make([]*Element, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}

		e, err := line2Element(ctx, line)
		if err != nil {
			return nil, errorx.Trace(err)
		}

		if ifSkip(e) {
			continue
		}

		res = append(res, e)
	}

	return res, nil
}

func ifSkip(e *Element) bool {
	if strings.Contains(e.Location, ".generated.go") {
		return true
	}

	// gomock
	if strings.Contains(e.Pkg, "mock/") {
		return true
	}

	return false
}

func line2Element(ctx context.Context, line string) (*Element, error) {
	tokens := strings.Split(line, " ")
	if len(tokens) != 4 {
		return nil, errorx.New("line %s format error", line)
	}

	e := &Element{
		Pkg:        path.Dir(tokens[3]),
		Fun:        tokens[2],
		Complexity: cast.ToInt(tokens[0]),
		Location:   tokens[3],
	}

	return e, nil
}

type Pair struct {
	Current *Element
	Base    *Element
}

func merge(current, base []*Element) []Pair {
	m := make(map[string]*Element, len(base))
	for _, e := range base {
		m[getKey(e)] = e
	}

	var res []Pair
	for _, e := range current {
		if e.Complexity <= 5 {
			continue
		}

		ori, ok := m[getKey(e)]
		if !ok {
			res = append(res, Pair{Current: e})
			continue
		}

		if e.Complexity == ori.Complexity {
			continue
		}

		res = append(res, Pair{
			Current: e,
			Base:    ori,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Current.Complexity > res[j].Current.Complexity
	})

	return res
}

func getKey(e *Element) string {
	return fmt.Sprintf("%s|%s", e.Pkg, e.Fun)
}
