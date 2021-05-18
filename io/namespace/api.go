package namespace

import (
	"os"
)

func New(name string, config ...Config) *Namespace {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c := Config{
		Path:      wd,
		ShardSize: 2 * 1024,
	}
	if len(config) > 0 {
		c = config[0]
	}

	ns := &Namespace{
		basePath:  c.Path,
		name:      name,
		shardSize: c.ShardSize,
	}

	ns.init()

	return ns
}
