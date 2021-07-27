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

	crcTable := crc32.MakeTable(23)

	for chr, chrRegions := range *regions {
		results[chr] = make([]AnnotatedRegionData, len(chrRegions))

		for i := range chrRegions {
			item := &results[chr][i]

			item.RegionData = &chrRegions[i]

			for _, file := range *files {
				correctedQValue := options.MaxQValue
				correctionPowerSize := 0

				for _, fileChrRegions := range *regions {
					correctionPowerSize += len(fileChrRegions)
				}

				if options.UseBonferroniCorrection {
					correctedQValue /= float32(correctionPowerSize)
				}

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

	for chr, chrRegions := range results {
		if len(chrRegions) == 0 {
			keysToDelete = append(keysToDelete, chr)

			continue
		}

		indexesToDelete := make([]int, 0)

		pRegion := &AnnotatedRegionData{}

		for i, item := range chrRegions {
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

		size := len(chrRegions)

		for _, index := range indexesToDelete {
			if index == size-1 {
				if len(chrRegions) > 1 {
					chrRegions = chrRegions[:index]
				} else {
					chrRegions = []AnnotatedRegionData{}
				}
			} else {
				chrRegions = append(chrRegions[:index], chrRegions[index+1:]...)
			}
		}

		if bed != nil {
			for i := range chrRegions {
				item := &results[chr][i]

				if genes := bed.SelectAllGenes(chr, item.Start, item.End); len(genes) > 0 {
					if options.GeneRegexFilter != nil {
						filterGenes(&genes, options.GeneRegexFilter)
					}

					err := renameGenes(&genes)

					if err != nil {
						fmt.Errorf("%s", err.Error())
					}

					if len(genes) > 0 {
						sort.Strings(genes)

						item.Genes = strings.Join(genes, writerOptions.FieldSeparator)
					} else {
						item.Genes = "-"
					}
				}
			}
		}

		results[chr] = chrRegions

		if len(chrRegions) == 0 {
			keysToDelete = append(keysToDelete, chr)
		}
	}

	for _, index := range keysToDelete {
		delete(results, index)
	}

	r = &results

	return
}
