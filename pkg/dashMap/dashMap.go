package dashMap

import (
	fb "github.com/Eugene-Usachev/fastbytes"
	"hash/maphash"
	"sync"
)

type DashMap[V any] struct {
	Lockers   []sync.RWMutex
	Maps      []map[string]V
	Size      uint64
	Seed      maphash.Seed
	HashFunc  func(key string) uint64
	IndexFunc func(key string) uint64
}

func NewDashMap[V any](size uint64) *DashMap[V] {
	self := &DashMap[V]{
		Lockers: make([]sync.RWMutex, size),
		Maps:    make([]map[string]V, size),
		Size:    size,
		Seed:    maphash.MakeSeed(),
	}
	for i := range self.Maps {
		self.Maps[i] = make(map[string]V)
	}
	self.HashFunc = func(key string) uint64 {
		return maphash.Bytes(self.Seed, fb.S2B(key))
	}
	self.IndexFunc = func(key string) uint64 {
		return maphash.Bytes(self.Seed, fb.S2B(key)) % size
	}

	return self
}

func (d *DashMap[V]) Get(key string) V {
	index := maphash.Bytes(d.Seed, fb.S2B(key)) % d.Size
	d.Lockers[index].RLock()
	v := d.Maps[index][key]
	d.Lockers[index].RUnlock()
	return v
}

func (d *DashMap[V]) Set(key string, value V) {
	index := maphash.Bytes(d.Seed, fb.S2B(key)) % d.Size
	d.Lockers[index].Lock()
	d.Maps[index][key] = value
	d.Lockers[index].Unlock()
}

func (d *DashMap[V]) Delete(key string) {
	index := maphash.Bytes(d.Seed, fb.S2B(key)) % d.Size
	d.Lockers[index].Lock()
	delete(d.Maps[index], key)
	d.Lockers[index].Unlock()
}

func (d *DashMap[V]) Len() int {
	var l int
	for i := 0; i < len(d.Maps); i++ {
		d.Lockers[i].RLock()
		l += len(d.Maps[i])
		d.Lockers[i].RUnlock()
	}
	return l
}

func (d *DashMap[V]) Clear() {
	for i := 0; i < len(d.Maps); i++ {
		d.Lockers[i].Lock()
		d.Maps[i] = make(map[string]V)
		d.Lockers[i].Unlock()
	}
}

func (d *DashMap[V]) ForEach(f func(key string, value V)) {
	for i := 0; i < len(d.Maps); i++ {
		d.Lockers[i].RLock()
		for k, v := range d.Maps[i] {
			f(k, v)
		}
		d.Lockers[i].RUnlock()
	}
}
