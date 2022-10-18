package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ericuni/errorx"
	"github.com/spf13/cast"
)

// 19 srv_xxx doXxx service/srv_xxx/doXxx.go:84:1
// 1 mock_dbops (*MockXxxRepo).DoXxx mock/dal/dbops/xxx.go:176:1
type Element struct {
	Fun        string // (*MockXxxRepo).DoXxx
	Complexity int    // 1
	File       string // mock/dal/dbops/xxx.go
	Pkg        string // mock_dbops
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
	if strings.Contains(e.File, ".generated.go") {
		return true
	}

	// gomock
	if strings.Contains(e.File, "mock/") {
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
		Fun:        tokens[2],
		Complexity: cast.ToInt(tokens[0]),
		Pkg:        tokens[1],
	}

	location := tokens[3]
	tokens = strings.Split(location, ":")
	if len(tokens) != 3 {
		return nil, errorx.New("location %s format error", location)
	}
	e.File = tokens[0]

	return e, nil
}

type Pair struct {
	Current *Element
	Base    *Element
}

// ignore those functions with complexity <= ignoreComplexity
func merge(current, base []*Element, ignoreComplexity int) []Pair {
	m := make(map[string]*Element, len(base))
	for _, e := range base {
		m[getKey(e)] = e
	}

	var res []Pair
	for _, e := range current {
		if e.Complexity <= ignoreComplexity {
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
	return fmt.Sprintf("%s|%s", e.File, e.Fun)
}
