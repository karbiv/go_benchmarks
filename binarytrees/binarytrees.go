// The Computer Language Benchmarks Game
// http://benchmarksgame.alioth.debian.org/
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
	arena []Node
	root  *Node
}

func (t *Tree) create(depth int) *Tree {
	var recur func(depth int) *Node
	t.arena = make([]Node, 0, (2 << depth))

	recur = func(d int) *Node {
		if d > 0 {
			t.arena = append(t.arena,
				Node{Left: recur(d - 1), Right: recur(d - 1)})
		} else {
			t.arena = append(t.arena, Node{})
		}
		return &t.arena[len(t.arena)-1]
	}
	recur(depth)
	t.root = &t.arena[len(t.arena)-1] // last elem is a tree root, after recur
	return t
}

func (t *Tree) check() int {
	var recur func(node *Node) int

	recur = func(node *Node) int {
		if node.Left == nil {
			return 1
		}
		return 1 + recur(node.Left) + recur(node.Right)
	}
	return recur(t.root)
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
	var longLivedTree Tree
	var group sync.WaitGroup
	var messages sync.Map

	if minDepth+2 > maxDepth {
		maxDepth = minDepth + 2
	}

	group.Add(2)
	go func() {
		var tree Tree
		messages.Store(0, fmt.Sprintf("stretch tree of depth %d\t check: %d",
			maxDepth+1, tree.create(maxDepth+1).check()))
		group.Done()
	}()
	go func() {
		longLivedTree.create(maxDepth)
		group.Done()
	}()

	for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
		iters := 1 << (maxDepth - (halfDepth * 2) + minDepth)
		group.Add(1)
		go func(depth, i, chk int) {
			var tree Tree
			for i := 0; i < iters; i++ {
				chk += tree.create(depth).check()
			}
			messages.Store(depth, fmt.Sprintf("%d\t trees of depth %d\t check: %d",
				i+1, depth, chk))
			group.Done()
		}(halfDepth*2, iters, 0)
	}

	group.Wait() // wait for all

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
		maxDepth, longLivedTree.check())
}
