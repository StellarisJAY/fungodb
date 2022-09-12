package engine

import "math/rand"

type node struct {
	key    string
	value  []byte
	prev   *node
	levels []*level
}

type level struct {
	next *node
	key  string
}

type SkipList struct {
	head         *node
	length       int
	maxLevel     int
	currentLevel int
}

func NewSkipList(maxLevel int) *SkipList {
	sl := &SkipList{&node{}, 0, maxLevel, 0}
	sl.head.levels = make([]*level, maxLevel)
	for i := 0; i < maxLevel; i++ {
		sl.head.levels[i] = new(level)
	}
	return sl
}

func newNode(key string, value []byte, maxLevels int) *node {
	countLevels := randomLevel(maxLevels)
	n := &node{key: key, value: value, levels: make([]*level, countLevels)}
	for i := 0; i < countLevels; i++ {
		n.levels[i] = &level{key: key}
	}
	return n
}

func randomLevel(maxLevels int) int {
	level := 1
	for float32(rand.Int31()&0xFFFF) < (0.25 * 0xFFFF) {
		level++
	}
	if level < maxLevels {
		return level
	}
	return maxLevels
}

func (l *SkipList) Put(key string, value []byte) {
	n := l.head
	// 遍历跳表路过的节点，这些节点在插入新节点时作为新节点的前驱
	passedNodes := make([]*node, l.currentLevel)
	// 从最顶层开始遍历，直到0层
	for i := l.currentLevel - 1; i >= 0; i-- {
		// 每一层判断下一个节点的key大小，小于当前key则可以跳过
		for n.levels[i].next != nil && n.levels[i].next.key < key {
			n = n.levels[i].next
		}
		// 记录路过的节点
		passedNodes[i] = n
		// 如果n的key与key相同，表示key已经存在，是更新操作
		if n.levels[i].next != nil && n.levels[i].next.key == key {
			n = n.levels[i].next
			break
		}
	}
	if n.key == key {
		n.value = value
		return
	}

	node := newNode(key, value, l.maxLevel)
	if len(node.levels) > l.currentLevel {
		// 从0层开始，每一层插入节点
		for i := 0; i < l.currentLevel; i++ {
			// 新节点在这一层的下一个，是遍历时路过这一层的最后一个节点的下一个节点
			node.levels[i].next = passedNodes[i].levels[i].next
			passedNodes[i].levels[i].next = node
		}
		// 新节点的层数超过了当前跳表的最大层数，这些部分的前驱节点是head节点
		for i := l.currentLevel; i < len(node.levels); i++ {
			l.head.levels[i].next = node
		}
		l.currentLevel = len(node.levels)
	} else {
		for i := 0; i < len(node.levels); i++ {
			node.levels[i].next = passedNodes[i].levels[i].next
			passedNodes[i].levels[i].next = node
		}
	}
	l.length++
}

func (l *SkipList) Remove(key string) {
	n := l.head
	passedNodes := make([]*node, l.currentLevel)
	for i := l.currentLevel - 1; i >= 0; i-- {
		for n.levels[i].next != nil && n.levels[i].next.key < key {
			n = n.levels[i].next
		}
		passedNodes[i] = n
	}
	if n.levels[0].next == nil || n.levels[0].next.key != key {
		return
	}
	remove := n.levels[0].next
	for i := 0; i < len(remove.levels); i++ {
		passedNodes[i].levels[i].next = remove.levels[i].next
		remove.levels[i].next = nil
	}
	for l.currentLevel > 0 && l.head.levels[l.currentLevel-1].next == nil {
		l.currentLevel--
	}
	l.length--
}
