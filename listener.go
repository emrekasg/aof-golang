package aof

// Listener  uns after aofChan receive a command
type Listener interface {
	Callback([][][]byte)
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
