package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
	Values() [][]any
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*Node
}

type cacheItem struct {
	key   Key
	value any
}

func NewCache(capacity int) *lruCache {
	if capacity < 1 {
		return nil
	}
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Node, capacity),
	}
}

func (cache *lruCache) Set(key Key, value any) bool {
	node, ok := cache.items[key]

	if !ok && cache.queue.Len() >= cache.capacity {
		delete(cache.items, cache.queue.Back().Value.(cacheItem).key)
		cache.queue.Remove(cache.queue.Back())
	}

	if ok {
		cache.queue.MoveToFront(node)
		if node.Value.(cacheItem).value != value {
			node.Value = cacheItem{key: key, value: value}
		}
	} else {
		cache.items[key] = cache.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
	}

	return ok
}

func (cache *lruCache) Get(key Key) (any, bool) {
	node, ok := cache.items[key]

	var value any

	if ok {
		cache.queue.MoveToFront(node)
		value = node.Value.(cacheItem).value
	}

	return value, ok
}

func (cache *lruCache) Clear() {
	cache.items = make(map[Key]*Node, cache.capacity)
	cache.queue = NewList()
}

func (cache *lruCache) Values() [][]any {
	result := make([][]any, 0, len(cache.items))
	currentNode := cache.queue.Front()

	for currentNode != nil {
		item := currentNode.Value.(cacheItem)
		result = append(result, []any{item.key, item.value})
		currentNode = currentNode.Next
	}

	return result
}
