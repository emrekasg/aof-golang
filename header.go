package aof

import (
	"encoding/binary"
	"time"
)

const (
	FileHeaderSize = SizeUint16 + // db version
		SizeUint64 // timestamp
)

type FileHeader struct {
	DbVersion uint16
	Timestamp uint64
}

func (fh *FileHeader) Encode() []byte {
	b := make([]byte, FileHeaderSize)
	binary.BigEndian.PutUint16(b[0:2], fh.DbVersion)
	binary.BigEndian.PutUint64(b[2:10], fh.Timestamp)

	return b
}

func DecodeFileHeader(b []byte) *FileHeader {
	return &FileHeader{
		DbVersion: binary.BigEndian.Uint16(b[0:2]),
		Timestamp: binary.BigEndian.Uint64(b[2:10]),
	}
}

func (fh *FileHeader) SizeInBytes() int {
	return FileHeaderSize
}

func (fh *FileHeader) Create(b []byte) {
	fh.DbVersion = uint16(1)
	fh.Timestamp = uint64(time.Now().UnixNano())
	// todo: i'll consider adding different header types (maybe size of data or sth?)
}
