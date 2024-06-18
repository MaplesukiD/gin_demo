package utils

import "github.com/google/uuid"

const NAME_PREFIX = "user_"

// RandomString 生成随机用户名
func RandomString() string {
	return (NAME_PREFIX + uuid.New().String())[:10]
}
