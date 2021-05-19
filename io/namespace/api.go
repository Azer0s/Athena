package namespace

import (
	"errors"
	"github.com/Azer0s/athena/io/shard"
	"os"
)

func New(name string, config ...Config) (*Namespace, error) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	c := Config{
		Path:      wd,
		ShardSize: 2 * 1024,
	}
	if len(config) > 0 {
		if config[0].Path != "" {
			c.Path = config[0].Path
		}

		if config[0].ShardSize != 0 {
			c.ShardSize = config[0].ShardSize
		}
	}

	ns := &Namespace{
		basePath:  c.Path,
		name:      name,
		shardSize: c.ShardSize,
	}

	err = ns.init()
	if err != nil {
		return nil, err
	}

	return ns, nil
}

func (n *Namespace) Put(id string, values map[string]interface{}) error {
	n.closingMu.RLock()
	defer n.closingMu.RUnlock()

	if n.closing {
		return errors.New("namespace is closing, can't start new write process")
	}

	n.openOps.Add(1)

	n.writeQueueMu.Lock()
	defer n.writeQueueMu.Unlock()

	n.writeQueue = append(n.writeQueue, &shard.Document{
		Id:     id,
		Values: values,
	})

	//TODO: put document into in-memory idx
	//TODO: idx document id to in-memory store

	return nil
}

func (n *Namespace) Name() string {
	return n.name
}

func (n *Namespace) ShardSize() int {
	return n.shardSize
}

func (n *Namespace) Close() error {
	n.closingMu.Lock()
	n.closing = true
	n.closingMu.Unlock()

	//now, no one should be able to enqueue something anymore, so we can wait until the queue len are 0
	//and close the goroutines

	n.openOps.Wait()

	n.closeWriter = true

	return nil
}

func (n *Namespace) Path() string {
	return n.path
}
