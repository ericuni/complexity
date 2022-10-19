package internal

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

type Item struct {
	Fun        string
	Complexity int
	File       string
	Pkg        string
	Line       int
	Column     int
}

func ParseComplexity(ctx context.Context, path string) ([]*Item, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, errorx.Trace(err)
	}

	lines := strings.Split(string(bs), "\n")

	res := make([]*Item, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}

		e, err := line2item(ctx, line)
		if err != nil {
			return nil, errorx.Trace(err)
		}

		if skipOrNot(e) {
			continue
		}

		res = append(res, e)
	}

	return res, nil
}

func skipOrNot(e *Item) bool {
	if strings.Contains(e.File, ".generated.go") {
		return true
	}

	// gomock
	if strings.Contains(e.File, "mock/") {
		return true
	}

	return false
}

func line2item(ctx context.Context, line string) (*Item, error) {
	// line format: <complexity> <package> <function> <file:line:column>
	// example 1: 19 srv_xxx doXxx service/srv_xxx/doXxx.go:84:1
	// example 2: 1 mock_dbops (*MockXxxRepo).DoXxx mock/dal/dbops/xxx.go:176:1
	tokens := strings.Split(line, " ")
	if len(tokens) != 4 {
		return nil, errorx.New("line %s format error", line)
	}

	item := &Item{
		Fun:        tokens[2],
		Complexity: cast.ToInt(tokens[0]),
		Pkg:        tokens[1],
	}

	location := tokens[3]
	tokens = strings.Split(location, ":")
	if len(tokens) != 3 {
		return nil, errorx.New("location %s format error", location)
	}
	item.File = tokens[0]
	item.Line = cast.ToInt(tokens[1])
	item.Column = cast.ToInt(tokens[2])

	return item, nil
}

type Pair struct {
	Current *Item
	Base    *Item
}

func Merge(base, current []*Item) []Pair {
	m := make(map[string]*Item, len(base))
	for _, item := range base {
		m[getKey(item)] = item
	}

	var res []Pair
	for _, item := range current {
		ori, ok := m[getKey(item)]
		if !ok {
			res = append(res, Pair{Current: item})
			continue
		}

		res = append(res, Pair{
			Current: item,
			Base:    ori,
		})
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Current.Complexity > res[j].Current.Complexity
	})

	return res
}

// we can move a function between files in the same directory(i.e. pkg)
func getKey(e *Item) string {
	return fmt.Sprintf("%s|%s", path.Dir(e.File), e.Fun)
}
