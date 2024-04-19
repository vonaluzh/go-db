package main

import (
	"fmt"
	"goDb/connection"
	"sort"
	"strconv"
	"strings"
)

const MaxDegree = 3

type Node struct {
	parent   *Node
	keys     []int
	children []*Node
}

func (node *Node) getRoot() *Node {
	var currNode = node

	for {
		if currNode.parent == nil {
			return currNode
		}

		currNode = currNode.parent
	}
}

func (node *Node) findSuitableLeafNode(key int) *Node {
	if len(node.children) < 1 {
		return node
	}

	var keysLength = len(node.keys)
	var childNode *Node

	for i := 0; i < keysLength; i++ {
		var currKey = node.keys[i]

		if key <= currKey {
			childNode = node.children[i]
			break
		} else if i == keysLength-1 {
			childNode = node.children[i+1]
		}
	}

	return childNode.findSuitableLeafNode(key)
}

func (node *Node) insert(key int) {
	node.getRoot().findSuitableLeafNode(key).doInsert(key)
}

func (node *Node) doInsert(key int) {
	if len(node.keys) < 2 {
		node.keys = append(node.keys, key)
		sort.Ints(node.keys)
		return
	}
	// splitting node

	node.keys = append(node.keys, key)
	sort.Ints(node.keys)

	var rightKeys []int
	var middleKey = node.keys[1]
	node.keys = node.keys[:1]

	var parentNode *Node
	var newSiblingNode *Node

	if len(node.children) < 1 { // scenario 1 "splitting leaf node"
		rightKeys = node.keys[1:3]
		newSiblingNode = &Node{keys: rightKeys}

		if node.parent == nil {
			parentNode = &Node{}
			parentNode.children = append(parentNode.children, node, newSiblingNode)
			node.parent = parentNode
		} else {
			parentNode = node.parent
			parentNode.children = append(parentNode.children, newSiblingNode)
		}

		newSiblingNode.parent = parentNode
	} else { // scenario 2 "splitting middle node"
		rightKeys = node.keys[2:3]
		newSiblingNode = &Node{keys: rightKeys}

		if node.parent == nil {
			parentNode = &Node{}
			parentNode.children = append(parentNode.children, node, newSiblingNode)
			node.parent = parentNode
		} else {
			parentNode = node.parent
			parentNode.children = append(parentNode.children, newSiblingNode)
		}

		newSiblingNode.parent = parentNode
		newSiblingNode.children = node.children[2:4]
		// Can we update pointers smarter?
		for _, children := range newSiblingNode.children {
			children.parent = newSiblingNode
		}

		node.children = node.children[:2]
	}

	parentNode.doInsert(middleKey)
}

func main() {
	var node = Node{}

	for i := 1; i < 11; i++ {
		node.insert(i)
	}

	var rootNode = node.getRoot()

	draw(rootNode, 0)

	connection.UpServer()
}

func draw(node *Node, level int) {
	var nodeKeys []string

	for _, key := range node.keys {
		nodeKeys = append(nodeKeys, strconv.Itoa(key))
	}

	fmt.Printf("%s|--[%s]\n", strings.Repeat("  ", level), strings.Join(nodeKeys, ", "))

	if len(node.children) > 0 {
		for _, child := range node.children {
			draw(child, level+1)
		}
	}
}
