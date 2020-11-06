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

type Tree struct {
	Left, Right *Tree
}

func createTree(depth int) *Tree {
	if depth > 0 {
		return &Tree{Left: createTree(depth - 1), Right: createTree(depth - 1)}
	}
	return nil
}

func checkTree(tree *Tree) int {
	if tree == nil {
		return 1
	}
	return 1 + checkTree(tree.Left) + checkTree(tree.Right)
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
	var longLivedTree *Tree
	var group sync.WaitGroup
	var messages sync.Map

	if minDepth+2 > maxDepth {
		maxDepth = minDepth + 2
	}

	group.Add(2)
	go func() {
		tree := createTree(maxDepth + 1)
		fmt.Printf("stretch tree of depth %d\t check: %d\n",
			maxDepth+1, checkTree(tree))
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
