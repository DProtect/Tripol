// Package stdlib embeds the official TRIPOL standard-library source.
package stdlib

import "embed"

// Sources contains the .tri module files at their original location.
//go:embed *.tri
var Sources embed.FS
