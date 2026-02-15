package util

import "github.com/google/uuid"

func GenerateUUIDStr() string {
	id := uuid.New().String()
	return id
}

func GenerateUUID() uuid.UUID {
	return uuid.New()
}
