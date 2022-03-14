package main

import (
// "gorm.io/gorm"
)

type Project struct {
	ID          uint `gorm:"primarykey"`
	Order       int
	Title       string
	Link        string
	Description string
	Stars       int
}
