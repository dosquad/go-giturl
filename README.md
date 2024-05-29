# go-giturl

[![CI](https://github.com/dosquad/go-giturl/actions/workflows/ci.yml/badge.svg)](https://github.com/dosquad/go-giturl/actions/workflows/ci.yml)
[![GoDoc](https://godoc.org/github.com/dosquad/go-giturl/?status.svg)](https://godoc.org/github.com/dosquad/go-giturl)
[![GitHub issues](https://img.shields.io/github/issues/dosquad/go-giturl)](https://github.com/dosquad/go-giturl/issues)
[![GitHub forks](https://img.shields.io/github/forks/dosquad/go-giturl)](https://github.com/dosquad/go-giturl/network)
[![GitHub stars](https://img.shields.io/github/stars/dosquad/go-giturl)](https://github.com/dosquad/go-giturl/stargazers)
[![GitHub license](https://img.shields.io/github/license/dosquad/go-giturl)](https://github.com/dosquad/go-giturl/blob/main/LICENSE)

Parse git remote URLs.

## Example

```golang
package main

import (
    "fmt"

    "github.com/dosquad/go-giturl"
)

// Output: slug: dosquad/go-giturl
func main() {
    u, err := giturl.Parse("git@github.com:dosquad/go-giturl.git")
    if err != nil {
        panic(err)
    }

    fmt.Printf("slug: %s", u.Slug())
}
```
