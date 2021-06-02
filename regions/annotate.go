package regions

import (
	"fmt"
	"github.com/overerd/clincnv_overlapper/models/bed"
	"github.com/overerd/clincnv_overlapper/models/clincnv"
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

	for chr, regions := range *regions {
		results[chr] = make([]AnnotatedRegionData, len(regions))

		for i := range regions {
			item := &results[chr][i]

			item.RegionData = &regions[i]

			if bed != nil {
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

			for _, file := range *files {
				if found, fileItem := file.FindClosestItem(chr, item.Start, item.End); found {
					sampleString := fmt.Sprintf(
						"%s (%d)",
						file.Options.Name,
						(*fileItem).CNChange,
					)

					if fileItem.QValue <= options.MaxQValue {
						item.Samples = append(item.Samples, sampleString)
					}
				}
			}

			sort.Strings(item.Samples)
		}
	}

	var keysToDelete []string

	for chr, regions := range results {
		if len(regions) == 0 {
			keysToDelete = append(keysToDelete, chr)

			continue
		}

		indexesToDelete := make([]int, 0)

		for i, item := range regions {
			if len(item.Samples) <= 1 && !options.SingleRunMode {
				indexesToDelete = append(indexesToDelete, i)
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
