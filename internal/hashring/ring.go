package hashring

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

type Ring interface {
	AddNode(nodeID string)

	RemoveNode(nodeID string)

	GetNode(key string) string

	GetNodes() []string

	Size() int
}

type HashRing struct {
	mu           sync.RWMutex
	ring         []uint32
	nodeMap      map[uint32]string
	nodes        map[string]bool
	virtualNodes int
}

func New(virtualNodes int) *HashRing {
	if virtualNodes <= 0 {
		virtualNodes = 150
	}
	return &HashRing{
		ring:         make([]uint32, 0),
		nodeMap:      make(map[uint32]string),
		nodes:        make(map[string]bool),
		virtualNodes: virtualNodes,
	}
}

func (hr *HashRing) AddNode(nodeID string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	if hr.nodes[nodeID] {
		return
	}

	hr.nodes[nodeID] = true

	for i := 0; i < hr.virtualNodes; i++ {
		vKey := fmt.Sprintf("%s#%d", nodeID, i)
		h := hashKey(vKey)
		hr.ring = append(hr.ring, h)
		hr.nodeMap[h] = nodeID
	}

	sort.Slice(hr.ring, func(i, j int) bool {
		return hr.ring[i] < hr.ring[j]
	})
}

func (hr *HashRing) RemoveNode(nodeID string) {
	hr.mu.Lock()
	defer hr.mu.Unlock()

	if !hr.nodes[nodeID] {
		return
	}

	delete(hr.nodes, nodeID)

	newRing := make([]uint32, 0, len(hr.ring))
	for _, h := range hr.ring {
		if hr.nodeMap[h] != nodeID {
			newRing = append(newRing, h)
		} else {
			delete(hr.nodeMap, h)
		}
	}
	hr.ring = newRing
}

func (hr *HashRing) GetNode(key string) string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	if len(hr.ring) == 0 {
		return ""
	}

	h := hashKey(key)

	idx := sort.Search(len(hr.ring), func(i int) bool {
		return hr.ring[i] >= h
	})

	if idx >= len(hr.ring) {
		idx = 0
	}

	return hr.nodeMap[hr.ring[idx]]
}

func (hr *HashRing) GetNodes() []string {
	hr.mu.RLock()
	defer hr.mu.RUnlock()

	nodes := make([]string, 0, len(hr.nodes))
	for id := range hr.nodes {
		nodes = append(nodes, id)
	}
	sort.Strings(nodes)
	return nodes
}

func (hr *HashRing) Size() int {
	hr.mu.RLock()
	defer hr.mu.RUnlock()
	return len(hr.nodes)
}

func hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}
