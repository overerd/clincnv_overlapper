package main

import (
	"fmt"
	"github.com/overerd/clincnv_overlapper/models/bed"
	"github.com/overerd/clincnv_overlapper/models/clincnv"
	"github.com/overerd/clincnv_overlapper/regions"
)

func readOverlaps(target string, files map[string]string, options regions.Options) (results *map[string][]regions.RegionData, parsedFiles *[]*clincnv.TableFile, err error) {
	println(fmt.Sprintf("\nReading %s ClinCNV samples with buffer size %d...", target, options.BufferSize))

	s, err := regions.CreateSummarizer(options)

	if err != nil {
		return
	}

	err = s.LoadFiles(files)

	if err != nil {
		return
	}

	parsedFiles = s.GetFiles()

	overlaps, err := s.CalculateOverlaps()

	if err != nil {
		return
	}

	results = overlaps

	return
}

func calculateOverlaps(
	bed *bed.File,
	normalsFileListPath, tumorsFileListPath string,
	options regions.Options,
	writerOptions regions.WriterOptions,
	calculateSeparateOverlaps bool,
) (
	results *map[string][]regions.AnnotatedRegionData,
	separateResults *map[string]*map[string][]regions.AnnotatedRegionData,
	err error,
) {
	var normals map[string]string
	var tumors map[string]string

	var normalRegions *map[string][]regions.RegionData
	var tumorRegions *map[string][]regions.RegionData
	var intermediateResults *map[string][]regions.RegionData

	var files *[]*clincnv.TableFile
	var tumorFiles *[]*clincnv.TableFile
	var normalFiles *[]*clincnv.TableFile

	normals, err = loadFileList(normalsFileListPath)

	if err != nil {
		return
	}

	minOverlapValue := options.MinOverlap

	if tumorsFileListPath != "" {
		tumors, err = loadFileList(tumorsFileListPath)

		if err != nil {
			return
		}

		tumorRegions, tumorFiles, err = readOverlaps("tumor", tumors, regions.Options{
			Debug:         options.Debug,
			MinOverlap:    options.MinOverlap,
			MaxQValue:     options.MaxQValue,
			SingleRunMode: options.SingleRunMode,
			BufferSize:    options.BufferSize,
			BedGeneIndex:  options.BedGeneIndex,

			MinLogLikelihood: options.MinLogLikelihood,

			MinMedianLogLikelihood: options.MinMedianLogLikelihood,

			UseBonferroniCorrection: options.UseBonferroniCorrection,
		})

		if err != nil {
			return
		}

		minOverlapValue = uint16(1)
	}

	normalRegions, normalFiles, err = readOverlaps("normal", normals, regions.Options{
		Debug:         options.Debug,
		MinOverlap:    minOverlapValue,
		MaxQValue:     options.MaxQValue,
		SingleRunMode: options.SingleRunMode,
		BufferSize:    options.BufferSize,

		MinLogLikelihood: options.MinLogLikelihood,

		MinMedianLogLikelihood: options.MinMedianLogLikelihood,

		UseBonferroniCorrection: options.UseBonferroniCorrection,
	})

	if err != nil {
		return
	}

	if tumorsFileListPath != "" {
		fmt.Println("\n\nFiltering tumor regions with normal...")

		regions.FilterOverlapsWithNormals(normalRegions, tumorRegions)

		intermediateResults = tumorRegions
		files = tumorFiles
	} else {
		intermediateResults = normalRegions
		files = normalFiles
	}

	annotatedRegions := regions.AnnotateOverlaps(bed, files, intermediateResults, options, writerOptions)

	results = annotatedRegions

	if calculateSeparateOverlaps {
		fmt.Print("Separating overlaps per sample file...\n\n")

		separateResults, err = regions.SeparateOverlaps(annotatedRegions, files, options)
	}

	return
}
