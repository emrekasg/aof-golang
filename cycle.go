package aof

import (
	"fmt"
	"os"
	"path/filepath"
)

// cycle creates a new segment and updates aof with the new segment.
func (aof *AOF) cycle() error {
	aof.sync.Lock()
	defer aof.sync.Unlock()

	// Ensure the current segment is not nil
	if aof.segment == nil {
		return fmt.Errorf("no current segment to cycle")
	}

	// Sync and close the current segment
	if err := aof.segment.Sync(); err != nil {
		return fmt.Errorf("error syncing current segment: %w", err)
	}
	if err := aof.segment.Close(); err != nil {
		return fmt.Errorf("error closing current segment: %w", err)
	}

	// Calculate the index for the new segment
	newIndex := uint64(1)
	if len(aof.segments) > 0 {
		newIndex = aof.segments[len(aof.segments)-1].index + 1
	}

	newPath := filepath.Join(aof.Options.Path, fmt.Sprintf("%014d", newIndex))

	// Create and open a new segment
	newSegment, err := os.OpenFile(newPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error creating new segment: %w", err)
	}

	// Update aof with the new segment
	aof.segment = newSegment
	aof.segments = append(aof.segments, &segment{
		path:  newPath,
		index: newIndex,
	})

	return nil
}
