// The Computer Language Benchmarks Game
// http://benchmarksgame.alioth.debian.org/
//
//
// Alexandr Karbivnichiy

package main

import (
	"flag"
	"fmt"
	"sort"
	"strconv"
	"sync"
)

type Node struct {
	Left, Right *Node
}

type Tree struct {
	chunks [][]Node
}

const lots = 64

// similar to typed-arena lib in Rust
func (t *Tree) Allot() *Node {
	if len(t.chunks[len(t.chunks)-1]) == lots {
		t.chunks = append(t.chunks, make([]Node, 0, lots))
	}
	chunk := append(t.chunks[len(t.chunks)-1], Node{})
	t.chunks[len(t.chunks)-1] = chunk
	return &chunk[len(chunk)-1]
}

func (t *Tree) Init() *Tree {
	t.chunks = [][]Node{make([]Node, 0, lots)}
	return t
}

func createTree(t *Tree, depth int) *Node {
	node := t.Allot()
	if depth > 0 {
		node.Left, node.Right = createTree(t, depth-1), createTree(t, depth-1)
	}
	return node
}

func checkTree(node *Node) int {
	if node.Left == nil {
		return 1
	}
	return 1 + checkTree(node.Left) + checkTree(node.Right)
}

func main() {
	n := 0
	flag.Parse()
	if flag.NArg() > 0 {
		n, _ = strconv.Atoi(flag.Arg(0))
	}
	run(n)
}

func run(maxDepth int) {
	const minDepth = 4
	var longLivedTree *Node
	var group sync.WaitGroup
	var messages sync.Map

	if minDepth+2 > maxDepth {
		maxDepth = minDepth + 2
	}

	group.Add(1)
	go func() {
		
		messages.Store(-1, fmt.Sprintf("stretch tree of depth %d\t check: %d",
			maxDepth+1, checkTree(createTree(new(Tree).Init(), maxDepth+1))))
		group.Done()
	}()
	group.Add(1)
	go func() {
		longLivedTree = createTree(new(Tree).Init(), maxDepth)
		group.Done()
	}()

	for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
		iters := 1 << (maxDepth - (halfDepth * 2) + minDepth)
		group.Add(1)
		go func(depth, i, chk int) {
			for i := 0; i < iters; i++ {
				chk += checkTree(createTree(new(Tree).Init(), depth))
			}
			messages.Store(depth, fmt.Sprintf("%d\t trees of depth %d\t check: %d",
				i, depth, chk))
			group.Done()
		}(halfDepth*2, iters, 0)
	}

	group.Wait()

	var idxs []int
	messages.Range(func(key, val interface{}) bool {
		idxs = append(idxs, key.(int))
		return true
	})
	sort.Ints(idxs)
	for _, idx := range idxs {
		msg, _ := messages.Load(idx)
		fmt.Println(msg)
	}

	fmt.Printf("long lived tree of depth %d\t check: %d\n",
		maxDepth, checkTree(longLivedTree))
}
