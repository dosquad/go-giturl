package giturl_test

import (
	"fmt"

	"github.com/dosquad/go-giturl"
)

func ExampleParse() {
	u, err := giturl.Parse("git@github.com:dosquad/go-giturl.git")
	if err != nil {
		panic(err)
	}

	fmt.Printf("slug: %s", u.Slug())
	// Output: slug: dosquad/go-giturl
}
