package syncmap

import "sync"

// Map - Sync map.
type Map interface {
	// Lock - Locks the entire map.
	Lock()

	// Unlock - Unlocks the map for access.
	Unlock()

	// Get - Returns value from map by key.
	// Constructs value if not found and returns.
	// Map is locked during construction of the value.
	Get(interface{}) (interface{}, error)

	// Set - Sets value by key.
	Set(interface{}, interface{})

	// Delete - Deletes value by it's key.
	Delete(interface{})

	// Pop - Pops value by it's key.
	// Finds, deletes from map and returns.
	// It does not create a value if one was not found.
	// Bool parameter is true if value was found in a map.
	Pop(interface{}) (interface{}, bool)

	// Values - Returns map values.
	// It is lock-free meaning a map should be locked first.
	Values() map[interface{}]interface{}

	// Clear - Clears the map.
	Clear()
}

// Func - Sync map func.
type Func func(interface{}) (interface{}, error)

// New - Creates a new sync map.
func New(fnc Func) Map {
	return &syncmap{
		fnc:    fnc,
		mutex:  new(sync.Mutex),
		values: make(map[interface{}]interface{}),
	}
}

type syncmap struct {
	fnc    Func
	mutex  *sync.Mutex
	values map[interface{}]interface{}
}

func (m *syncmap) Get(key interface{}) (value interface{}, err error) {
	m.mutex.Lock()
	value = m.values[key]
	if value == nil {
		value, err = m.fnc(key)
		if err == nil {
			m.values[key] = value
		}
	}
	m.mutex.Unlock()
	return
}

func (m *syncmap) Set(key interface{}, value interface{}) {
	m.mutex.Lock()
	m.values[key] = value
	m.mutex.Unlock()
}

func (m *syncmap) Delete(key interface{}) {
	m.mutex.Lock()
	delete(m.values, key)
	m.mutex.Unlock()
}

func (m *syncmap) Pop(key interface{}) (value interface{}, ok bool) {
	m.mutex.Lock()
	value, ok = m.values[key]
	if ok {
		delete(m.values, key)
	}
	m.mutex.Unlock()
	return
}

func (m *syncmap) Clear() {
	m.mutex.Lock()
	m.values = make(map[interface{}]interface{})
	m.mutex.Unlock()
}

func (m *syncmap) Lock() {
	m.mutex.Lock()
}

func (m *syncmap) Unlock() {
	m.mutex.Unlock()
}

func (m *syncmap) Values() map[interface{}]interface{} {
	return m.values
}
