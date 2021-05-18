package shard

func New(path string) (*Shard, error) {
	s := &Shard{
		path:   path,
		idIdx:  nil,
		handle: nil,
	}

	err := s.init()
	if err != nil {
		return nil, err
	}

	return s, nil
}
