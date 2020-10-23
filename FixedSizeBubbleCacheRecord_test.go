package fsbcache

import (
	"time"

	"github.com/vault-thirteen/tester"

	"testing"
)

func Test_Check(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var record *FixedSizeBubbleCacheRecord
	var err error

	// Test #1. Null Data.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "123",
		Data: nil,
	}
	err = record.Check()
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrDataIsEmpty)

	// Test #2. Empty UID.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "",
		Data: 123,
	}
	err = record.Check()
	aTest.MustBeAnError(err)
	aTest.MustBeEqual(err.Error(), ErrUIDIsEmpty)

	// Test #3. Normal.
	record = &FixedSizeBubbleCacheRecord{
		UID:  "123",
		Data: 123,
	}
	err = record.Check()
	aTest.MustBeNoError(err)
}

func Test_UpdateDataAndLAT(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var record *FixedSizeBubbleCacheRecord

	// Test #1.
	var tsNow = uint(time.Now().Unix())
	record = &FixedSizeBubbleCacheRecord{
		Data:           101,
		lastAccessTime: tsNow,
	}
	time.Sleep(time.Second)
	record.UpdateDataAndLAT(202)
	aTest.MustBeEqual(record.Data, 202)
	aTest.MustBeEqual(record.lastAccessTime-tsNow, uint(1))
}

func Test_UpdateLAT(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var record *FixedSizeBubbleCacheRecord

	// Test #1.
	var tsNow = uint(time.Now().Unix())
	record = &FixedSizeBubbleCacheRecord{
		lastAccessTime: tsNow,
	}
	time.Sleep(time.Second)
	record.UpdateLAT()
	aTest.MustBeEqual(record.lastAccessTime-tsNow, uint(1))
}

func Test_updateLATWithCurrentTime(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var record *FixedSizeBubbleCacheRecord

	// Test #1.
	var tsNow = uint(time.Now().Unix())
	record = &FixedSizeBubbleCacheRecord{
		lastAccessTime: tsNow,
	}
	time.Sleep(time.Second)
	record.updateLATWithCurrentTime()
	aTest.MustBeEqual(record.lastAccessTime-tsNow, uint(1))
}

func Test_isActual(t *testing.T) {
	var aTest *tester.Test = tester.New(t)
	var record *FixedSizeBubbleCacheRecord

	// Test #1. Positive.
	record = &FixedSizeBubbleCacheRecord{}
	record.UpdateLAT()
	time.Sleep(time.Second * 2)
	aTest.MustBeEqual(record.isActual(3), true)
	aTest.MustBeEqual(record.isActual(1), false)
}
