package aof

import "sync"

type CmdLine = [][]byte

type AofCmd struct {
	CmdLine CmdLine
	DbIndex int
	Wg      *sync.WaitGroup
}
