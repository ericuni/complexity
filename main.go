package main

import (
	"context"
	"flag"

	"github.com/ericuni/errorx"
	"github.com/golang/glog"
)

var (
	flagBase    string
	flagCurrent string
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "./")
	flag.StringVar(&flagBase, "base", "", "base path")
	flag.StringVar(&flagCurrent, "current", "", "current path")
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

	merged := merge(current, base)

	for _, pair := range merged {
		if pair.Base != nil {
			glog.Infof("pkg: %s fun: %s complexity: %d(%d)", pair.Current.Pkg, pair.Current.Fun, pair.Current.Complexity,
				pair.Current.Complexity-pair.Base.Complexity)
		} else {
			glog.Infof("pkg: %s fun: %s complexity: %d(%d)", pair.Current.Pkg, pair.Current.Fun, pair.Current.Complexity,
				pair.Current.Complexity)
		}
	}

	return nil
}
