package msys_cmd_vars

import (
	"bytes"
	"sync"
)

// supported linux distro/flavor values
const (
	UBUNTU1604 = 0
	UBUNTU1804 = 1
	CENTOS7 = 2
	UNKNOWN = 255
)

type SafeWriter struct {
	buffer    *bytes.Buffer
	lines     []string
	*sync.Mutex   //embed mutex into struct itself
}

func NewSafeWriter() *SafeWriter {
	return &SafeWriter{
		&bytes.Buffer{},
		[]string{},
		&sync.Mutex{},
	}
}

// io.Writer write
func (sw *SafeWriter) Write(p []byte) (int, error) {
	sw.Lock()
	defer sw.Unlock()
	return sw.buffer.Write(p)
}

func (sw *SafeWriter) GetBytes() []byte {
	sw.Lock()
	defer sw.Unlock()
	return sw.buffer.Bytes()
}

