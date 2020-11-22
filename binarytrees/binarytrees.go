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
	Next *Next
}

type Next struct {
	left, right Node
}

func createTree(depth int) Node {
	if depth > 1 {
		return Node{&Next{createTree(depth - 1), createTree(depth - 1)}}
	}
	return Node{&Next{Node{}, Node{}}}
}

func checkTree(p Node) int {
	sum := 1
	current := p.Next
	for current != nil {
		sum += checkTree(current.right) + 1
		current = current.left.Next
	}
	return sum
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
	var longLivedTree Node
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
		go func(depth, iters, chk int) {
			for i := 0; i < iters; i++ {
				chk += checkTree(createTree(depth))
			}
			messages.Store(depth, fmt.Sprintf("%d\t trees of depth %d\t check: %d",
				iters, depth, chk))
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
