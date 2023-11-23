package aof

// Listener going to run after aofChan receive a command
// and before the command is written to the aof file
type Listener interface {
	Callback([]CmdLine)
}

func (aof *AOF) AddListener(listener Listener) {
	aof.sync.Lock()
	defer aof.sync.Unlock()
	aof.listeners[listener] = struct{}{}
}

func (aof *AOF) RemoveListener(listener Listener) {
	aof.sync.Lock()
	defer aof.sync.Unlock()
	delete(aof.listeners, listener)
}
