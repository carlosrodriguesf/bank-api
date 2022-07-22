package auth

import (
	"net/http"
	"strings"
)

func getTokenFromRequest(req *http.Request) string {
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return ""
	}

	splitted := strings.Split(authorization, " ")
	if len(splitted) != 2 {
		return ""
	}
	if strings.ToLower(splitted[0]) != "bearer" {
		return ""
	}

	return splitted[1]
}
