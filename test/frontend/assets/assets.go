package frontend

import "embed"

//go:embed index.html
//go:embed script.js
//go:embed style.css
//go:embed dist/*
//go:embed wailsjs/*
var Assets embed.FS
