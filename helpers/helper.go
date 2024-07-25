package helpers

import (
	"net/url"
	"strings"
)

func ResolveUrl(ur string) string {
	urs := strings.Split(ur, "/")
	if len(urs) == 0 {
		return ur
	}
	urs[len(urs)-1] = url.QueryEscape(urs[len(urs)-1])
	return strings.Join(urs, "/")
}
