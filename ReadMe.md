# Fixed Size Bubble Cache.


## Short Description.

This Package provides a fixed Size Bubble Cache Functionality.

## Full Description.

The Cache Stores Information about the N most active Records, where 'N' is a 
fixed Number of Records. New Records are placed at the Top, old Records are 
removed from the Bottom of the Cache. Request for an existing Record moves the 
Record to the Top Position.

The Cache Object has two Parameters:

	*	Capacity, Maximum Size (N, mentioned above);
	*	Time-to-Live Settings of a Record (Period is set in Seconds).

When we add a Record to the Cache, if an incoming Record already exists in the 
Cache, it is moved from its existing Position to the Top of the Cache. The Term 
'exists' means that there is a Record in the Cache with the same UID as the UID 
of the inserted Record.

Each Record has a 'UID' and a 'Data' Field.
'UID' is used for Indexing. 'Data' is used to store some useful Information.

If an incoming Record is new (does not exist in the Cache), it is added to the 
Top of the Cache. If the Cache is already at its maximum Size, then the oldest 
Record, which is located at the Bottom of the Cache, is removed. 

For Example, if the Size of the Cache (N) is Five (5), then the following 
Examples are correct: <br />
[ghi] + [abc,def,ghi,jkl,xyz] => [ghi,abc,def,jkl,xyz]. <br />
[xxx] + [abc,def,ghi,jkl,xyz] => [xxx,abc,def,ghi,jkl]. <br />

When a User requests a Value (by its UID) from the Cache, we first, check its 
Existence in the Cache's List, and then we check the Record's TTL (Time To 
Live). If the requested Record exists but is outdated, we remove it from the 
Cache.

The Removals are done in a "Lazy" Style: either when the Record is requested, or
when a new Record arrives and we have no free Space to store old Records. This 
is done to save much of the CPU Time. We check TTL only when it is necessary.

## Installation.

Import Commands:
```
go get -u "github.com/vault-thirteen/FixedSizeBubbleCache"
```

## Usage.

```
import "github.com/vault-thirteen/FixedSizeBubbleCache"
```
