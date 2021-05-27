package regions

import "sort"

func (s *Overlapper) calculatePositions() (positions map[string][]uint64, err error) {
	m := make(map[string]map[uint64]bool, 0)

	for _, file := range s.files {
		for chr, items := range file.Chromosomes {
			if _, ok := m[chr]; !ok {
				m[chr] = make(map[uint64]bool, 0)
			}

			for _, item := range items {
				m[chr][item.Start] = true
				m[chr][item.End] = true
			}
		}
	}

	positions = make(map[string][]uint64, len(m))

	for chr, innerMap := range m {
		i := 0

		if _, ok := positions[chr]; !ok {
			positions[chr] = make([]uint64, len(innerMap))
		}

		for k := range innerMap {
			positions[chr][i] = k
			i++
		}
	}

	for _, innerMap := range positions {
		sort.Slice(innerMap, func(i, j int) bool {
			return innerMap[i] < innerMap[j]
		})
	}

	return
}
