package views

import "embed"

//go:embed *html
//go:embed javascript/*js
//go:embed components/*html
//go:embed css/*css
var Static embed.FS
