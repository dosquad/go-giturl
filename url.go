package giturl

import (
	"net/url"
	"strings"
)

type URL struct {
	*url.URL

	PathUser     *url.Userinfo
	PathRelative string
}

func (u *URL) Slug() string {
	return strings.TrimPrefix(
		strings.TrimSuffix(
			strings.TrimSuffix(u.PathRelative, "/"),
			".git",
		),
		"/",
	)
}

func (u *URL) String() string {
	// if u.PathRelative != "" {
	// 	u2 := *u.URL
	// 	u2.Path = u.PathRelative

	// 	// if u.PathUser.String() != "" {
	// 	// 	u.User = url.User(u.hostUser)
	// 	// } else if strings.HasPrefix(u.PathRelative, "/~") || strings.HasPrefix(u.PathRelative, "~") {
	// 	// 	u.User = nil
	// 	// }

	// 	return u2.String()
	// }

	if u.Scheme == "" {
		return strings.TrimPrefix(u.URL.String(), "//")
	}

	return u.URL.String()
}
