package concurrent

import (
	"sync/atomic"
	"unsafe"
)

//lock free queue, based on the implement of golang package
//only support one producer to push head, multi consumers to pop tail

type queue struct {
	headTail uint64
	vals     []eface
}

type eface struct {
	typ, val unsafe.Pointer
}

const dequeueBits = 32
const dequeueLimit = (1 << dequeueBits) / 4

type dequeueNil *struct{}

func (d *queue) unpack(ptrs uint64) (head, tail uint32) {
	const mask = 1<<dequeueBits - 1
	head = uint32((ptrs >> dequeueBits) & mask)
	tail = uint32(ptrs & mask)
	return
}

func (d *queue) pack(head, tail uint32) uint64 {
	const mask = 1<<dequeueBits - 1
	return (uint64(head) << dequeueBits) |
		uint64(tail&mask)
}

func (d *queue) pushHead(val interface{}) bool {
	ptrs := atomic.LoadUint64(&d.headTail)
	head, tail := d.unpack(ptrs)
	if (tail+uint32(len(d.vals)))&(1<<dequeueBits-1) == head {
		// Queue is full.
		return false
	}
	slot := &d.vals[head&uint32(len(d.vals)-1)]

	typ := atomic.LoadPointer(&slot.typ)
	if typ != nil {
		return false
	}

	if val == nil {
		val = dequeueNil(nil)
	}
	*(*interface{})(unsafe.Pointer(slot)) = val
	atomic.AddUint64(&d.headTail, 1<<dequeueBits)
	return true
}

func (d *queue) popHead() (interface{}, bool) {
	var slot *eface
	for {
		ptrs := atomic.LoadUint64(&d.headTail)
		head, tail := d.unpack(ptrs)
		if tail == head {
			// Queue is empty.
			return nil, false
		}
		head--
		ptrs2 := d.pack(head, tail)
		if atomic.CompareAndSwapUint64(&d.headTail, ptrs, ptrs2) {
			// We successfully took back slot.
			slot = &d.vals[head&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*interface{})(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}
	*slot = eface{}
	return val, true
}

func (d *queue) popTail() (interface{}, bool) {
	var slot *eface
	for {
		ptrs := atomic.LoadUint64(&d.headTail)
		head, tail := d.unpack(ptrs)
		if tail == head {
			// Queue is empty.
			return nil, false
		}

		ptrs2 := d.pack(head, tail+1)
		if atomic.CompareAndSwapUint64(&d.headTail, ptrs, ptrs2) {
			// Success.
			slot = &d.vals[tail&uint32(len(d.vals)-1)]
			break
		}
	}

	val := *(*interface{})(unsafe.Pointer(slot))
	if val == dequeueNil(nil) {
		val = nil
	}

	slot.val = nil
	atomic.StorePointer(&slot.typ, nil)
	// At this point pushHead owns the slot.

	return val, true
}

func NewLightQueue() *LightQueue {
	return new(LightQueue)
}

func (c *LightQueue) PushHead(val interface{}) bool {
	c.pushHead(val)
	return true
}

func (c *LightQueue) PopHead() (interface{}, bool) {
	return c.popHead()
}

func (c *LightQueue) PopTail() (interface{}, bool) {
	return c.popTail()
}

type LightQueue struct {
	head *lightQueueElt
	tail *lightQueueElt
}

type lightQueueElt struct {
	queue
	next, prev *lightQueueElt
}

func storePoolChainElt(pp **lightQueueElt, v *lightQueueElt) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(pp)), unsafe.Pointer(v))
}

func loadLightQueueElt(pp **lightQueueElt) *lightQueueElt {
	return (*lightQueueElt)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(pp))))
}

func (c *LightQueue) pushHead(val interface{}) {
	d := c.head
	if d == nil {
		// Initialize the chain.
		const initSize = 8 // Must be a power of 2
		d = new(lightQueueElt)
		d.vals = make([]eface, initSize)
		c.head = d
		storePoolChainElt(&c.tail, d)
	}

	if d.pushHead(val) {
		return
	}

	// The current dequeue is full. Allocate a new one of twice
	// the size.
	newSize := len(d.vals) * 2
	if newSize >= dequeueLimit {
		// Can't make it any bigger.
		newSize = dequeueLimit
	}

	d2 := &lightQueueElt{prev: d}
	d2.vals = make([]eface, newSize)
	c.head = d2
	storePoolChainElt(&d.next, d2)
	d2.pushHead(val)
}

func (c *LightQueue) popHead() (interface{}, bool) {
	d := c.head
	for d != nil {
		if val, ok := d.popHead(); ok {
			return val, ok
		}
		d = loadLightQueueElt(&d.prev)
	}
	return nil, false
}

func (c *LightQueue) popTail() (interface{}, bool) {
	d := loadLightQueueElt(&c.tail)
	if d == nil {
		return nil, false
	}
	for {
		d2 := loadLightQueueElt(&d.next)

		if val, ok := d.popTail(); ok {
			return val, ok
		}

		if d2 == nil {
			return nil, false
		}
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.tail)), unsafe.Pointer(d), unsafe.Pointer(d2)) {
			storePoolChainElt(&d2.prev, nil)
		}
		d = d2
	}
}
