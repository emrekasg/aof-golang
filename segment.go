package aof

import (
	"os"
)

type segment struct {
	path  string // path to the segment file
	index uint64 // index of the segment file

	// todo:
	// cbuf  []byte   // cached entries of buffer
	// cpos  []bpos   // cached position of the buffer

	// currentBlockNumber uint64     // current block number
	// currentBlockSize   uint64     // current block size
	// blockPool          *sync.Pool // block pool
}

// type bpos struct {
// 	pos int // position of the entry in the segment file
// 	end int // length of the entry
// }

func (aof *AOF) createSegment(path string, index uint64) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	aof.segment = file
	aof.segments = append(aof.segments, &segment{
		path:  path,
		index: index,
	})

	return nil
}
