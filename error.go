// Fixed Size Bubble Cache.

package fsbcache

// Error Messages.
const (
	ErrDataIsEmpty    = `'Data' Field is not set`
	ErrRecordIsNotSet = `Record is not set`
	ErrUIDIsEmpty     = `'UID' Field is not set`
	ErrCacheZeroSize  = "Cache Size is Zero"
	//
	ErrfRecordWithUidIsNotFound = `Record with UID='%v' is not found`
	ErrfRecordWithUidIsOutdated = `Record with UID='%v' is outdated`
	ErrIntegrityCheckFailure    = `Integrity Check Failure`
	//
	ErrTypeCast = "Type Cast Failure"
)
