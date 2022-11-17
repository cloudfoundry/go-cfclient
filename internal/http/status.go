package http

import "net/http"

func IsResponseRedirect(statusCode int) bool {
	switch statusCode {
	case http.StatusTemporaryRedirect, http.StatusPermanentRedirect, http.StatusMovedPermanently, http.StatusFound, http.StatusSeeOther:
		return true
	}
	return false
}
