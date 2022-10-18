package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ericuni/complexity/internal"
	"github.com/ericuni/errorx"
	"github.com/golang/glog"
)

var (
	flagBase             string
	flagCurrent          string
	flagIgnoreComplexity int
	flagMaxComplexity    int
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.StringVar(&flagBase, "base", "./base", "base path")
	flag.StringVar(&flagCurrent, "current", "./current", "current path")
	flag.IntVar(&flagIgnoreComplexity, "ignore_complexity", 5, "do not display those functions with complexity <= ignore_complexity")
	flag.IntVar(&flagMaxComplexity, "max_complexity", 20, "when there is a function whose complexity >= max_complexity, exit with status 1")
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

	pairs := internal.CompareAndMerge(current, base, flagIgnoreComplexity)

	var buffer bytes.Buffer
	for _, pair := range pairs {
		buffer.Reset()
		buffer.WriteString(fmt.Sprintf("%s %s %d", pair.Current.File, pair.Current.Fun, pair.Current.Complexity))

		diff := pair.Current.Complexity
		if pair.Base != nil {
			diff = pair.Current.Complexity - pair.Base.Complexity
		}

		if diff > 0 {
			buffer.WriteString(fmt.Sprintf("(+%d)", diff))
		} else {
			buffer.WriteString(fmt.Sprintf("(-%d)", -diff))
		}
		fmt.Println(buffer.String())
	}

	if len(pairs) > 0 && pairs[0].Current.Complexity >= flagMaxComplexity {
		os.Exit(1)
	}

	return nil
}
