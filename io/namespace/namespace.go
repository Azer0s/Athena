package namespace

import (
	"github.com/Azer0s/athena/io/shard"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"sync"
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

	closingMu *sync.RWMutex
	closing   bool

	closeWriter  bool
	writeQueueMu *sync.Mutex
	writeQueue   []*shard.Document

	openOps *sync.WaitGroup
}

func (n *Namespace) readLoop() {
	var doc *shard.Document

	for len(n.writeQueue) == 0 {
		//chill out for a while
	}

	n.writeQueueMu.Lock()
	l := len(n.writeQueue)
	doc, n.writeQueue = n.writeQueue[0], n.writeQueue[1:l]
	n.writeQueueMu.Unlock()

	var shardToWrite *shard.Shard = nil

	for _, s := range n.shards {
		if s.Size() <= n.shardSize {
			shardToWrite = s
			break
		}
	}

	if shardToWrite == nil {
		var err error
		shardToWrite, err = n.createShard()

		if err != nil {
			panic(err)
		}
	}

	err := shardToWrite.Write(doc.Id, doc.Values)
	if err != nil {
		panic(err)
	}

	n.openOps.Done()
}

func (n *Namespace) setupWriter() {
	go func() {
		for !n.closeWriter {
			n.readLoop()
		}
	}()
}

func (n *Namespace) init() error {
	n.openOps = &sync.WaitGroup{}

	n.closingMu = &sync.RWMutex{}
	n.writeQueueMu = &sync.Mutex{}

	n.path = filepath.Join(n.basePath, n.name)
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

	shardDocument, err := n.metaShard.Read("shards")
	if err != nil {
		err := n.metaShard.Write("shards", map[string]interface{}{
			"value": []string{},
		})
		if err != nil {
			return err
		}

		err = n.metaShard.Write("max_shardsize", map[string]interface{}{
			"value": n.shardSize,
		})
		if err != nil {
			return err
		}

		_, err = n.createShard()
		if err != nil {
			return err
		}

		shardDocument, err = n.metaShard.Read("shards")
		if err != nil {
			return err
		}
	}

	shardSizeDocument, err := n.metaShard.Read("max_shardsize")
	if err != nil {
		return err
	}
	n.shardSize = shardSizeDocument.Get("value").(int)

	shards := shardDocument.Get("value").([]interface{})

	for _, s := range shards {
		openedShard, err := shard.New(filepath.Join(n.path, s.(string)))
		if err != nil {
			return err
		}

		n.shards = append(n.shards, openedShard)
	}

	n.idIdx = make(map[string]*shard.Shard)

	for _, s := range n.shards {
		for _, doc := range s.Documents() {
			n.idIdx[doc] = s
		}
	}

	n.setupWriter()

	return nil
}

func (n *Namespace) createShard() (*shard.Shard, error) {
	name := uuid.NewString() + ".ath"
	s, err := shard.New(filepath.Join(n.path, name))
	if err != nil {
		return nil, err
	}

	n.shards = append(n.shards, s)

	document, err := n.metaShard.Read("shards")
	if err != nil {
		panic(err)
	}

	shards := document.Get("value").([]interface{})
	shards = append(shards, name)

	err = n.metaShard.Delete("shards")
	if err != nil {
		panic(err)
	}

	err = n.metaShard.Write("shards", map[string]interface{}{
		"value": shards,
	})
	if err != nil {
		panic(err)
	}

	return s, nil
}
