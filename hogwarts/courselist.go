//go:build !solution

package hogwarts

import "errors"

type Colour int

var path []string

const (
	White Colour = 0
	Grey  Colour = 1
	Black Colour = 2
)

func dfs(node string, adjacency map[string][]string, colour map[string]Colour, previous map[string]string) bool {
	colour[node] = Grey

	for _, neighbour := range adjacency[node] {
		if colour[neighbour] == White {
			previous[neighbour] = node

			if dfs(neighbour, adjacency, colour, previous) {
				return true
			}
		} else if colour[neighbour] == Grey {
			return true
		}
	}

	colour[node] = Black

	return false
}

func topSort(node string, adjacency map[string][]string, isVisited map[string]bool) {
	if !isVisited[node] {
		isVisited[node] = true

		for _, neighbour := range adjacency[node] {
			topSort(neighbour, adjacency, isVisited)
		}

		path = append(path, node)
	} else {
		return
	}
}

func reversePath() []string {
	var result []string

	for i := len(path) - 1; i >= 0; i-- {
		result = append(result, path[i])
	}

	return result
}

func GetCourseList(prereqs map[string][]string) []string {
	var adjacency = make(map[string][]string)
	var previous = make(map[string]string)
	var isVisited = make(map[string]bool)
	var colour = make(map[string]Colour)

	for k, v := range prereqs {
		for _, prereq := range v {
			adjacency[prereq] = append(adjacency[prereq], k)
		}
	}

	for _, v := range adjacency {
		for _, prereq := range v {
			if dfs(prereq, adjacency, colour, previous) {
				panic(errors.New("cycle dependency in prerequisites"))
			}
		}
	}

	for node := range adjacency {
		if !isVisited[node] {
			topSort(node, adjacency, isVisited)
		}
	}

	return reversePath()
}
