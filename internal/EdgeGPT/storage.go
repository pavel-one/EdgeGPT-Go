package EdgeGPT

import (
	"EdgeGPT-Go/config"
	"errors"
	"fmt"
	"time"
)

type Storage map[string]*GPT

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) GetOrSet(key string) (*GPT, error) {
	var gpt *GPT

	gpt, err := s.Get(key)
	if err == nil {
		return gpt, nil
	}

	conf, err := config.NewGpt()
	if err != nil {
		return nil, fmt.Errorf("didn't create GPT config: %s", err)
	}

	gpt, err = NewGPT(conf)
	if err != nil {
		return nil, fmt.Errorf("didn't init GPT service: %s", err)
	}

	s.Add(gpt, key)

	return gpt, nil
}

func (s *Storage) Add(gpt *GPT, key string) {
	(*s)[key] = gpt
}

func (s *Storage) Get(key string) (*GPT, error) {
	v, ok := (*s)[key]
	if !ok {
		return nil, errors.New("not set session")
	}

	if time.Now().After(v.ExpiredAt) {
		if err := s.Remove(key); err != nil {
			return nil, err
		}
		return nil, errors.New("session is expired")
	}

	return v, nil
}

func (s *Storage) Remove(key string) error {
	so := *s
	_, ok := so[key]
	if !ok {
		return errors.New("not set session")
	}

	delete(so, key)

	return nil
}
