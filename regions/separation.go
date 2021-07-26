package regions

import (
	"fmt"
	"github.com/overerd/clincnv_overlapper/models"
	"github.com/overerd/clincnv_overlapper/models/clincnv"
)

func SeparateOverlaps(
	regions *map[string][]AnnotatedRegionData,
	files *[]*clincnv.TableFile,
	options Options,
) (r *map[string]*map[string][]AnnotatedRegionData, err error) {
	results := make(map[string]*map[string][]AnnotatedRegionData)

	var prevState byte
	var prevItem *clincnv.Item
	var items []*clincnv.Item
	var data AnnotatedRegionData
	var interval models.ChromosomeRegion

	for _, file := range *files {
		fmt.Println(fmt.Sprintf(" [*] %s", file.Options.Name))

		sampleRegions := make(map[string][]AnnotatedRegionData)

		results[file.Options.Name] = &sampleRegions

		correctionPowerSize := 0

		for _, regions := range *regions {
			correctionPowerSize += len(regions)
		}

		for chr, chrRegions := range *regions {
			prevItem = &clincnv.Item{}
			prevState = 255

			correctedQValue := options.MaxQValue

			if options.UseBonferroniCorrection {
				correctedQValue /= float32(correctionPowerSize)
			}

			for _, region := range chrRegions {
				items = file.FindIntervalItems(chr, region.Start, region.End)

				for _, item := range items {
					if item.Start >= region.End || item.End <= region.Start || item.QValue > correctedQValue || item.LogLikelihood < 1 {
						continue
					}

					interval = models.ChromosomeRegion{
						Start: item.Start,
						End:   item.End,
					}

					if interval.Start < region.Start {
						interval.Start = region.Start
					}

					if interval.End > region.End {
						interval.End = region.End
					}

					if prevItem.End == interval.End && prevState == item.CNChange {
						continue
					}

					prevItem = item
					prevState = item.CNChange

					data = AnnotatedRegionData{
						RegionData: &RegionData{
							ChromosomeRegion: interval,
						},
						Samples: []string{fmt.Sprintf("%d", item.CNChange)},
						Genes:   "",
					}

					sampleRegions[chr] = append(sampleRegions[chr], data)
				}
			}
		}
	}

	r = &results

	return
}
