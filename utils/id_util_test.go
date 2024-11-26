package utils

import (
	"fmt"
	"testing"
)

func Test_encodeToShortID(t *testing.T) {
	sf := NewSnowflake(1)
	id := sf.NextID()
	shortID := encodeToShortID(id)

	fmt.Println(shortID)
}
