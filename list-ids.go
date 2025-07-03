package nstore

import (
	"os"
	"path/filepath"
)

func (s *NStorage[T]) IDs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.cache))
	for id := range s.cache {
		ids = append(ids, id)
	}
	if s.loaded {
		return ids
	}

	files, err := os.ReadDir(s.folder)
	if err != nil {
		return ids
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		id := file.Name()[:len(file.Name())-len(".json")]
		if _, exists := s.cache[id]; !exists {
			ids = append(ids, id)
		}
	}
	return ids
}
