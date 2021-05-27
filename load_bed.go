package main

import (
	"github.com/overerd/clincnv_overlapper/models/bed"
	"github.com/overerd/clincnv_overlapper/regions"
)

func loadBEDFile(options regions.Options) (result *bed.File, err error) {
	result = bed.LoadBEDFile(bed.FileParserOptions{
		Path: options.BedPath,

		GeneColumnIndex: options.BedGeneIndex,
	})

	err = result.Load()

	if err != nil {
		return
	}

	return
}
