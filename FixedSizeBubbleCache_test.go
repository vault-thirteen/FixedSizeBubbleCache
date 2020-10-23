// Fixed Size Bubble Cache.

package fsbcache

import (
	"fmt"
	"testing"
	"time"

	"github.com/vault-thirteen/tester"
)

func Test_NewFixedSizeBubbleCache(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache *FixedSizeBubbleCache

	// Test #1. Non-Zero Capacity.
	cache = NewFixedSizeBubbleCache(2, 60)
	aTest.MustBeEqual(cache.capacity, uint(2))

	// Test #2. Zero Capacity.
	cache = NewFixedSizeBubbleCache(0, 60)
	aTest.MustBeEqual(cache.capacity, uint(1))
}

func Test_initialize(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache *FixedSizeBubbleCache = new(FixedSizeBubbleCache)
	cache.initialize(10, 60)

	// Test #1.
	aTest.MustBeEqual(cache.top, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.size, uint(0))
	aTest.MustBeEqual(cache.capacity, uint(10))
	aTest.MustBeDifferent((cache.recordsByUID), nil)
	aTest.MustBeEqual(len(cache.recordsByUID), int(0))
	aTest.MustBeEqual(cache.recordTTL, uint(60))
}

func Test_AddRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(10, 60)
	var err error

	// Test #1. Null Record.
	err = cache.AddRecord(nil)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrRecordIsNotSet)

	// Test #2. Bad Record.
	err = cache.AddRecord(
		&FixedSizeBubbleCacheRecord{
			Data: nil,
		},
	)
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)

	// Test #3. Normal Record.
	err = cache.AddRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 101,
			UID:  "101",
		},
	)
	aTest.MustBeNoError(err)
}

func Test_addRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(2, 60)

	// Test #1. A new Record.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 101,
			UID:  "101",
		},
	)
	aTest.MustBeEqual(cache.size, uint(1))
	aTest.MustBeEqual(cache.top.UID, "101")
	aTest.MustBeEqual(cache.top.Data, 101)

	// Test #2. A Record with an existing UID, a Top Record.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 102,
			UID:  "101",
		},
	)
	aTest.MustBeEqual(cache.size, uint(1))
	aTest.MustBeEqual(cache.top.UID, "101")
	aTest.MustBeEqual(cache.top.Data, 102)

	// Test #3. A Record with a new UID.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 200,
			UID:  "200",
		},
	)
	aTest.MustBeEqual(cache.size, uint(2))
	aTest.MustBeEqual(cache.top.UID, "200")
	aTest.MustBeEqual(cache.top.Data, 200)

	// Test #4. A Record with an existing UID, not a Top Record.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 103,
			UID:  "101",
		},
	)
	aTest.MustBeEqual(cache.size, uint(2))
	aTest.MustBeEqual(cache.top.UID, "101")
	aTest.MustBeEqual(cache.top.Data, 103)

	// Test #5. A Record is new and the Cache is full.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 300,
			UID:  "300",
		},
	)
	aTest.MustBeEqual(cache.size, uint(2))
	aTest.MustBeEqual(cache.top.UID, "300")
	aTest.MustBeEqual(cache.top.Data, 300)
}

func Test_moveExistingRecordToTop(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(3, 60)

	// Test #1. Size = 2.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 1,
			UID:  "First",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 2,
			UID:  "Second",
		},
	)
	cache.moveExistingRecordToTop(cache.bottom)
	aTest.MustBeEqual(cache.top.UID, "First")
	aTest.MustBeEqual(cache.top.Data, 1)
	aTest.MustBeEqual(cache.bottom.UID, "Second")
	aTest.MustBeEqual(cache.bottom.Data, 2)

	// Test #2. Size > 2, a Bottom.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 3,
			UID:  "Third",
		},
	)
	cache.moveExistingRecordToTop(cache.bottom)
	aTest.MustBeEqual(cache.top.UID, "Second")
	aTest.MustBeEqual(cache.top.Data, 2)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)

	// Test #3. Size > 2, Not a Bottom.
	cache.moveExistingRecordToTop(cache.bottom.upperRecord)
	aTest.MustBeEqual(cache.top.UID, "Third")
	aTest.MustBeEqual(cache.top.Data, 3)
	aTest.MustBeEqual(cache.top.lowerRecord.UID, "Second")
	aTest.MustBeEqual(cache.top.lowerRecord.Data, 2)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)
}

func Test_unlinkTopRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(3, 60)

	// Test #1.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 1,
			UID:  "First",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 2,
			UID:  "Second",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 3,
			UID:  "Third",
		},
	)
	var oldTop *FixedSizeBubbleCacheRecord = cache.top
	cache.unlinkTopRecord()
	//
	aTest.MustBeEqual(oldTop.UID, "Third")
	aTest.MustBeEqual(oldTop.Data, 3)
	aTest.MustBeEqual(oldTop.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(oldTop.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	//
	aTest.MustBeEqual(cache.top.UID, "Second")
	aTest.MustBeEqual(cache.top.Data, 2)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, cache.bottom)
	aTest.MustBeEqual(cache.bottom.upperRecord, cache.top)
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_unlinkMiddleRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(3, 60)

	// Test #1.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 1,
			UID:  "First",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 2,
			UID:  "Second",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 3,
			UID:  "Third",
		},
	)
	var oldMiddle *FixedSizeBubbleCacheRecord = cache.top.lowerRecord
	cache.unlinkMiddleRecord(oldMiddle)
	//
	aTest.MustBeEqual(oldMiddle.UID, "Second")
	aTest.MustBeEqual(oldMiddle.Data, 2)
	aTest.MustBeEqual(oldMiddle.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(oldMiddle.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	//
	aTest.MustBeEqual(cache.top.UID, "Third")
	aTest.MustBeEqual(cache.top.Data, 3)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, cache.bottom)
	aTest.MustBeEqual(cache.bottom.upperRecord, cache.top)
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_unlinkBottomRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(3, 60)

	// Test #1.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 1,
			UID:  "First",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 2,
			UID:  "Second",
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			Data: 3,
			UID:  "Third",
		},
	)
	var oldBottom *FixedSizeBubbleCacheRecord = cache.bottom
	cache.unlinkBottomRecord()
	//
	aTest.MustBeEqual(oldBottom.UID, "First")
	aTest.MustBeEqual(oldBottom.Data, 1)
	aTest.MustBeEqual(oldBottom.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(oldBottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	//
	aTest.MustBeEqual(cache.top.UID, "Third")
	aTest.MustBeEqual(cache.top.Data, 3)
	aTest.MustBeEqual(cache.bottom.UID, "Second")
	aTest.MustBeEqual(cache.bottom.Data, 2)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, cache.bottom)
	aTest.MustBeEqual(cache.bottom.upperRecord, cache.top)
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_linkTopRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(3, 60)
	var record *FixedSizeBubbleCacheRecord

	// Test #1. Empty Cache.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "First",
		Data: 1,
	}
	cache.linkTopRecord(record)
	cache.size++
	//
	aTest.MustBeEqual(cache.top.UID, "First")
	aTest.MustBeEqual(cache.top.Data, 1)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))

	// Test #2. Non-empty Cache.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "Second",
		Data: 2,
	}
	cache.linkTopRecord(record)
	//
	aTest.MustBeEqual(cache.top.UID, "Second")
	aTest.MustBeEqual(cache.top.Data, 2)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, cache.bottom)
	aTest.MustBeEqual(cache.bottom.upperRecord, cache.top)
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_linkBottomRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var cache = NewFixedSizeBubbleCache(3, 60)
	var record *FixedSizeBubbleCacheRecord

	// Test #1. Empty Cache.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "First",
		Data: 1,
	}
	cache.linkBottomRecord(record)
	cache.size++
	//
	aTest.MustBeEqual(cache.top.UID, "First")
	aTest.MustBeEqual(cache.top.Data, 1)
	aTest.MustBeEqual(cache.bottom.UID, "First")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))

	// Test #2. Non-empty Cache.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "Second",
		Data: 2,
	}
	cache.linkBottomRecord(record)
	//
	aTest.MustBeEqual(cache.top.UID, "First")
	aTest.MustBeEqual(cache.top.Data, 1)
	aTest.MustBeEqual(cache.bottom.UID, "Second")
	aTest.MustBeEqual(cache.bottom.Data, 2)
	//
	aTest.MustBeEqual(cache.top.lowerRecord, cache.bottom)
	aTest.MustBeEqual(cache.bottom.upperRecord, cache.top)
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_Clear(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var err error

	// Test #1. Broken Cache.
	var cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.top.lowerRecord = nil
	err = cache.Clear()
	aTest.MustBeAnError(err)

	// Test #2. Normal Cache.
	cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	err = cache.Clear()
	aTest.MustBeNoError(err)
	//
	aTest.MustBeEqual(cache.size, uint(0))
	aTest.MustBeEqual(cache.capacity, uint(3))
	aTest.MustBeEqual(cache.recordTTL, uint(60))
	//
	aTest.MustBeEqual(cache.top, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_isIntegral(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	// Test #1. Bad Map.
	var cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.recordsByUID["1"] = nil
	aTest.MustBeEqual(cache.isIntegral(), false)

	// Test #2. Size > Capacity.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.size = 999
	aTest.MustBeEqual(cache.isIntegral(), false)

	// Test #2. Top and Bottom are non-null when Size is 0.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.size = 0
	aTest.MustBeEqual(cache.isIntegral(), false)
	cache.top = nil
	aTest.MustBeEqual(cache.isIntegral(), false)
	cache.bottom = nil
	aTest.MustBeEqual(cache.isIntegral(), true)

	// Test #3. Size is 1.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.top = nil
	aTest.MustBeEqual(cache.isIntegral(), false)
	//
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.bottom = nil
	aTest.MustBeEqual(cache.isIntegral(), false)
	//
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.top = &FixedSizeBubbleCacheRecord{
		UID:  "Fake",
		Data: 999,
	}
	aTest.MustBeEqual(cache.isIntegral(), false)
	//
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.top.upperRecord = &FixedSizeBubbleCacheRecord{
		UID:  "Fake",
		Data: 999,
	}
	aTest.MustBeEqual(cache.isIntegral(), false)
	//
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.bottom.lowerRecord = &FixedSizeBubbleCacheRecord{
		UID:  "Fake",
		Data: 999,
	}
	aTest.MustBeEqual(cache.isIntegral(), false)
	//
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	aTest.MustBeEqual(cache.isIntegral(), true)

	// Test #6. Edges.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)
	cache.top.upperRecord = &FixedSizeBubbleCacheRecord{
		UID:  "Fake",
		Data: 999,
	}
	aTest.MustBeEqual(cache.isIntegral(), false)
	cache.top.upperRecord = nil
	cache.bottom.lowerRecord = &FixedSizeBubbleCacheRecord{
		UID:  "Fake",
		Data: 999,
	}
	aTest.MustBeEqual(cache.isIntegral(), false)
	cache.bottom.lowerRecord = nil

	// Test #5. Forward Loop.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)
	var tmp = &FixedSizeBubbleCacheRecord{
		UID:         "Fake",
		Data:        999,
		lowerRecord: cache.bottom,
	}
	cache.bottom.upperRecord.lowerRecord = tmp
	aTest.MustBeEqual(cache.isIntegral(), false)
	tmp.lowerRecord = nil
	aTest.MustBeEqual(cache.isIntegral(), false)
	cache.top.lowerRecord = nil
	aTest.MustBeEqual(cache.isIntegral(), false)

	// Test #6. Reverse Loop.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)
	tmp = &FixedSizeBubbleCacheRecord{
		UID:         "Fake",
		Data:        999,
		upperRecord: cache.top,
	}
	cache.top.lowerRecord.upperRecord = tmp
	aTest.MustBeEqual(cache.isIntegral(), false)
	tmp.upperRecord = nil
	aTest.MustBeEqual(cache.isIntegral(), false)
	cache.bottom.upperRecord = nil
	aTest.MustBeEqual(cache.isIntegral(), false)

	// Test #7. OK.
	cache = NewFixedSizeBubbleCache(3, 1)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)
	aTest.MustBeEqual(cache.isIntegral(), true)
}

func Test_deleteRecord(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var err error

	// Test #1. Null Record.
	var cache = NewFixedSizeBubbleCache(3, 60)
	err = cache.deleteRecord(nil, false)
	aTest.MustBeAnError(err)

	// Test #2. Zero Size.
	err = cache.deleteRecord(&FixedSizeBubbleCacheRecord{}, false)
	aTest.MustBeAnError(err)

	// Test #3. Record is not known to exist and does not exist.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	err = cache.deleteRecord(&FixedSizeBubbleCacheRecord{}, false)
	aTest.MustBeAnError(err)

	// Test #4. Record exists.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	err = cache.deleteRecord(
		&FixedSizeBubbleCacheRecord{
			UID: "1",
		},
		false,
	)
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(cache.top, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.size, uint(0))

	// Test #5. Normal Deletion of a Top.
	cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	err = cache.deleteRecord(cache.top, false)
	aTest.MustBeNoError(err)
	//
	aTest.MustBeEqual(cache.top.UID, "1")
	aTest.MustBeEqual(cache.top.Data, 1)
	aTest.MustBeEqual(cache.bottom.UID, "1")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.top.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	//
	aTest.MustBeEqual(cache.size, uint(1))

	// Test #6. Normal Deletion of a Bottom.
	cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	err = cache.deleteRecord(cache.bottom, false)
	aTest.MustBeNoError(err)
	//
	aTest.MustBeEqual(cache.top.UID, "2")
	aTest.MustBeEqual(cache.top.Data, 2)
	aTest.MustBeEqual(cache.bottom.UID, "2")
	aTest.MustBeEqual(cache.bottom.Data, 2)
	//
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.top.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	//
	aTest.MustBeEqual(cache.size, uint(1))

	// Test #6. Normal Deletion of a Middle.
	cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)
	err = cache.deleteRecord(cache.top.lowerRecord, false)
	aTest.MustBeNoError(err)
	//
	aTest.MustBeEqual(cache.top.UID, "3")
	aTest.MustBeEqual(cache.top.Data, 3)
	aTest.MustBeEqual(cache.bottom.UID, "1")
	aTest.MustBeEqual(cache.bottom.Data, 1)
	//
	aTest.MustBeEqual(cache.top.upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.top.lowerRecord, cache.bottom)
	aTest.MustBeEqual(cache.bottom.upperRecord, cache.top)
	aTest.MustBeEqual(cache.bottom.lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
	//
	aTest.MustBeEqual(cache.size, uint(2))
}

func Test_RecordUIDExists(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	var cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)

	// Test #1.
	aTest.MustBeEqual(cache.RecordUIDExists("2"), true)
	aTest.MustBeEqual(cache.RecordUIDExists("8"), false)
}

func Test_recordUIDExists(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	var cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)

	// Test #1.
	aTest.MustBeEqual(cache.recordUIDExists("2"), true)
	aTest.MustBeEqual(cache.recordUIDExists("8"), false)
}

func Test_getRecordByUID(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	var cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	var now = uint(time.Now().Unix())
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)

	// Test #1.
	var record *FixedSizeBubbleCacheRecord
	var err error
	record, err = cache.getRecordByUID("2")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(record.UID, "2")
	aTest.MustBeEqual(record.Data, 2)
	aTest.MustBeEqual(record.lastAccessTime, now)
	aTest.MustBeEqual(record.upperRecord, cache.top)
	aTest.MustBeEqual(record.lowerRecord, cache.bottom)

	// Test #2.
	record, err = cache.getRecordByUID("999")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(
		err.Error(),
		fmt.Sprintf(ErrfRecordWithUidIsNotFound, "999"),
	)
}

func Test_DeleteRecordByUID(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	var cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)

	// Test #1.
	var err error
	err = cache.DeleteRecordByUID("999")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(
		err.Error(),
		fmt.Sprintf(ErrfRecordWithUidIsNotFound, "999"),
	)

	// Test #2. Error Emulation.
	cache.size = 0
	err = cache.DeleteRecordByUID("1")
	aTest.MustBeAnError(err)

	// Test #3.
	cache.size = 1
	err = cache.DeleteRecordByUID("1")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(cache.size, uint(0))
	aTest.MustBeEqual(cache.top, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_ListAllRecordValues(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	// Test #1. Non-empty Cache.
	var cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	aTest.MustBeEqual(
		cache.ListAllRecordValues(),
		[]interface{}{
			int(2),
			int(1),
		},
	)

	// Test #2. An empty Cache.
	cache = NewFixedSizeBubbleCache(3, 60)
	aTest.MustBeEqual(
		cache.ListAllRecordValues(),
		[]interface{}{},
	)
}

func Test_ListAllRecords(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	// Test #1. An empty Cache.
	var cache = NewFixedSizeBubbleCache(3, 60)
	aTest.MustBeEqual(
		cache.ListAllRecords(),
		[]*FixedSizeBubbleCacheRecord{},
	)

	// Test #2. Non-empty Cache.
	cache = NewFixedSizeBubbleCache(3, 60)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "3",
			Data: 3,
		},
	)
	var records []*FixedSizeBubbleCacheRecord = cache.ListAllRecords()
	aTest.MustBeEqual(len(records), 3)
	//
	aTest.MustBeEqual(records[0].UID, "3")
	aTest.MustBeEqual(records[0].Data, 3)
	aTest.MustBeEqual(records[0].upperRecord, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(records[0].lowerRecord, cache.top.lowerRecord)
	//
	aTest.MustBeEqual(records[1].UID, "2")
	aTest.MustBeEqual(records[1].Data, 2)
	aTest.MustBeEqual(records[1].upperRecord, cache.top)
	aTest.MustBeEqual(records[1].lowerRecord, cache.bottom)
	//
	aTest.MustBeEqual(records[2].UID, "1")
	aTest.MustBeEqual(records[2].Data, 1)
	aTest.MustBeEqual(records[2].upperRecord, cache.bottom.upperRecord)
	aTest.MustBeEqual(records[2].lowerRecord, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_GetActualRecordDataByUID(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var data interface{}
	var err error

	// Test #1. Non-existent Record.
	var cache = NewFixedSizeBubbleCache(3, 2)
	data, err = cache.GetActualRecordDataByUID("999")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(data, nil)

	// Test #2. An existent Record, a Top.
	var now = uint(time.Now().Unix())
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	time.Sleep(time.Second)
	data, err = cache.GetActualRecordDataByUID("1")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(data, 1)
	aTest.MustBeEqual(cache.top.lastAccessTime, now+1)

	// Test #3. An existent Record, not a Top.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "2",
			Data: 2,
		},
	)
	data, err = cache.GetActualRecordDataByUID("1")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(data, 1)
	aTest.MustBeEqual(cache.top.lastAccessTime, now+1)

	// Test #4. An existent Record, outdated.
	cache = NewFixedSizeBubbleCache(3, 2)
	now = uint(time.Now().Unix())
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	time.Sleep(time.Second * 3)
	data, err = cache.GetActualRecordDataByUID("1")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), fmt.Sprintf(ErrfRecordWithUidIsOutdated, "1"))
	aTest.MustBeEqual(cache.size, uint(0))
	aTest.MustBeEqual(cache.top, (*FixedSizeBubbleCacheRecord)(nil))
	aTest.MustBeEqual(cache.bottom, (*FixedSizeBubbleCacheRecord)(nil))
}

func Test_GetRecordTTL(t *testing.T) {
	var aTest *tester.Test = tester.New(t)

	// Test #1.
	var cache = NewFixedSizeBubbleCache(3, 60)
	aTest.MustBeEqual(cache.GetRecordTTL(), uint(60))
}

func Test_IsRecordUIDActive(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var recordIsActive bool
	var err error

	// Test #1. Non-existent Record.
	var cache = NewFixedSizeBubbleCache(3, 1)
	recordIsActive, err = cache.IsRecordUIDActive("999")
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), fmt.Sprintf(ErrfRecordWithUidIsNotFound, "999"))

	// Test #2. An existent Record.
	cache.addRecord(
		&FixedSizeBubbleCacheRecord{
			UID:  "1",
			Data: 1,
		},
	)
	time.Sleep(time.Second * 2)
	recordIsActive, err = cache.IsRecordUIDActive("1")
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(recordIsActive, false)
}
