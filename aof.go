package aof

import (
	"bytes"
	"context"
	"os"
	"sync"
	"time"
)

type AOFOptions struct {
	Path  string // absolute path to the aof log directory
	Fsync string // file synchronization options
	// todo: blockCache  uint32
	// todo: bytesPerSync uint64
}

const (
	aofQueueSize = 1 << 20  // 1 MB size of the channel used to communicate
	aofMaxSize   = 20 << 20 // 20 MB, max size of an aof file

	FsyncAlways      = "always"   // if you want to be sure that every single operation is written to disk (it's really so slow but guarantees data integrity)
	FsyncNo          = "no"       // if you don't care about data loss in case of a crash
	FsyncEverySecond = "everysec" // (default)
)

var (
	ErrEOF = os.ErrClosed

	DefaultOptions = &AOFOptions{
		Path:  "/tmp/letly_aof",
		Fsync: FsyncAlways,
	}
)

type AOF struct {
	sync sync.RWMutex

	ctx       context.Context
	cancel    context.CancelFunc
	aofChan   chan *AofCmd
	Options   *AOFOptions
	listeners map[Listener]struct{}
	buffer    [][][]byte
	segment   *os.File
	segments  []*segment

	reader *Reader
}

type AofCmd struct {
	CmdLine [][]byte
	DbIndex int
	Wg      *sync.WaitGroup
}

func NewAOF(aofOptions *AOFOptions) (*AOF, error) {
	var err error

	if aofOptions == nil {
		aofOptions = DefaultOptions
	}

	aof := &AOF{
		Options:   aofOptions,
		aofChan:   make(chan *AofCmd, aofQueueSize),
		listeners: make(map[Listener]struct{}),
	}

	aof.aofChan = make(chan *AofCmd, aofQueueSize)
	aof.listeners = make(map[Listener]struct{})

	aof.ctx, aof.cancel = context.WithCancel(context.Background())

	if err = os.MkdirAll(aof.Options.Path, 0755); err != nil {
		return nil, err
	}

	if err = aof.loadFiles(); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case <-aof.ctx.Done():
				return
			default:
				if FsyncEverySecond == aof.Options.Fsync {
					aof.fsync() // fsync every second
				}
				aof.Retention()
				time.Sleep(time.Second * 1)
			}
		}
	}()

	return aof, nil

}

// WriteAof writes a new entry to the Append-Only Log.
func (aof *AOF) WriteAof(p *AofCmd) error {
	aof.sync.Lock()
	defer aof.sync.Unlock()
	aof.buffer = aof.buffer[:0] // clear the buffer

	aof.buffer = append(aof.buffer, p.CmdLine)

	// convert buffer to bytes
	bufferInBytes := aof.bufferToBytes()

	// write to aof file in byte format
	if _, err := aof.segment.Write(bufferInBytes); err != nil {
		return err
	}

	if FsyncAlways == aof.Options.Fsync {
		aof.segment.Sync()
	}

	// notify listeners
	for listener := range aof.listeners {
		listener.Callback(aof.buffer)
	}

	return nil
}

// bufferToBytes converts the buffer to bytes and adds a file header.
func (aof *AOF) bufferToBytes() []byte {
	var buffer []byte
	for _, cmdLine := range aof.buffer {
		buffer = append(buffer, bytes.Join(cmdLine, []byte(" "))...)
	}

	var fileHeader FileHeader
	fileHeader.Create(buffer)

	fileHeaderInBytes := fileHeader.Encode()[:] // Encode to bytes
	fileHeaderInBytes = append(fileHeaderInBytes, ' ')

	buffer = append(fileHeaderInBytes, buffer...)
	buffer = append(buffer, '\n')

	return buffer
}

func (aof *AOF) GetReader(buffer []byte) (*Reader, error) {
	aof.sync.RLock()
	defer aof.sync.RUnlock()

	var err error

	aof.reader, err = NewReader(aof.Options.Path)
	if err != nil {
		return nil, err
	}

	return aof.reader, nil
}
