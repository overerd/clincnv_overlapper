package regions

import (
	"github.com/overerd/clincnv_overlapper/models/clincnv"
)

type Overlapper struct {
	options Options

	files []*clincnv.TableFile

	regions []RegionData
}

func (s *Overlapper) CalculateOverlaps() (results *map[string][]RegionData, err error) {
	println("\nCalculating overlaps...")

	positions, err := s.calculatePositions()

	if err != nil {
		return
	}

	regions := s.calculateRegions(&positions)

	s.calculateRegionPowers(&regions)

	s.filterPoweredRegions(&regions)

	results = &regions

	return
}

func (s *Overlapper) LoadFiles(paths map[string]string) (err error) {
	err = s.loadCNVFiles(paths)

	if s.options.Debug {
		return
	}

	println()

	return
}

func CreateSummarizer(options Options) (s *Overlapper, err error) {
	err = options.Validate()

	if err != nil {
		return
	}

	s = &Overlapper{options: options}

	return
}
