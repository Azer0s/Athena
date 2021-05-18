package namespace

import (
	"github.com/Azer0s/athena/io/shard"
	"os"
	"path/filepath"
)

type Config struct {
	Path      string
	ShardSize int
}

type Namespace struct {
	basePath  string
	name      string
	path      string
	shardSize int
	shards    []*shard.Shard
	idIdx     map[string]*shard.Shard
	metaShard *shard.Shard
}

func (n *Namespace) init() error {
	n.path = filepath.Join(n.path, n.name)
	_, err := os.Open(n.path)

	if os.IsNotExist(err) {
		err := os.Mkdir(n.path, os.ModePerm)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	n.metaShard, err = shard.New(filepath.Join(n.path, "meta.ath"))
	if err != nil {
		return err
	}

	document, err := n.metaShard.Read("shards")
	if err != nil {
		err := n.metaShard.Write("shard", map[string]interface{}{
			"value": []string{},
		})
		if err != nil {
			return err
		}

		document, err = n.metaShard.Read("shards")
		if err != nil {
			return err
		}
	}

	for _, s := range document.Get("shards").([]string) {
		shrd, err := shard.New(s)
		if err != nil {
			return err
		}

		n.shards = append(n.shards, shrd)
	}

	n.idIdx = make(map[string]*shard.Shard)

	for _, s := range n.shards {
		for _, doc := range s.Documents() {
			n.idIdx[doc] = s
		}
	}

	return nil
}

func (n *Namespace) Name() string {
	return n.name
}

func (n *Namespace) ShardSize() int {
	return n.shardSize
}
