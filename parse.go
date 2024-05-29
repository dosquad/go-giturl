package giturl

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

var ErrUnknownSyntax = errors.New("unknown URL syntax")

var (
	gitURIRegex    = regexp.MustCompile(`(?P<user>\S+)@(?P<host>[\S\.]+):(?P<path>\S+)`)
	transportRegex = regexp.MustCompile(`^(?P<scheme>ssh|git|http(s|)|ftp(s|))://`)

	tildeUserRegex = regexp.MustCompile(
		`^((?P<scheme>ssh|git)://|)((?P<user>\S+)@|)(?P<host>\S+)(:(?P<port>\d+)|)[/:]~(?P<pathuser>[^/]+)/(?P<path>.*)`,
	)
)

func Parse(rawURL string) (*URL, error) {
	if tildeUserRegex.MatchString(rawURL) {
		return parseTildeURI(rawURL)
	}

	if transportRegex.MatchString(rawURL) {
		return parseURI(rawURL)
	}

	if gitURIRegex.MatchString(rawURL) {
		return parseGitURI(rawURL)
	}

	if strings.HasPrefix(rawURL, "/") || strings.HasPrefix(rawURL, "file://") {
		return parseURI(rawURL)
	}

	return nil, ErrUnknownSyntax
}

func parseURI(rawURL string) (*URL, error) {
	u, err := url.Parse(rawURL)
	relPath := ""
	if u != nil {
		relPath = u.Path
	}
	return &URL{
		URL:          u,
		PathUser:     &url.Userinfo{},
		PathRelative: relPath,
	}, err
}

func parseGitURI(rawURL string) (*URL, error) {
	out := gitURIRegex.FindStringSubmatch(rawURL)
	if len(out) == 4 { //nolint:mnd // matches in regexp
		return &URL{
			URL: &url.URL{
				Scheme: "git+ssh",
				User:   url.User(out[1]),
				Host:   out[2],
				Path:   out[3],
			},
			PathUser:     &url.Userinfo{},
			PathRelative: out[3],
		}, nil
	}
	return nil, ErrUnknownSyntax
}

func parseTildeURIBase(rawURL string) (*URL, error) {
	u, err := parseURI(rawURL)
	if err != nil {
		return nil, err
	}

	if u.Host == "" {
		return nil, ErrUnknownSyntax
	}

	userPath := strings.TrimPrefix(u.Path, "/")
	if after, ok := strings.CutPrefix(userPath, "~"); ok {
		sp := strings.SplitN(after, "/", 2) //nolint:mnd // tokenisation
		u.PathUser = url.User(sp[0])
		// u.originalPath = u.Path
		u.PathRelative = sp[1]
	}

	return u, nil
}

func parseTildeURI(rawURL string) (*URL, error) {
	u, err := parseTildeURIBase(rawURL)
	if err == nil {
		return u, nil
	}

	m := tildeUserRegex.FindStringSubmatch(rawURL)
	u = &URL{
		URL:      &url.URL{},
		PathUser: &url.Userinfo{},
	}
	for idx, name := range tildeUserRegex.SubexpNames() {
		if idx == 0 {
			continue
		}
		if m[idx] == "" {
			continue
		}
		switch name {
		case "scheme":
			u.URL.Scheme = m[idx]
		case "pathuser":
			u.PathUser = url.User(m[idx])
			if u.URL.Path != "" {
				u.URL.Path = "~" + m[idx] + "/" + u.URL.Path
			}
		case "user":
			u.URL.User = url.User(m[idx])
		case "path":
			u.PathRelative = m[idx]
			if u.PathUser.String() != "" {
				u.URL.Path = "~" + u.PathUser.Username() + "/" + m[idx]
				continue
			}

			u.URL.Path = m[idx]
		case "host":
			if len(u.URL.Host) > 0 {
				u.URL.Host = m[idx] + u.URL.Host
				continue
			}

			u.URL.Host = m[idx]
		case "port":
			u.URL.Host += ":" + m[idx]
		}
	}

	return u, nil
}
