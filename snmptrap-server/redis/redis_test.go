package redis

import (
	"testing"
)

const (
	subject = "test"
	content = "test1234"
)

func Test_WriteMail(t *testing.T) {
	tos := []string{"abc@example.com", "def@example.com"}
	InitRedisConnPool()
	WriteMail(tos, subject, content)
	RedisConnPool.Close()
}
