package regions

func FilterOverlapsWithNormals(normals *map[string][]RegionData, tumors *map[string][]RegionData) {
	normalsStartPositions := make(map[string]map[uint64]*RegionData, 0)
	normalsEndPositions := make(map[string]map[uint64]*RegionData, 0)

	for chr, regions := range *normals {
		normalsStartPositions[chr] = make(map[uint64]*RegionData, 0)
		normalsEndPositions[chr] = make(map[uint64]*RegionData, 0)

		for _, region := range regions {
			normalsStartPositions[chr][region.Start] = &region
			normalsEndPositions[chr][region.End] = &region
		}
	}

	for chr, regions := range *tumors {
		i := 0
		size := len(regions)

		for i < size {
			chrNormalPositions := normalsStartPositions[chr]

			if _, ok := chrNormalPositions[regions[i].Start]; ok {
				if i == 0 {
					if len(regions) > 1 {
						regions = regions[1:]
					} else {
						regions = []RegionData{}
					}
				} else {
					regions = append(regions[:i], regions[i+1:]...)
				}

				size--
			} else {
				i++
			}
		}

		(*tumors)[chr] = regions
	}
}
