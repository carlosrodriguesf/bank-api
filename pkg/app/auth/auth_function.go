package auth

import "fmt"

func getSessionCacheKey(token string) string {
	return fmt.Sprintf(cacheKeySession, token)
}
