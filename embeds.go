package main

import "embed"

//go:embed templates/*
var templates embed.FS

//go:embed static/*
var static embed.FS
