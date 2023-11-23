package aof

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

var (
	ErrCorrupt = errors.New("corrupt aof file")
	ErrClosed  = errors.New("aof closed")
	ErrAof     = errors.New("aof error")
)

func (aof *AOF) fsync() {
	aof.sync.Lock()
	defer aof.sync.Unlock()

	aof.segment.Sync()
}

func (aof *AOF) loadFiles() error {
	files, err := os.ReadDir(aof.Options.Path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || len(file.Name()) < 14 {
			continue
		}

		index, err := strconv.ParseUint(file.Name()[:14], 10, 64)
		if err != nil || index == 0 {
			continue
		}

		aof.segments = append(aof.segments, &segment{
			path:  filepath.Join(aof.Options.Path, file.Name()),
			index: index,
		})
	}

	if len(aof.segments) == 0 {
		segmentPath := filepath.Join(aof.Options.Path, fmt.Sprintf("%014d", 1))
		return aof.createSegment(segmentPath, 1)
	}

	lastSegment := aof.segments[len(aof.segments)-1]
	return aof.createSegment(lastSegment.path, lastSegment.index)
}
