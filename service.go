// Copyright 2020 morgine.com. All rights reserved.

package cfg

import (
	"github.com/morgine/service"
)

type Service struct {
	storageService StorageService
	self           service.Provider
}

func (s *Service) New(ctn *service.Container) (value interface{}, err error) {
	var store Storage
	store, err = s.storageService.Get(ctn)
	if err != nil {
		return nil, err
	}
	var data []byte
	data, err = store.GetTomlData()
	if err != nil {
		return nil, err
	}
	var env Env
	env, err = Decode(data)
	if err != nil {
		return nil, err
	}
	return env, nil
}

// Get singleton data
func (s *Service) Get(ctn *service.Container) (envs Env, err error) {
	es, er := ctn.Get(&s.self)
	if er != nil {
		return nil, er
	} else {
		return es.(Env), nil
	}
}

func NewService(storageService StorageService) *Service {
	s := &Service{storageService: storageService}
	s.self = s
	return s
}
