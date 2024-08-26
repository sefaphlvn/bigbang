package helper

import "github.com/gin-gonic/gin"

func GetFromContext[T any](c *gin.Context, key string) (T, bool) {
	value, exists := c.Get(key)
	if !exists {
		var zero T
		return zero, false
	}

	typedValue, ok := value.(T)
	if !ok {
		var zero T
		return zero, false
	}

	return typedValue, true
}
