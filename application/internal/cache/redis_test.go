package cache

import (
	"testing"
)

// docker run --rm --name testredis -d -p 6379:6379 redis
func TestRedis(t *testing.T) {
	client := NewRedisCache("localhost:6379")

	want := "101"
	err := client.Set("TestRedis", want)
	if err != nil {
		t.Errorf("redis key set failed, err: %s", err.Error())
	}

	got, err := client.Get("TestRedis")
	if want != got {
		t.Errorf("failed: got %s, want %s", got, want)
	}
}
