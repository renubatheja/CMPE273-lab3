package consistentHashing

/**
A consistent hashing algo for REST servers 
*/

import (
	"fmt"
	_ "math"
	"sort"
)


type HashKey uint32
type HashKeyOrder []HashKey

func (h HashKeyOrder) Len() int           { return len(h) }
func (h HashKeyOrder) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h HashKeyOrder) Less(i, j int) bool { return h[i] < h[j] }


type HashRing struct {
	ring       map[HashKey]string
	sortedKeys []HashKey
	nodes      []string
}


func New(nodes []string) *HashRing {
	hashRing := &HashRing{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
	}
	hashRing.generateCircle()
	return hashRing
}


func (h *HashRing) generateCircle() {
	j := 0
	for _, node := range h.nodes {
			nodeKey := fmt.Sprintf("%s-%d", node, j)
			bKey := HashMD5(nodeKey)
			j++
			for i := 0; i < 3; i++ {
				key := hashVal(bKey[i*4 : i*4+4])
				h.ring[key] = node
				h.sortedKeys = append(h.sortedKeys, key)
			}
	}

	sort.Sort(HashKeyOrder(h.sortedKeys))
}


func (h *HashRing) GetNode(stringKey string) (node string, ok bool) {
	pos, ok := h.GetNodePos(stringKey)
	if !ok {
		return "", false
	}
	return h.ring[h.sortedKeys[pos]], true
}


func (h *HashRing) GetNodePos(stringKey string) (pos int, ok bool) {
	if len(h.ring) == 0 {
		return 0, false
	}

	key := h.GenKey(stringKey)

	nodes := h.sortedKeys
	pos = sort.Search(len(nodes), func(i int) bool { return nodes[i] > key })

	if pos == len(nodes) {
		// Wrap the search, should return first node
		return 0, true
	} else {
		return pos, true
	}
}


func (h *HashRing) GenKey(key string) HashKey {
	bKey := HashMD5(key)
	return hashVal(bKey[0:4])
}


func (h *HashRing) AddNode(node string) *HashRing {
	nodes := make([]string, len(h.nodes), len(h.nodes)+1)
	copy(nodes, h.nodes)
	nodes = append(nodes, node)

	hashRing := &HashRing{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
	}
	hashRing.generateCircle()
	return hashRing
	
}


func (h *HashRing) RemoveNode(node string) *HashRing {
	nodes := make([]string, 0)
	for _, eNode := range h.nodes {
		if eNode != node {
			nodes = append(nodes, eNode)
		}
	}

	hashRing := &HashRing{
		ring:       make(map[HashKey]string),
		sortedKeys: make([]HashKey, 0),
		nodes:      nodes,
	}
	hashRing.generateCircle()
	return hashRing
}


func hashVal(bKey []byte) HashKey {
	return ((HashKey(bKey[3]) << 24) |
		(HashKey(bKey[2]) << 16) |
		(HashKey(bKey[1]) << 8) |
		(HashKey(bKey[0])))
}


