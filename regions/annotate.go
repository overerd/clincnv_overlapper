package regions

import (
	"fmt"
	"github.com/overerd/clincnv_overlapper/models/bed"
	"github.com/overerd/clincnv_overlapper/models/clincnv"
	"hash/crc32"
	"regexp"
	"sort"
	"strings"
)

func filterGenes(genes *[]string, filter *regexp.Regexp) {
	var indexesToFilter []int

	for i, gene := range *genes {
		if filter.MatchString(gene) {
			indexesToFilter = append(indexesToFilter, i)
		}
	}

	size := len(*genes)

	if size == len(indexesToFilter) {
		*genes = []string{"-"}

		return
	}

	if len(indexesToFilter) == 0 {
		return
	}

	sort.Slice(indexesToFilter, func(i, j int) bool {
		return indexesToFilter[i] > indexesToFilter[j]
	})

	for _, index := range indexesToFilter {
		if index == size-1 {
			*genes = (*genes)[:index]
		} else {
			if index == 0 {
				*genes = (*genes)[1:]
			} else {
				*genes = append((*genes)[:index], (*genes)[index+1:]...)
			}
		}
	}
}

func renameGenes(genes *[]string) (err error) {
	r, err := regexp.Compile("^ref\\|")

	if err != nil {
		return err
	}

	for i, gene := range *genes {
		(*genes)[i] = r.ReplaceAllString(gene, "")
	}

	return
}

func AnnotateOverlaps(
	bed *bed.File,
	files *[]*clincnv.TableFile,
	regions *map[string][]RegionData,
	options Options,
	writerOptions WriterOptions,
) (r *map[string][]AnnotatedRegionData) {
	results := make(map[string][]AnnotatedRegionData)

	correctionPowerSize := 0

	crcTable := crc32.MakeTable(23)

	for _, regions := range *regions {
		correctionPowerSize += len(regions)
	}

	for chr, regions := range *regions {
		results[chr] = make([]AnnotatedRegionData, len(regions))

		correctedQValue := options.MaxQValue

		if options.UseBonferroniCorrection {
			correctedQValue /= float32(correctionPowerSize)
		}

		for i := range regions {
			item := &results[chr][i]

			item.RegionData = &regions[i]

			for _, file := range *files {
				if found, fileItem := file.FindClosestItem(chr, item.Start, item.End); found {
					sampleString := fmt.Sprintf(
						"%s (%d)",
						file.Options.Name,
						(*fileItem).CNChange,
					)

					if fileItem.QValue <= correctedQValue && fileItem.LogLikelihood >= 1 {
						item.Samples = append(item.Samples, sampleString)
					}
				}
			}

			sort.Strings(item.Samples)

			item.SamplesHash = crc32.Checksum([]byte(strings.Join(item.Samples, ";")), crcTable)
		}
	}

	var keysToDelete []string

	minOverlaps := int(options.MinOverlap)

	for chr, regions := range results {
		if len(regions) == 0 {
			keysToDelete = append(keysToDelete, chr)

			continue
		}

		indexesToDelete := make([]int, 0)

		pRegion := &AnnotatedRegionData{}

		for i, item := range regions {
			if len(item.Samples) <= minOverlaps && !options.SingleRunMode {
				indexesToDelete = append(indexesToDelete, i)
			} else {
				if pRegion.SamplesHash == item.SamplesHash && pRegion.End >= item.Start {
					indexesToDelete = append(indexesToDelete, i)
				} else {
					pRegion = &item
				}
			}
		}

		if len(indexesToDelete) == 0 {
			continue
		}

		sort.Slice(indexesToDelete, func(i, j int) bool {
			return indexesToDelete[i] > indexesToDelete[j]
		})

		size := len(regions)

		for _, index := range indexesToDelete {
			if index == size-1 {
				if len(regions) > 1 {
					regions = regions[:index]
				} else {
					regions = []AnnotatedRegionData{}
				}
			} else {
				regions = append(regions[:index], regions[index+1:]...)
			}
		}

		if bed != nil {
			for i := range regions {
				item := &results[chr][i]

				if genes := bed.SelectAllGenes(chr, item.Start, item.End); len(genes) > 0 {
					if options.GeneRegexFilter != nil {
						filterGenes(&genes, options.GeneRegexFilter)
					}

					err := renameGenes(&genes)

					if err != nil {
						fmt.Errorf("%s", err.Error())
					}

					sort.Strings(genes)

					item.Genes = strings.Join(genes, writerOptions.FieldSeparator)
				}
			}
		}

		results[chr] = regions

		if len(regions) == 0 {
			keysToDelete = append(keysToDelete, chr)
		}
	}

	for _, index := range keysToDelete {
		delete(results, index)
	}

	r = &results

	return
}
