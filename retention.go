package aof

// Checks if the segment file is greater than the maximum size.
// If it is, then create a new segment file and close the current one.
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
