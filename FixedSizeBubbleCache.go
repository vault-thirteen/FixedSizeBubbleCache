// Fixed Size Bubble Cache.

package fsbcache

import (
	"errors"
	"fmt"
	"sync"
)

// A fixed-Size Bubble Cache.
//
// It is called 'Bubble' while the old Records are lifted upwards when they are
// requested. This Process reminds the Air Bubbles travelling vertically in the
// Water. The most actively used Records are stored at the Top of the Cache.
// The least actively used Records are stored at the Bottom of the Cache.
// The maximum Size of the Cache (the Cache's Records Count) is fixed.
// This Structure reminds a classic Stack which receives new Items at the Top.
type FixedSizeBubbleCache struct {

	// The Top Record is the most fresh Record of the Cache.
	//
	// It might be either a newly added Record or an existing Record which has
	// just been requested and thus became the most fresh One.
	top *FixedSizeBubbleCacheRecord

	// The Bottom Record is a Record with the oldest Access Time.
	//
	// When a new Record is added to the Cache and the Cache is at its maximum
	// Size, the Bottom Record is removed from the Cache to keep its Size
	// constant.
	bottom *FixedSizeBubbleCacheRecord

	// The current Size of the Cache, the Count of the Cache's Records.
	size uint

	// The Capacity is the maximum Size of the Cache.
	//
	// A new Record added when the Cache is at its maximum Size, will remove the
	// Bottom Record from the Cache.
	capacity uint

	// An internal List of Records that may be fast requested by their unique
	// Identifier, a UID.
	recordsByUID sync.Map

	// Record's Time-To-Live (TTL) is the Period of Time, after which the
	// Record is considered outdated. If a Record requested from the Cache is
	// outdated, it is removed from the Cache. Record's TTL is measured in
	// Seconds.
	recordTTL uint
}

// Creates a new fixed-Size Bubble Cache.
func NewFixedSizeBubbleCache(
	capacity uint,
	recordTTL uint,
) (cache *FixedSizeBubbleCache) {
	if capacity == 0 {
		capacity++
	}
	cache = new(FixedSizeBubbleCache)
	cache.initialize(capacity, recordTTL)
	return
}

// Initializes the Cache.
func (c *FixedSizeBubbleCache) initialize(
	capacity uint,
	recordTTL uint,
) {
	c.top = nil
	c.bottom = nil
	c.size = 0
	c.capacity = capacity
	c.recordTTL = recordTTL
}

// Checks the Record's Parameters and adds it to the fixed-Size Bubble Cache.
func (c *FixedSizeBubbleCache) AddRecord(
	record *FixedSizeBubbleCacheRecord,
) (err error) {

	// Checks.
	if record == nil {
		return errors.New(ErrRecordIsNotSet)
	}
	err = record.Check()
	if err != nil {
		return
	}

	// Addition.
	c.addRecord(record)
	return
}

// Adds a Record to the Cache.
//
//	If the Record with a specified UID already exists in the Cache,
// 		then that existing Record is moved to the Top Position and its
//		Contents are updated.
//	If the Record with a specified UID does not exist in the Cache,
//		then a new Top Record is inserted into the Cache.
//	If the Cache is at its maximum Size (Size is equal to Capacity) and a new
//	Record must be added,
//		then the Bottom Record is removed from the Cache.
func (c *FixedSizeBubbleCache) addRecord(
	addedRecord *FixedSizeBubbleCacheRecord,
) {
	var existingRecordIfc interface{}
	var uidExists bool
	existingRecordIfc, uidExists = c.recordsByUID.Load(addedRecord.UID)
	if uidExists {
		var existingRecord *FixedSizeBubbleCacheRecord
		var ok bool
		existingRecord, ok = existingRecordIfc.(*FixedSizeBubbleCacheRecord)
		if !ok {
			panic(ErrTypeCast)
		}
		if existingRecord != c.top {
			c.moveExistingRecordToTop(existingRecord)
		}
		c.top.UpdateDataAndLAT(addedRecord.Data)
		return
	}
	if c.size == c.capacity {
		var removedRecord *FixedSizeBubbleCacheRecord = c.unlinkBottomRecord()
		c.recordsByUID.Delete(removedRecord.UID)
	}
	c.linkTopRecord(addedRecord)

	c.recordsByUID.Store(addedRecord.UID, addedRecord)
	c.top.UpdateDataAndLAT(addedRecord.Data)
	if c.size != c.capacity {
		c.size++ // We can not increase the Size prior to Linking.
	}
}

// Moves the existing Record to the Top.
// The Record must be a non-Top Record.
// This Method only manipulates the Links, so it does not check the Integrity
// (you must do the Checks beforehand), it does not touch the fast-Access Map,
// it does not touch the Size Counter.
func (c *FixedSizeBubbleCache) moveExistingRecordToTop(
	existingRecord *FixedSizeBubbleCacheRecord,
) {
	// When the Record is not the Top, then the Cache Size is 2 or more.
	if c.size == 2 {
		c.unlinkBottomRecord()
	} else { // size > 2.
		if existingRecord == c.bottom {
			c.unlinkBottomRecord()
		} else {
			c.unlinkMiddleRecord(existingRecord)
		}
	}
	c.linkTopRecord(existingRecord)
}

// Unlinks the Top Record and returns it.
// This Method only manipulates the Links, so it does not check the Integrity
// (you must do the Checks beforehand), it does not touch the fast-Access Map,
// it does not touch the Size Counter.
func (c *FixedSizeBubbleCache) unlinkTopRecord() (oldTop *FixedSizeBubbleCacheRecord) {
	oldTop = c.top
	c.top = oldTop.lowerRecord
	c.top.upperRecord = nil
	oldTop.lowerRecord = nil
	oldTop.upperRecord = nil
	return
}

// Unlinks the Middle Record and returns it.
// This Method only manipulates the Links, so it does not check the Integrity
// (you must do the Checks beforehand), it does not touch the fast-Access Map,
// it does not touch the Size Counter.
func (c *FixedSizeBubbleCache) unlinkMiddleRecord(
	record *FixedSizeBubbleCacheRecord,
) *FixedSizeBubbleCacheRecord {
	record.upperRecord.lowerRecord = record.lowerRecord
	record.lowerRecord.upperRecord = record.upperRecord
	record.upperRecord = nil
	record.lowerRecord = nil
	return record
}

// Unlinks the Bottom Record and returns it.
// This Method only manipulates the Links, so it does not check the Integrity
// (you must do the Checks beforehand), it does not touch the fast-Access Map,
// it does not touch the Size Counter.
func (c *FixedSizeBubbleCache) unlinkBottomRecord() (oldBottom *FixedSizeBubbleCacheRecord) {
	oldBottom = c.bottom
	c.bottom = oldBottom.upperRecord
	c.bottom.lowerRecord = nil
	oldBottom.upperRecord = nil
	oldBottom.lowerRecord = nil
	return
}

// Inserts (connects) the Top Record and returns it.
// This Method only manipulates the Links, so it does not check the Integrity
// (you must do the Checks beforehand), it does not touch the fast-Access Map,
// it does not touch the Size Counter.
func (c *FixedSizeBubbleCache) linkTopRecord(
	newTop *FixedSizeBubbleCacheRecord,
) *FixedSizeBubbleCacheRecord {
	if c.size == 0 {
		c.top = newTop
		c.bottom = newTop
		newTop.upperRecord = nil
		newTop.lowerRecord = nil
	} else {
		c.top.upperRecord = newTop
		newTop.upperRecord = nil
		newTop.lowerRecord = c.top
		c.top = newTop
	}
	return newTop
}

// Inserts (connects) the Bottom Record and returns it.
// This Method only manipulates the Links, so it does not check the Integrity
// (you must do the Checks beforehand), it does not touch the fast-Access Map,
// it does not touch the Size Counter.
func (c *FixedSizeBubbleCache) linkBottomRecord(
	newBottom *FixedSizeBubbleCacheRecord,
) *FixedSizeBubbleCacheRecord {
	if c.size == 0 {
		c.top = newBottom
		c.bottom = newBottom
		newBottom.upperRecord = nil
		newBottom.lowerRecord = nil
	} else {
		c.bottom.lowerRecord = newBottom
		newBottom.upperRecord = c.bottom
		newBottom.lowerRecord = nil
		c.bottom = newBottom
	}
	return newBottom
}

// Deletes all Records from the Cache.
// As opposed to other Deletion Methods, this Method uses the Integrity Check.
func (c *FixedSizeBubbleCache) Clear() (err error) {

	// Before deleting the Records, we must ensure that Cache is not broken.
	// Broken Cache Deletion would cost us a lot of Memory Leaks!
	if !c.isIntegral() {
		err = errors.New(ErrIntegrityCheckFailure)
		return
	}

	// Delete all Items from Top to Bottom.
	for c.size > 0 {
		err = c.deleteRecord(c.bottom, false)
		if err != nil {
			return
		}
	}
	return
}

// Checks the Integrity of the Cache.
// This is a Self-Check Function intended to find Anomalies.
// This Function is not intended to be used in an ordinary Case.
// Returns 'true' if the Cache is in a good Shape.
func (c *FixedSizeBubbleCache) isIntegral() bool {

	// Check Fast Access Register.
	var ok bool = true
	var nilSearcher = func(key, value interface{}) bool {
		if value == nil {
			ok = false
			return false
		}
		return true
	}
	c.recordsByUID.Range(nilSearcher)
	if !ok {
		return false
	}

	// Prepare Data.
	var top *FixedSizeBubbleCacheRecord = c.top
	var bottom *FixedSizeBubbleCacheRecord = c.bottom
	var size uint = c.size

	// Capacity Check.
	if size > c.capacity {
		return false
	}

	// Empty List?
	if size == 0 {
		if top != nil {
			return false
		}
		if bottom != nil {
			return false
		}
		return true
	}

	// Single-Item List?
	if size == 1 {
		if top == nil {
			return false
		}
		if bottom == nil {
			return false
		}
		if top != bottom {
			return false
		}
		if top.upperRecord != nil {
			return false
		}
		if bottom.lowerRecord != nil {
			return false
		}
		return true
	}

	// List has two or more Items.

	// Check the Top.
	if top.upperRecord != nil {
		return false
	}
	// Check the Bottom.
	if bottom.lowerRecord != nil {
		return false
	}

	// Try to inspect all the Items from Top to Bottom.
	// This checks Connectivity by the 'lower' Pointer.
	var cursor *FixedSizeBubbleCacheRecord
	var cursorLowerRecord *FixedSizeBubbleCacheRecord
	var cursorUpperRecord *FixedSizeBubbleCacheRecord
	cursor = top
	cursorLowerRecord = cursor.lowerRecord
	var i uint = 1
	var sizeAnomaly bool
	for cursorLowerRecord != nil {
		cursor = cursorLowerRecord
		cursorLowerRecord = cursor.lowerRecord
		i++
		// Defence against Self-Loop Anomaly.
		if i > size {
			sizeAnomaly = true
			break
		}
	}
	if sizeAnomaly {
		// Size Anomaly can happen if we either have a Self-Loop in the Chain
		// or the Corner Item for some Reason is not the End of the Chain.
		return false
	}
	if i != size {
		return false
	}
	// We have stopped the Search at the first Break in the Chain.
	// Are we really there where we should be?
	if cursor != bottom {
		// We have found a broken Connection.
		return false
	}

	// Now, try to inspect all Items in a reversed Order.
	// This checks Connectivity by the 'upper' Pointer.
	cursor = bottom
	cursorUpperRecord = cursor.upperRecord
	i = 1
	for cursorUpperRecord != nil {
		cursor = cursorUpperRecord
		cursorUpperRecord = cursor.upperRecord
		i++
		// Defence against Self-Loop Anomaly.
		if i > size {
			sizeAnomaly = true
			break
		}
	}
	if sizeAnomaly {
		// Size Anomaly can happen if we either have a Self-Loop in the Chain
		// or the Corner Item for some Reason is not the End of the Chain.
		return false
	}
	if i != size {
		return false
	}
	// We have stopped the Search at the first Break in the Chain.
	// Are we really there where we should be?
	if cursor != top {
		// We have found a broken Connection.
		return false
	}

	// All Clear.
	return true
}

// Deletes a Record from the Cache.
func (c *FixedSizeBubbleCache) deleteRecord(
	record *FixedSizeBubbleCacheRecord,
	recordIsKnownToExist bool, // A Flag to avoid the Existence Check.
) (err error) {

	// Fool Check.
	if record == nil {
		err = errors.New(ErrRecordIsNotSet)
		return
	}
	if c.size == 0 {
		err = errors.New(ErrCacheZeroSize)
		return
	}
	if !recordIsKnownToExist {
		var recordExists bool = c.recordUIDExists(record.UID)
		if !recordExists {
			err = fmt.Errorf(ErrfRecordWithUidIsNotFound, record.UID)
			return
		}
	}

	// A single-Record Cache.
	if c.size == 1 {
		c.top = nil
		c.bottom = nil
		c.size--
		c.recordsByUID.Delete(record.UID)
		return
	}

	// A normal Deletion, Size is >= 2.
	if record == c.top { // Removed Record is the Top Record.
		c.unlinkTopRecord()
	} else if record == c.bottom { // Removed Record is the Bottom Record.
		c.unlinkBottomRecord()
	} else { // Removed Record is somewhere in the Middle.
		c.unlinkMiddleRecord(record)
	}
	c.size--
	c.recordsByUID.Delete(record.UID)
	return
}

// Checks whether the specified Record's UID exists in the Cache's List.
func (c *FixedSizeBubbleCache) RecordUIDExists(
	uid FixedSizeBubbleCacheRecordUID,
) (uidExists bool) {
	return c.recordUIDExists(uid)
}

// Checks whether the specified Record's UID exists in the Cache's List.
func (c *FixedSizeBubbleCache) recordUIDExists(
	uid FixedSizeBubbleCacheRecordUID,
) (uidExists bool) {
	_, uidExists = c.recordsByUID.Load(uid)
	return
}

// Gets the Record from the Cache's internal List.
func (c *FixedSizeBubbleCache) getRecordByUID(
	uid FixedSizeBubbleCacheRecordUID,
) (record *FixedSizeBubbleCacheRecord, err error) {
	var recordIsFound bool
	var recordIfc interface{}
	recordIfc, recordIsFound = c.recordsByUID.Load(uid)
	if !recordIsFound {
		err = fmt.Errorf(ErrfRecordWithUidIsNotFound, uid)
		return
	}
	var ok bool
	record, ok = recordIfc.(*FixedSizeBubbleCacheRecord)
	if !ok {
		panic(ErrTypeCast)
	}
	return
}

// Deletes a Record specified by its UID from the Cache.
func (c *FixedSizeBubbleCache) DeleteRecordByUID(
	uid FixedSizeBubbleCacheRecordUID,
) (err error) {
	var record *FixedSizeBubbleCacheRecord
	record, err = c.getRecordByUID(uid)
	if err != nil {
		return
	}
	err = c.deleteRecord(record, true)
	if err != nil {
		return
	}
	return
}

// Lists the Values of all Records of the Cache.
func (c *FixedSizeBubbleCache) ListAllRecordValues() (values []interface{}) {
	values = make([]interface{}, c.size)
	if c.size == 0 {
		return
	}

	// Get the first Item.
	var record *FixedSizeBubbleCacheRecord = c.top
	values[0] = record.Data

	// Get all the rest Items.
	var i uint
	for i = 1; i < c.size; i++ {
		record = record.lowerRecord
		values[i] = record.Data
	}
	return
}

// Lists all the Records of the Cache.
func (c *FixedSizeBubbleCache) ListAllRecords() (records []*FixedSizeBubbleCacheRecord) {
	records = make([]*FixedSizeBubbleCacheRecord, c.size)
	if c.size == 0 {
		return records
	}

	// Get the first Item.
	var record *FixedSizeBubbleCacheRecord = c.top
	records[0] = record

	// Get all the rest Items.
	var i uint
	for i = 1; i < c.size; i++ {
		record = record.lowerRecord
		records[i] = record
	}
	return
}

// Gets the Record's Data by its UID. Moves the Record to the Top of the List
// and refreshes its LAT. If the Record is outdated, deletes it and returns an
// Error.
func (c *FixedSizeBubbleCache) GetActualRecordDataByUID(
	uid FixedSizeBubbleCacheRecordUID,
) (data interface{}, err error) {

	// Get the Record.
	var record *FixedSizeBubbleCacheRecord
	record, err = c.getRecordByUID(uid)
	if err != nil {
		return
	}

	// Check the TTL. Is the Record Outdated ?
	if !record.isActual(c.recordTTL) {
		err = c.deleteRecord(record, true)
		if err != nil {
			return
		}
		err = fmt.Errorf(ErrfRecordWithUidIsOutdated, uid)
		return
	}

	// Move the requested Record to the Top.
	if record != c.top {
		c.moveExistingRecordToTop(record)
	}
	c.top.UpdateLAT()

	data = record.Data
	return
}

// Returns the 'RecordTTL' Parameter of the Cache.
func (c *FixedSizeBubbleCache) GetRecordTTL() uint {
	return c.recordTTL
}

// Checks whether the specified Record's UID exists in the Cache and
// the Record with such UID is still active (not outdated).
func (c *FixedSizeBubbleCache) IsRecordUIDActive(
	uid FixedSizeBubbleCacheRecordUID,
) (recordIsActive bool, err error) {

	// Get the Record.
	var record *FixedSizeBubbleCacheRecord
	record, err = c.getRecordByUID(uid)
	if err != nil {
		return
	}

	recordIsActive = record.isActual(c.recordTTL)
	return
}
