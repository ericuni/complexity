package main

import (
	"context"
	"flag"

	"github.com/golang/glog"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "./")
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
	return nil
}
