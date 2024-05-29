package giturl_test

import (
	"testing"

	"github.com/dosquad/go-giturl"
)

//nolint:lll,gocognit // Test array
func TestGitURLs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		rawURL                                                       string
		scheme, user, pathuser, host, port, path, pathRelative, slug string
	}{
		// ssh://[<user>@]<host>[:<port>]/<path-to-git-repo>
		{"ssh://foo@githost.example:1234/path/to/git/repo", "ssh", "foo", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"ssh://githost.example:1234/path/to/git/repo", "ssh", "", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"ssh://foo@githost.example/path/to/git/repo", "ssh", "foo", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"ssh://githost.example/path/to/git/repo", "ssh", "", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},

		// git://<host>[:<port>]/<path-to-git-repo>
		{"git://githost.example:1234/path/to/git/repo", "git", "", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"git://githost.example/path/to/git/repo", "git", "", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},

		// http[s]://<host>[:<port>]/<path-to-git-repo>
		{"http://githost.example:1234/path/to/git/repo", "http", "", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"https://githost.example:1234/path/to/git/repo", "https", "", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"http://githost.example/path/to/git/repo", "http", "", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"https://githost.example/path/to/git/repo", "https", "", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},

		// ftp[s]://<host>[:<port>]/<path-to-git-repo>
		{"ftp://githost.example:1234/path/to/git/repo", "ftp", "", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"ftps://githost.example:1234/path/to/git/repo", "ftps", "", "", "githost.example", "1234", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"ftp://githost.example/path/to/git/repo", "ftp", "", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},
		{"ftps://githost.example/path/to/git/repo", "ftps", "", "", "githost.example", "", "/path/to/git/repo", "/path/to/git/repo", "path/to/git/repo"},

		// [<user>@]<host>:/<path-to-git-repo>
		{"git@github.com:dosquad/go-giturl.git", "git+ssh", "git", "", "github.com", "", "dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},

		// ssh://[<user>@]<host>[:<port>]/~<user>/<path-to-git-repo>
		{"ssh://github.com:1111/~foo/dosquad/go-giturl.git", "ssh", "", "foo", "github.com", "1111", "/~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},
		{"ssh://github.com/~foo/dosquad/go-giturl.git", "ssh", "", "foo", "github.com", "", "/~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},
		{"ssh://bar@github.com/~foo/dosquad/go-giturl.git", "ssh", "bar", "foo", "github.com", "", "/~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},

		// git://<host>[:<port>]/~<user>/<path-to-git-repo>
		{"git://github.com:1111/~foo/dosquad/go-giturl.git", "git", "", "foo", "github.com", "1111", "/~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},
		{"git://github.com/~foo/dosquad/go-giturl.git", "git", "", "foo", "github.com", "", "/~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},

		// [<user>@]<host>:~<user>/<path-to-git-repo>
		{"github.com:~foo/dosquad/go-giturl.git", "git", "", "foo", "github.com", "", "~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},
		{"bar@github.com:~foo/dosquad/go-giturl.git", "git", "bar", "foo", "github.com", "", "~foo/dosquad/go-giturl.git", "dosquad/go-giturl.git", "dosquad/go-giturl"},

		// /path/to/repo.git/
		{"/path/to/repo.git/", "file", "", "", "", "", "/path/to/repo.git/", "/path/to/repo.git/", "path/to/repo"},

		// file:///path/to/repo.git/
		{"file:///path/to/repo.git/", "file", "", "", "", "", "/path/to/repo.git/", "/path/to/repo.git/", "path/to/repo"},
	}

	for _, tt := range tests {
		t.Run(tt.rawURL, func(t *testing.T) {
			t.Parallel()

			u, err := giturl.Parse(tt.rawURL)
			if err != nil {
				t.Errorf("giturl.Parse: error, got '%s'", err)
				return
			}

			if v := u.User.Username(); v != tt.user {
				t.Errorf("url.User.Username(): got '%s', want '%s'", v, tt.user)
			}

			if v := u.PathUser.Username(); v != tt.pathuser {
				t.Errorf("url.PathUser.Username(): got '%s', want '%s'", v, tt.pathuser)
			}

			expectHost := tt.host
			if tt.port != "" {
				expectHost += ":" + tt.port
			}

			if v := u.Host; v != expectHost {
				t.Errorf("u.Host: got '%s', want '%s'", v, expectHost)
			}

			if v := u.Hostname(); v != tt.host {
				t.Errorf("url.Hostname(): got '%s', want '%s'", v, tt.host)
			}

			if v := u.Port(); v != tt.port {
				t.Errorf("url.Port(): got '%s', want '%s'", v, tt.port)
			}

			if v := u.Path; v != tt.path {
				t.Errorf("url.Path: got '%s', want '%s'", v, tt.path)
			}

			if v := u.PathRelative; v != tt.pathRelative {
				t.Errorf("url.PathRelative: got '%s', want '%s'", v, tt.pathRelative)
			}

			if v := u.Slug(); v != tt.slug {
				t.Errorf("url.Slug(): got '%s', want '%s'", v, tt.slug)
			}
		})
	}
}

// func TestGitURLString(t *testing.T) {
// 	tests := []struct {
// 		rawURL   string
// 		toString string
// 	}{
// 		// ssh://[<user>@]<host>[:<port>]/<path-to-git-repo>
// 		{
// 			"ssh://foo@githost.example:1234/path/to/git/repo",
// 			"ssh://foo@githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"ssh://githost.example:1234/path/to/git/repo",
// 			"ssh://githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"ssh://foo@githost.example/path/to/git/repo",
// 			"ssh://foo@githost.example/path/to/git/repo",
// 		},
// 		{
// 			"ssh://githost.example/path/to/git/repo",
// 			"ssh://githost.example/path/to/git/repo",
// 		},

// 		// git://<host>[:<port>]/<path-to-git-repo>
// 		{
// 			"git://githost.example:1234/path/to/git/repo",
// 			"git://githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"git://githost.example/path/to/git/repo",
// 			"git://githost.example/path/to/git/repo",
// 		},

// 		// http[s]://<host>[:<port>]/<path-to-git-repo>
// 		{
// 			"http://githost.example:1234/path/to/git/repo",
// 			"http://githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"https://githost.example:1234/path/to/git/repo",
// 			"https://githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"http://githost.example/path/to/git/repo",
// 			"http://githost.example/path/to/git/repo",
// 		},
// 		{
// 			"https://githost.example/path/to/git/repo",
// 			"https://githost.example/path/to/git/repo",
// 		},

// 		// ftp[s]://<host>[:<port>]/<path-to-git-repo>
// 		{
// 			"ftp://githost.example:1234/path/to/git/repo",
// 			"ftp://githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"ftps://githost.example:1234/path/to/git/repo",
// 			"ftps://githost.example:1234/path/to/git/repo",
// 		},
// 		{
// 			"ftp://githost.example/path/to/git/repo",
// 			"ftp://githost.example/path/to/git/repo",
// 		},
// 		{
// 			"ftps://githost.example/path/to/git/repo",
// 			"ftps://githost.example/path/to/git/repo",
// 		},

// 		// [<user>@]<host>:/<path-to-git-repo>
// 		{
// 			"git@github.com:dosquad/go-giturl.git",
// 			"git+ssh://git@github.com/dosquad/go-giturl.git",
// 		},

// 		// ssh://[<user>@]<host>[:<port>]/~<user>/<path-to-git-repo>
// 		{
// 			"ssh://github.com:1111/~foo/dosquad/go-giturl.git",
// 			"ssh://github.com:1111/~foo/dosquad/go-giturl.git",
// 		},
// 		{
// 			"ssh://github.com/~foo/dosquad/go-giturl.git",
// 			"ssh://github.com/~foo/dosquad/go-giturl.git",
// 		},
// 		{
// 			"ssh://bar@github.com/~foo/dosquad/go-giturl.git",
// 			"ssh://bar@github.com/~foo/dosquad/go-giturl.git",
// 		},

// 		// git://<host>[:<port>]/~<user>/<path-to-git-repo>
// 		{
// 			"git://github.com:1111/~foo/dosquad/go-giturl.git",
// 			"git://github.com:1111/~foo/dosquad/go-giturl.git",
// 		},
// 		{
// 			"git://github.com/~foo/dosquad/go-giturl.git",
// 			"git://github.com/~foo/dosquad/go-giturl.git",
// 		},

// 		// [<user>@]<host>:~<user>/<path-to-git-repo>
// 		{
// 			"github.com:~foo/dosquad/go-giturl.git",
// 			"github.com:~foo/dosquad/go-giturl.git",
// 		},
// 		{
// 			"bar@github.com:~foo/dosquad/go-giturl.git",
// 			"bar@github.com:~foo/dosquad/go-giturl.git",
// 		},

// 		// /path/to/repo.git/
// 		{
// 			"/path/to/repo.git/",
// 			"/path/to/repo.git/",
// 		},

// 		// file:///path/to/repo.git/
// 		{
// 			"file:///path/to/repo.git/",
// 			"file:///path/to/repo.git/",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.rawURL, func(t *testing.T) {
// 			u, err := giturl.Parse(tt.rawURL)
// 			if err != nil {
// 				t.Errorf("giturl.Parse: error, got '%s'", err)
// 				return
// 			}

// 			if v := u.String(); v != tt.toString {
// 				t.Errorf("u.String(): got '%s', want '%s'", v, tt.toString)
// 			}
// 		})
// 	}
// }
