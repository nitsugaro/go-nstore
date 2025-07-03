package nstore

import "time"

type Metadata struct {
	ID         string    `json:"id"`
	Rev        string    `json:"rev"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func (m *Metadata) GetMetadata() *Metadata {
	return m
}

type IMetadata interface {
	GetMetadata() *Metadata
}
