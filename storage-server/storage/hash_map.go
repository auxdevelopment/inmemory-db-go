package storage

import (
	"hash/fnv"
	"sync"
)

type HashMap struct {
	buckets           []*LinkedList
	loadFactor        float32
	numberStoredItems int32
	availableBuckets  int32
	rwLock            sync.RWMutex
}

func NewHashMap() *HashMap {
	capacity := int32(16)

	buckets := make([]*LinkedList, capacity)

	for i := int32(0); i < capacity; i++ {
		buckets[i] = NewLinkedList()
	}

	return &HashMap{
		buckets:           buckets,
		availableBuckets:  capacity,
		loadFactor:        0.75,
		numberStoredItems: 0,
		rwLock:            sync.RWMutex{},
	}
}

func (hashMap *HashMap) GetAll() ([]*KeyValuePair, error) {
	var pairs []*KeyValuePair

	hashMap.rwLock.RLock()
	defer hashMap.rwLock.RUnlock()

	for i := int32(0); i < hashMap.availableBuckets; i++ {
		current := hashMap.buckets[i].Head

		for current != nil {
			pairs = append(pairs, &current.Pair)

			current = current.Next
		}
	}

	return pairs, nil
}

func (hashMap *HashMap) Get(key string) (*KeyValuePair, error) {

	index := hashMap.getIndexForKey(key)

	hashMap.rwLock.RLock()
	defer hashMap.rwLock.RUnlock()

	bucket := hashMap.buckets[index]

	return bucket.Get(key)
}

func (hashMap *HashMap) putPair(pair KeyValuePair) error {
	index := hashMap.getIndexForKey(pair.Key)
	bucket := hashMap.buckets[index]

	fetchedPair, err := bucket.Get(pair.Key)

	if err != nil {
		addError := bucket.Add(pair)

		if addError != nil {
			return addError
		}

		hashMap.numberStoredItems++

		return nil
	}

	fetchedPair.Value = pair.Value
	hashMap.numberStoredItems++

	return nil
}

func (hashMap *HashMap) Put(pair KeyValuePair) error {
	hashMap.rwLock.Lock()
	defer hashMap.rwLock.Unlock()

	if hashMap.needsRehash() {
		err := hashMap.rehash()

		if err != nil {
			return err
		}
	}

	return hashMap.putPair(pair)
}

func (hashMap *HashMap) Remove(key string) bool {

	index := hashMap.getIndexForKey(key)

	hashMap.rwLock.Lock()
	defer hashMap.rwLock.Unlock()

	bucket := hashMap.buckets[index]

	if bucket.RemoveByKey(key) {
		hashMap.numberStoredItems--

		return true
	}

	return false
}

func (hashMap *HashMap) getIndexForKey(key string) int {
	algo := fnv.New32()
	algo.Write([]byte(key))

	return int(algo.Sum32()) % len(hashMap.buckets)
}

func (hashMap *HashMap) needsRehash() bool {
	max := hashMap.loadFactor * float32(hashMap.availableBuckets)

	return float32(hashMap.numberStoredItems) > max
}

func (hashMap *HashMap) rehash() error {
	// Remove all items from map and store them temporarily
	allEntries := make([]KeyValuePair, hashMap.numberStoredItems)
	entriesIndex := 0

	for bucketIndex := int32(0); bucketIndex < hashMap.availableBuckets; bucketIndex++ {
		bucket := hashMap.buckets[bucketIndex]

		if bucket.Head == nil {
			continue
		}

		bucketSize := bucket.Size()

		for bucketItemIndex := 0; bucketItemIndex < bucketSize; bucketItemIndex++ {
			removed, err := bucket.pop()
			hashMap.numberStoredItems--

			if err != nil {
				return err
			}

			allEntries[entriesIndex] = removed
			entriesIndex++
		}
	}

	// Resize internal array (double current size)

	hashMap.availableBuckets *= 2
	hashMap.buckets = make([]*LinkedList, hashMap.availableBuckets)

	for i := int32(0); i < hashMap.availableBuckets; i++ {
		hashMap.buckets[i] = NewLinkedList()
	}

	// Redistribute previously stored KeyValuePairs
	for entriesIndex = 0; entriesIndex < len(allEntries); entriesIndex++ {
		if err := hashMap.putPair(allEntries[entriesIndex]); err != nil {
			return err
		}
	}

	return nil
}
