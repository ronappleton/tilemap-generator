package maputils

import "sort"

// Adjacency lists the neighboring tile hashes in the four cardinal directions.
type Adjacency struct {
	Top    []string `json:"top"`
	Bottom []string `json:"bottom"`
	Left   []string `json:"left"`
	Right  []string `json:"right"`
}

// BuildAdjacency returns, for each tile ID, the hashes of neighbouring tiles in
// each direction based on the provided mapping grid.
func BuildAdjacency(tiles []Tile, mapping [][]int) map[int]Adjacency {
	hashByID := make(map[int]string)
	for _, t := range tiles {
		hashByID[t.ID] = t.Hash
	}
	type builder struct {
		top, bottom, left, right map[string]struct{}
	}
	builders := make(map[int]*builder)
	for _, t := range tiles {
		builders[t.ID] = &builder{
			top:    map[string]struct{}{},
			bottom: map[string]struct{}{},
			left:   map[string]struct{}{},
			right:  map[string]struct{}{},
		}
	}
	rows := len(mapping)
	if rows == 0 {
		return nil
	}
	cols := len(mapping[0])
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			id := mapping[y][x]
			b := builders[id]
			if y > 0 {
				b.top[hashByID[mapping[y-1][x]]] = struct{}{}
			}
			if y < rows-1 {
				b.bottom[hashByID[mapping[y+1][x]]] = struct{}{}
			}
			if x > 0 {
				b.left[hashByID[mapping[y][x-1]]] = struct{}{}
			}
			if x < cols-1 {
				b.right[hashByID[mapping[y][x+1]]] = struct{}{}
			}
		}
	}
	res := make(map[int]Adjacency)
	for id, b := range builders {
		res[id] = Adjacency{
			Top:    sortedKeys(b.top),
			Bottom: sortedKeys(b.bottom),
			Left:   sortedKeys(b.left),
			Right:  sortedKeys(b.right),
		}
	}
	return res
}

func sortedKeys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
