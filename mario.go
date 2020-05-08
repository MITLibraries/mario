package mario

import "github.com/markbates/pkger"

// This file seems to be needed to make pkger work with the new
// project layout. I think it may be related to
// https://github.com/markbates/pkger/issues/86.
func init() {
	pkger.Include("/config")
}
