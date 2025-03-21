package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateId() string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var id strings.Builder

	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Generate first part of ID (6 characters)
	for i := 0; i < 6; i++ {
		id.WriteByte(characters[rand.Intn(len(characters))])
	}

	id.WriteByte('-')

	// Generate second part of ID (8 characters)
	for i := 0; i < 8; i++ {
		id.WriteByte(characters[rand.Intn(len(characters))])
	}

	id.WriteByte('-')

	// Generate third part of ID (16 characters)
	for i := 0; i < 16; i++ {
		id.WriteByte(characters[rand.Intn(len(characters))])
	}

	return id.String()
}

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func TimePtr(t time.Time) *time.Time {
	return &t
}
