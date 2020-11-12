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

func createTree(depth int, nodes *[]Node) *Node {
	var r func(depth int) *Node
	r = func(d int) *Node {
		if d > 0 {
			*nodes = append(*nodes, Node{Left: r(d - 1), Right: r(d - 1)})
		} else {
			*nodes = append(*nodes, Node{})
		}
		return &(*nodes)[len(*nodes)-1]
	}
	return r(depth)
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

	group.Add(2)
	go func() {
		nodes := make([]Node, 0)
		messages.Store(-1, fmt.Sprintf("stretch tree of depth %d\t check: %d",
			maxDepth+1, checkTree(createTree(maxDepth+1, &nodes))))
		group.Done()
	}()
	go func() {
		nodes := make([]Node, 0)
		longLivedTree = createTree(maxDepth, &nodes)
		group.Done()
	}()

	for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
		iters := 1 << (maxDepth - (halfDepth * 2) + minDepth)
		group.Add(1)
		go func(depth, i, chk int) {
			nodes := make([]Node, 0)
			for i := 0; i < iters; i++ {
				chk += checkTree(createTree(depth, &nodes))
				nodes = nodes[:0]
			}
			messages.Store(depth, fmt.Sprintf("%d\t trees of depth %d\t check: %d",
				i, depth, chk))
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
		maxDepth, checkTree(longLivedTree))
}
