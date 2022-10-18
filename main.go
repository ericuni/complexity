package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"

	"github.com/ericuni/errorx"
	"github.com/golang/glog"
)

var (
	flagBase             string
	flagCurrent          string
	flagIgnoreComplexity int
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.StringVar(&flagBase, "base", "./base", "base path")
	flag.StringVar(&flagCurrent, "current", "./current", "current path")
	flag.IntVar(&flagIgnoreComplexity, "ignore_complexity", 5, "ignore those functions with complexity <=")
	flag.Parse()
}

func main() {
	defer glog.Flush()

	ctx := context.Background()
	if err := run(ctx); err != nil {
		glog.Errorf("run error %+v", err)
		return
	}
	glog.Infoln("done")
}

func run(ctx context.Context) error {
	if flagBase == "" || flagCurrent == "" {
		return errorx.New("base or current path empty")
	}

	base, err := parseCyclomatic(ctx, flagBase)
	if err != nil {
		return errorx.Trace(err)
	}

	current, err := parseCyclomatic(ctx, flagCurrent)
	if err != nil {
		return errorx.Trace(err)
	}

	merged := merge(current, base, flagIgnoreComplexity)

	var buffer bytes.Buffer
	for _, pair := range merged {
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

	return nil
}
