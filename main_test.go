package main

import (
	"github.com/overerd/clincnv_overlapper/regions"
	"testing"
)

func TestApp(t *testing.T) {
	// for now only checks if it runs successfully
	options := regions.Options{
		BedPath:             "test/test.bed",
		NormalsFileListPath: "test/normals.txt",
		TumorsFileListPath:  "test/tumors.txt",

		MinOverlap: 2,
		MaxQValue:  1.0,

		BedGeneIndex: 4,

		BufferSize: 100 * 1024,
	}

	writerOptions := regions.WriterOptions{
		Path:           "/dev/null",
		Separator:      "\t",
		FieldSeparator: "; ",
	}

	run(options, writerOptions)

	println()
}
