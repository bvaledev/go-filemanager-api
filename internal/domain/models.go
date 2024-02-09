package domain

import "io"

type FileInfo struct {
	Name string
	File io.Reader
}
