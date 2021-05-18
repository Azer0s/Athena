package shard

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Azer0s/athena/marshaller"
	"github.com/Azer0s/athena/util"
	"io"
	"os"
)

type Shard struct {
	path   string
	idIdx  map[string]*DocumentInfo
	handle *os.File
}

func (s *Shard) init() (err error) {
	s.handle, err = os.OpenFile(s.path, os.O_RDWR, os.ModePerm)

	if os.IsNotExist(err) {
		s.handle, err = os.OpenFile(s.path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)

		if err != nil {
			return
		}
	}

	s.idIdx = make(map[string]*DocumentInfo)
	err = s.buildIndices()
	if err != nil {
		return
	}

	return
}

func (s *Shard) buildIndices() error {
	defer func() {
		_, err := s.handle.Seek(0, 2)
		if err != nil {
			//something is seriously wrong if we can't seek the file to end
			panic(err)
		}
	}()

	seek, err := s.handle.Seek(0, 0)
	if err != nil {
		return err
	}

	stat, _ := s.handle.Stat()

	for seek < stat.Size() {
		info := &DocumentInfo{}

		start, err := s.handle.Seek(0, 1)
		if err != nil {
			return err
		}

		info.Start = start

		idBytes, err := util.ReadBytesLen(s.handle)
		if err != nil {
			return err
		}

		contentLen, err := util.ReadLen(s.handle)
		if err != nil {
			return err
		}
		info.Len = contentLen

		offset, _ := s.handle.Seek(0, 1)
		info.Pos = offset

		seek, err = s.handle.Seek(int64(contentLen), 1)
		if err != nil {
			return err
		}

		s.idIdx[string(idBytes)] = info
	}

	return nil
}

func (s *Shard) Write(id string, values map[string]interface{}) error {
	//flush the file before doing anything
	err := s.handle.Sync()
	if err != nil {
		return err
	}

	defer func() {
		_, err := s.handle.Seek(0, 2)
		if err != nil {
			//something is seriously wrong if we can't seek the file to end
			panic(err)
		}
	}()

	start, err := s.handle.Seek(0, 1)
	if err != nil {
		return err
	}

	buff := &bytes.Buffer{}

	//write id to buffer
	util.StrToLenBytes(id, buff)

	valBuff := &bytes.Buffer{}

	for k, v := range values {
		//write key to buffer
		util.StrToLenBytes(k, valBuff)

		vBuffer := &bytes.Buffer{}
		err := marshaller.Marshal(v, vBuffer)
		if err != nil {
			return err
		}

		valBuff.Write(vBuffer.Bytes())
	}

	util.BytesToLenBytes(valBuff.Bytes(), buff)

	_, err = s.handle.Write(buff.Bytes())
	if err != nil {
		return err
	}

	offset, err := s.handle.Seek(int64(-valBuff.Len()), 1)
	if err != nil {
		return err
	}

	info := &DocumentInfo{
		Start: start,
		Pos:   offset,
		Len:   uint32(valBuff.Len()),
	}

	s.idIdx[id] = info

	return nil
}

func (s *Shard) Read(id string) (*Document, error) {
	//flush the file before doing anything
	err := s.handle.Sync()
	if err != nil {
		return nil, err
	}

	defer func() {
		_, err := s.handle.Seek(0, 2)
		if err != nil {
			//something is seriously wrong if we can't seek the file to end
			panic(err)
		}
	}()

	info, ok := s.idIdx[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("document with id %s not found in shard", id))
	}

	_, err = s.handle.Seek(info.Pos, 0)
	if err != nil {
		return nil, err
	}

	contentBytes := make([]byte, info.Len)
	_, err = s.handle.Read(contentBytes)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(contentBytes)

	res := make(map[string]interface{})

	for {
		keyBytes, err := util.ReadBytesLen(buff)
		if errors.Is(io.EOF, err) {
			return &Document{
				Id:     id,
				Values: res,
			}, nil
		}

		if err != nil {
			return nil, err
		}

		typeBytes, err := util.ReadBytesLen(buff)
		if err != nil {
			return nil, err
		}

		valBytes, err := util.ReadBytesLen(buff)
		if err != nil {
			return nil, err
		}

		val, err := marshaller.UnmarshalWithoutLen(typeBytes, valBytes)
		if err != nil {
			return nil, err
		}

		res[string(keyBytes)] = val
	}
}

func (s *Shard) Delete(id string) error {
	//flush the file before doing anything
	err := s.handle.Sync()
	if err != nil {
		return err
	}

	defer func() {
		//rebuild indices
		err = s.buildIndices()
		if err != nil {
			panic(err)
		}
	}()

	info, ok := s.idIdx[id]
	if !ok {
		return errors.New(fmt.Sprintf("document with id %s not found in shard", id))
	}

	stat, _ := s.handle.Stat()

	//seek to end of document
	_, err = s.handle.Seek(info.Pos+int64(info.Len), 0)
	if err != nil {
		return err
	}

	buff := make([]byte, stat.Size()-(info.Pos+int64(info.Len)))
	_, err = s.handle.Read(buff)
	if err != nil {
		return err
	}

	//seek to the beginning of the document
	_, err = s.handle.Seek(info.Start, 0)
	if err != nil {
		return err
	}

	//overwrite bytes
	_, err = s.handle.Write(buff)
	if err != nil {
		return err
	}

	targetSize := stat.Size() - (info.Pos - info.Start) - int64(info.Len)
	err = s.handle.Truncate(targetSize)

	delete(s.idIdx, id)

	return err
}

func (s *Shard) Documents() []string {
	ids := make([]string, 0)
	for k := range s.idIdx {
		ids = append(ids, k)
	}

	return ids
}
