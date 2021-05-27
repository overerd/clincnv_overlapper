package regions

func (s *Overlapper) filterPoweredRegions(regions *map[string][]RegionData) {
	for chr, items := range *regions {
		i := 0
		size := len(items)

		for i < size {
			if items[i].Power < s.options.MinOverlap {
				if i == 0 {
					if len(items) > 1 {
						items = items[1:]
					} else {
						items = []RegionData{}
					}
				} else {
					items = append(items[:i], items[i+1:]...)
				}

				size--
			} else {
				i++
			}
		}

		(*regions)[chr] = items
	}
}
