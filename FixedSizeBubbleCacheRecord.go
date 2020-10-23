package fsbcache

import (
	"errors"
	"time"
)

type FixedSizeBubbleCacheRecord struct {

	// A unique Identifier of the Record.
	UID FixedSizeBubbleCacheRecordUID

	// Some useful Data stored in the Record.
	Data interface{}

	// Time of the last Access to the Record.
	lastAccessTime uint

	// A Pointer to an upper Record.
	upperRecord *FixedSizeBubbleCacheRecord

	// A Pointer to a lower Record.
	lowerRecord *FixedSizeBubbleCacheRecord
}

// Checks the Record before insertion into the Cache.
func (r *FixedSizeBubbleCacheRecord) Check() (err error) {

	// Check the 'Data' Field.
	if r.Data == nil {
		return errors.New(ErrDataIsEmpty)
	}

	// Check the 'UID' Field.
	if len(r.UID) == 0 {
		return errors.New(ErrUIDIsEmpty)
	}
	return
}

// Updates the Record's Data with the Data provided and the Last Access Time
// with the current Time.
func (r *FixedSizeBubbleCacheRecord) UpdateDataAndLAT(
	data interface{},
) {
	r.Data = data
	r.updateLATWithCurrentTime()
}

// Updates the Record's Last Access Time with the current Time.
func (r *FixedSizeBubbleCacheRecord) UpdateLAT() {
	r.updateLATWithCurrentTime()
}

// Updates the Record's Last Access Time with the current Time.
func (r *FixedSizeBubbleCacheRecord) updateLATWithCurrentTime() {
	r.lastAccessTime = uint(time.Now().Unix())
}

// Checks whether the Record is outdated or not.
func (r *FixedSizeBubbleCacheRecord) isActual(
	ttl uint,
) bool {
	return uint(time.Now().Unix()) < r.lastAccessTime+ttl
}
