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

func createTree(depth int) *Node {
	var recur func(depth int) *Node
	nodes := make([]Node, 0, (2 << depth))
	recur = func(d int) *Node {
		if d > 0 {
			nodes = append(nodes, Node{Left: recur(d - 1), Right: recur(d - 1)})
		} else {
			nodes = append(nodes, Node{})
		}
		return &nodes[len(nodes)-1]
	}
	recur(depth)
	return &nodes[len(nodes)-1]
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
		messages.Store(-1, fmt.Sprintf("stretch tree of depth %d\t check: %d",
			maxDepth+1, checkTree(createTree(maxDepth+1))))
		group.Done()
	}()
	go func() {
		longLivedTree = createTree(maxDepth)
		group.Done()
	}()

	for halfDepth := minDepth / 2; halfDepth < maxDepth/2+1; halfDepth++ {
		iters := 1 << (maxDepth - (halfDepth * 2) + minDepth)
		group.Add(1)
		go func(depth, i, chk int) {
			for i := 0; i < iters; i++ {
				chk += checkTree(createTree(depth))
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
