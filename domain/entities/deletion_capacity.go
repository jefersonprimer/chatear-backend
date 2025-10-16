package entities

import (
	"time"
)

// DeletionCapacity represents the daily deletion capacity limits
type DeletionCapacity struct {
	Day      time.Time `json:"day"`
	Count    int       `json:"count"`
	MaxLimit int       `json:"max_limit"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewDeletionCapacity creates a new deletion capacity entry
func NewDeletionCapacity(day time.Time, maxLimit int) *DeletionCapacity {
	return &DeletionCapacity{
		Day:      day,
		Count:    0,
		MaxLimit: maxLimit,
		UpdatedAt: time.Now(),
	}
}

// CanDelete checks if a deletion can be performed
func (dc *DeletionCapacity) CanDelete() bool {
	return dc.Count < dc.MaxLimit
}

// IncrementCount increments the deletion count
func (dc *DeletionCapacity) IncrementCount() {
	dc.Count++
	dc.UpdatedAt = time.Now()
}
