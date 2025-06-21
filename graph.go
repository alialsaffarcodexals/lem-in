package main

// BFS finds the shortest path between start and end using breadth-first search.
// It returns the path as a slice of room names including start and end.
func BFS(g *Graph, start, end string) ([]string, bool) {
	type node struct {
		name string
		prev *node
	}
	visited := make(map[string]bool)
	q := []node{{name: start}}
	visited[start] = true
	var endNode *node
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.name == end {
			endNode = &cur
			break
		}
		for _, nb := range g.Links[cur.name] {
			if !visited[nb] {
				visited[nb] = true
				q = append(q, node{name: nb, prev: &cur})
			}
		}
	}
	if endNode == nil {
		return nil, false
	}
	var path []string
	for n := endNode; n != nil; n = n.prev {
		path = append([]string{n.name}, path...)
	}
	return path, true
}
