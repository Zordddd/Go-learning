package main

import (
	"fmt"
)

func DeepFirstSerch(graph map[string][]string, currentDot string) {
    // if graph[currentDot] == nil {
    //     fmt.Println("Wrong start")
    //     os.Exit(1)
    // }

    visited[currentDot] = true
    if len(visited) != 1 {
        fmt.Printf(" -> %s", currentDot)
    } else {
        fmt.Printf("%s", currentDot)
    }

    for _, v := range graph[currentDot] {
        if ok := visited[v]; !ok {
            DeepFirstSerch(graph, v)
        } else {
            continue
        }
    }
}

var visited = map[string]bool{}
var Prereqs = map[string][]string{
    "algorithms":            {"data structures"},
    "calculus":              {"linear algebra"},
    "compilers":             {"data structures", "formal languages", "computer organization"},
    "data structures":       {"discrete math"},
    "databases":             {"data structures"},
    "discrete math":         {"intro to programming"},
    "formal languages":      {"discrete math"},
    "networks":              {"operating systems"},
    "operating systems":     {"data structures", "computer organization"},
    "programming languages": {"data structures", "computer organization"},
}