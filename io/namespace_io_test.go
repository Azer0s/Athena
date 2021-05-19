package io

import (
	"github.com/Azer0s/athena/io/namespace"
	"github.com/Azer0s/athena/io/shard"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestNewNamespace(t *testing.T) {
	n, _ := namespace.New("test", namespace.Config{
		ShardSize: 2048,
	})

	defer func(n *namespace.Namespace) {
		err := n.Close()
		if err != nil {
			panic(err)
		}

		err = os.RemoveAll(n.Path())
		if err != nil {
			panic(err)
		}
	}(n)

	for i := 0; i < 100; i++ {
		err := n.Put(uuid.NewString(), shard.Values{
			"foo":  "bar",
			"x":    "y",
			"test": 16.0,
			"b":    true,
		})
		if err != nil {
			return
		}
	}
}
