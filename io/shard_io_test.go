package io

import (
	"github.com/Azer0s/athena/io/shard"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewShard(t *testing.T) {
	s, _ := shard.New("foo")

	id := uuid.NewString()

	val := shard.Values{
		"foo":  "bar",
		"x":    "y",
		"test": 16.0,
		"b":    true,
	}
	err := s.Write(id, val)
	if err != nil {
		return
	}

	res, err := s.Read(id)
	if err != nil {
		return
	}

	assert.Equal(t, res.Values, val)

	os.Remove("foo")
}

func TestNewShardDelete(t *testing.T) {
	s, _ := shard.New("foo")

	id1 := uuid.NewString()

	val := shard.Values{
		"foo":  "bar",
		"x":    "y",
		"test": 16.0,
		"b":    true,
	}
	err := s.Write(id1, val)
	if err != nil {
		return
	}

	id2 := uuid.NewString()
	err = s.Write(id2, val)
	if err != nil {
		return
	}

	idToDelete := uuid.NewString()
	err = s.Write(idToDelete, val)
	if err != nil {
		return
	}

	id3 := uuid.NewString()
	err = s.Write(id3, val)
	if err != nil {
		return
	}

	id4 := uuid.NewString()
	err = s.Write(id4, val)
	if err != nil {
		return
	}

	err = s.Delete(idToDelete)

	res, err := s.Read(id1)
	if err != nil {
		return
	}
	assert.Equal(t, res.Values, val)

	res, err = s.Read(id2)
	if err != nil {
		return
	}
	assert.Equal(t, res.Values, val)

	res, err = s.Read(id3)
	if err != nil {
		return
	}
	assert.Equal(t, res.Values, val)

	res, err = s.Read(id4)
	if err != nil {
		return
	}
	assert.Equal(t, res.Values, val)

	os.Remove("foo")
}
