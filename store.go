package nstore

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

type NStorage[T IMetadata] struct {
	folder string
	cache  map[string]T
	mu     sync.RWMutex
	loaded bool
}

func New[T IMetadata](folder string) (*NStorage[T], error) {
	if err := os.MkdirAll(folder, 0755); err != nil {
		return nil, err
	}
	return &NStorage[T]{
		folder: folder,
		cache:  make(map[string]T),
	}, nil
}

func (s *NStorage[T]) Save(obj T) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	meta := obj.GetMetadata()
	if meta == nil {
		val := reflect.ValueOf(obj)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		field := val.FieldByName("Metadata")
		if field.IsValid() && field.CanSet() && field.IsNil() {
			newMeta := &Metadata{}
			field.Set(reflect.ValueOf(newMeta))
			meta = newMeta
		}
	}

	if meta.ID == "" {
		meta.ID = uuid.NewString()
		meta.CreatedAt = time.Now()
	}
	meta.ModifiedAt = time.Now()
	meta.Rev = uuid.NewString()

	s.cache[meta.ID] = obj

	path := filepath.Join(s.folder, meta.ID+".json")
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (s *NStorage[T]) Load(id string) (T, error) {
	s.mu.RLock()
	obj, ok := s.cache[id]
	loaded := s.loaded
	s.mu.RUnlock()

	if ok {
		return obj, nil
	}

	if loaded {
		var zero T
		return zero, os.ErrNotExist
	}

	path := filepath.Join(s.folder, id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		var zero T
		return zero, err
	}

	var loadedObj T
	if err := json.Unmarshal(data, &loadedObj); err != nil {
		var zero T
		return zero, err
	}

	s.mu.Lock()
	s.cache[id] = loadedObj
	s.mu.Unlock()

	return loadedObj, nil
}

func (s *NStorage[T]) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.cache, id)
	path := filepath.Join(s.folder, id+".json")
	return os.Remove(path)
}

func (s *NStorage[T]) ListOfCache() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]T, 0, len(s.cache))
	for _, v := range s.cache {
		result = append(result, v)
	}
	return result
}
