package regions

import (
	"github.com/overerd/clincnv_overlapper/models"
)

type RegionData struct {
	models.ChromosomeRegion
	Power uint16
}

type AnnotatedRegionData struct {
	*RegionData

	Samples []string

	SamplesHash uint32

	Genes string
}

func (s *Overlapper) calculateRegions(positions *map[string][]uint64) (regions map[string][]RegionData) {
	regions = make(map[string][]RegionData, len(*positions))

	for chr, items := range *positions {
		size := len(items) - 1

		if size == 0 {
			size = 1
		}

		regions[chr] = make([]RegionData, size)

		j := 0

		for i := 0; i < size; i++ {
			regions[chr][j] = RegionData{
				ChromosomeRegion: models.ChromosomeRegion{
					Start: items[i],
					End:   items[i+1],
				},

				Power: 0,
			}

			j++
		}
	}

	return
}

func (s *Overlapper) calculateRegionPowers(regions *map[string][]RegionData) {
	for _, file := range s.files {
		correctionPowerSize := 0
		correctedQValue := s.options.MaxQValue

		for _, items := range file.Chromosomes {
			correctionPowerSize += len(items)
		}

		if correctionPowerSize == 0 {
			continue
		}

		if s.options.UseBonferroniCorrection {
			correctedQValue /= float32(correctionPowerSize)
		}

		for chr, items := range file.Chromosomes {
			inputRegions := (*regions)[chr]

			for _, fileRegion := range items {
				for j := range inputRegions {
					region := &inputRegions[j]

					if fileRegion.Start <= region.Start && fileRegion.End >= region.End && fileRegion.QValue <= correctedQValue && fileRegion.LogLikelihood >= 1 {
						region.Power++
					}
				}
			}
		}
	}
}
