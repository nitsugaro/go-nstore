package nstore

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

/* Checks if exists on cache or disk if it has not yet been loaded */
func (s *NStorage[T]) Exists(id string) bool {
	s.mu.RLock()
	_, ok := s.cache[id]
	loaded := s.loaded
	s.mu.RUnlock()

	if ok {
		return true
	}
	if loaded {
		return false
	}

	path := filepath.Join(s.folder, id+".json")
	_, err := os.Stat(path)
	return err == nil
}

func (s *NStorage[T]) LoadFromDisk() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache = make(map[string]T)
	files, err := os.ReadDir(s.folder)
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.folder, file.Name()))
		if err != nil {
			return err
		}
		var obj T
		if err := json.Unmarshal(data, &obj); err != nil {
			return err
		}
		id := obj.GetMetadata().ID
		if id == "" {
			return errors.New("object without ID in file " + file.Name())
		}
		s.cache[id] = obj
	}
	s.loaded = true
	return nil
}
