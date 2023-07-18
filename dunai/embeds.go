package dunai

import "embed"

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS

//go:embed data/cv.yaml
var cvYaml []byte
