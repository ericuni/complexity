package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ericuni/complexity/internal"
	"github.com/ericuni/errorx"
	"github.com/golang/glog"
)

var (
	flagBase          string
	flagCurrent       string
	flagMinComplexity int
	flagMaxComplexity int
)

func init() {
	flag.Set("logtostderr", "true")
	flag.StringVar(&flagBase, "base", "./base", "base path")
	flag.StringVar(&flagCurrent, "current", "./current", "current path")
	flag.IntVar(&flagMinComplexity, "min_complexity", 5, "do not display those functions with complexity < min_complexity")
	flag.IntVar(&flagMaxComplexity, "max_complexity", 20, "when there is a function whose complexity > max_complexity, exit with status 1")
	flag.Parse()
}

func main() {
	defer glog.Flush()

	ctx := context.Background()
	if err := run(ctx); err != nil {
		glog.Errorf("run error %+v", err)
		os.Exit(-1)
		return
	}
	glog.Infoln("done")
}

func run(ctx context.Context) error {
	if flagBase == "" || flagCurrent == "" {
		return errorx.New("base or current path empty")
	}

	base, err := internal.ParseComplexity(ctx, flagBase)
	if err != nil {
		return errorx.Trace(err)
	}

	current, err := internal.ParseComplexity(ctx, flagCurrent)
	if err != nil {
		return errorx.Trace(err)
	}

	pairs := internal.Merge(base, current)

	for _, pair := range pairs {
		if pair.Current.Complexity < flagMinComplexity {
			continue
		}

		diff := getComplexity(pair)
		if diff == "" {
			continue
		}

		fmt.Println(pair.Current.File, pair.Current.Fun, diff)
	}

	if len(pairs) > 0 && pairs[0].Current.Complexity > flagMaxComplexity {
		os.Exit(1)
	}

	return nil
}

func getComplexity(pair internal.Pair) string {
	diff := pair.Current.Complexity
	if pair.Base != nil {
		diff = pair.Current.Complexity - pair.Base.Complexity
	}

	if diff == 0 {
		return ""
	}

	if diff > 0 {
		return fmt.Sprintf("%d(+%d)", pair.Current.Complexity, diff)
	}
	return fmt.Sprintf("%d(-%d)", pair.Current.Complexity, -diff)
}
