// Copyright 2020 morgine.com. All rights reserved.

package cfg

import (
	"github.com/morgine/service"
	"io/ioutil"
)

type Storage interface {
	GetTomlData() ([]byte, error)
}

type StorageService interface {
	service.Provider
	Get(ctn *service.Container) (Storage, error)
}

type MemoryStorage string

type MemoryStorageService struct {
	storage MemoryStorage
	self    service.Provider
}

func NewMemoryStorageService(configData string) *MemoryStorageService {
	s := &MemoryStorageService{
		storage: MemoryStorage(configData),
		self:    nil,
	}
	s.self = s
	return s
}

func (s MemoryStorage) GetTomlData() ([]byte, error) {
	return []byte(s), nil
}

func (m MemoryStorageService) New(ctn *service.Container) (value interface{}, err error) {
	return Storage(m.storage), nil
}

func (m MemoryStorageService) Get(ctn *service.Container) (Storage, error) {
	storage, _ := ctn.Get(&m.self)
	return storage.(Storage), nil
}

type FileStorage string

type FileStorageService struct {
	storage FileStorage
	self    service.Provider
}

func NewFileStorageService(filename string) *FileStorageService {
	s := &FileStorageService{
		storage: FileStorage(filename),
		self:    nil,
	}
	s.self = s
	return s
}

func (f FileStorage) GetTomlData() ([]byte, error) {
	return ioutil.ReadFile(string(f))
}

func (f FileStorageService) New(ctn *service.Container) (value interface{}, err error) {
	return Storage(f.storage), nil
}

func (f FileStorageService) Get(ctn *service.Container) (Storage, error) {
	storage, _ := ctn.Get(&f.self)
	return storage.(Storage), nil
}
