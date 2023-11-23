package aof

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

const (
	readerBufferSize = 32 << 10 // 32 KB
)

type Reader struct {
	sync.RWMutex
	segments     []*segment
	currentIndex int
	file         *os.File // Current open file
	bufferSize   int      // Size of each read chunk
}

func NewReader(folderPath string) (*Reader, error) {
	if folderPath == "" {
		return nil, fmt.Errorf("folder path cannot be empty")
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	var segments []*segment
	for _, file := range files {
		if !file.IsDir() {
			segments = append(segments, &segment{path: filepath.Join(folderPath, file.Name())})
		}
	}

	return &Reader{
		segments:   segments,
		bufferSize: readerBufferSize,
	}, nil
}

// io.Reader Interface Implementation
func (r *Reader) Read(p []byte) (int, error) {
	r.Lock()
	defer r.Unlock()

	for {
		// If no file is open or the end of the current file is reached, open the next file
		if r.file == nil {
			if r.currentIndex >= len(r.segments) {
				return 0, io.EOF // No more segments to read
			}
			var err error
			r.file, err = os.Open(r.segments[r.currentIndex].path)
			if err != nil {
				return 0, fmt.Errorf("error opening segment file: %w", err)
			}
		}

		n, err := r.file.Read(p)
		if err != nil {
			// Close the current file on error or EOF and move to the next segment
			r.file.Close()
			r.file = nil
			r.currentIndex++

			if err == io.EOF && r.currentIndex < len(r.segments) {
				continue // EOF reached, but more files are available
			}
			return n, err
		}

		return n, nil // Data read into p, return number of bytes read
	}
}
