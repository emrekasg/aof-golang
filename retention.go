package aof

// we need retention to be able to check if size of the aof file is greater than 20MB
// if it is, we need to create a new segment and write to it
func (aof *AOF) Retention() error {
	stat, err := aof.segment.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > aofMaxSize {
		if err := aof.cycle(); err != nil {
			return err
		}
	}

	return nil
}
