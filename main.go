package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/overerd/clincnv_overlapper/models/bed"
	"github.com/overerd/clincnv_overlapper/regions"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	version = "0.0.6"
	build   = "0"
)

const AppTitle = "clincnv-overlapper"
const AppDescription = "It provides a way to calculate overlaps in multiple ClinCNV tables."

func invokeError(err error) {
	if err == nil {
		return
	}

	println("\n")
	log.Fatalln(err.Error())
}

func main() {
	appString := fmt.Sprintf("%s v%s (build %s)\n", AppTitle, version, build)

	fmt.Println(fmt.Sprintf("Starting %s...", appString))

	options, writerOptions, dirWriterOptions, err := setup()

	if err != nil {
		invokeError(err)
	}

	invokeError(run(options, writerOptions, dirWriterOptions))

	println()
}

func setup() (options regions.Options, writerOptions, dirWriterOptions regions.WriterOptions, err error) {
	parser := argparse.NewParser(os.Args[0], AppDescription)

	normalsFileListPath := parser.String("n", "normals", &argparse.Options{
		Required: true,
		Help:     "path to a file with '\n' separated file names to each ClinCNV table",
	})

	tumorsFileListPath := parser.String("t", "tumors", &argparse.Options{
		Required: false,
		Help:     "path to a file with '\n' separated file names to each ClinCNV table",
	})

	bedPath := parser.String("b", "bed", &argparse.Options{
		Required: false,
		Help:     "path to gene annotated BED-file",
	})

	bedGeneColumnIndex := parser.Int("", "bed-gene-index", &argparse.Options{
		Required: false,
		Help:     "genes column index for bed file",
		Default:  bed.PosGenes,
	})

	minOverlap := parser.Int("m", "min-overlap", &argparse.Options{
		Required: false,
		Help:     "minimum overlap to filter gathered regions for each chromosome",
		Default:  2,
	})

	bufferSize := parser.Int("r", "buffer-size", &argparse.Options{
		Required: false,
		Help:     "buffer size for clincnv files parser",
		Default:  100 * 1024,
	})

	outputPath := parser.String("o", "output", &argparse.Options{
		Required: true,
		Help:     "output regions file",
	})

	outputDirPath := parser.String("d", "output-directory", &argparse.Options{
		Required: false,
		Help:     "output regions file",
	})

	outputSeparator := parser.String("s", "output-separator", &argparse.Options{
		Required: false,
		Help:     "output file separator",
		Default:  "\t",
	})

	geneFiltersFile := parser.String("x", "genes-exclude-filter", &argparse.Options{
		Required: false,
		Help:     "file with genes substr filter on each line",
		Default:  "",
	})

	outputFieldSeparator := parser.String("f", "output-field-separator", &argparse.Options{
		Required: false,
		Help:     "output file field separator",
		Default:  " | ",
	})

	maxQValue := parser.Float("q", "max-qvalue", &argparse.Options{
		Required: false,
		Help:     "maximum qvalue of CNV region",
		Default:  1.0,
	})

	singleFlag := parser.Flag("", "run-single", &argparse.Options{
		Required: false,
		Help:     "run with single target in normals",
	})

	invokeError(parser.Parse(os.Args))

	if *minOverlap < 1 {
		invokeError(errors.New("-m|--min-overlap should be >= 1"))
	}

	if *outputDirPath != "" {
		info, err := os.Stat(*outputDirPath)

		if os.IsNotExist(err) {
			err = errors.New(fmt.Sprintf("Path '%s' does not exist", *outputDirPath))

			return options, writerOptions, dirWriterOptions, err
		}

		if !info.IsDir() {
			err = errors.New(fmt.Sprintf("Path '%s' is not a directory", *outputDirPath))

			return options, writerOptions, dirWriterOptions, err
		}
	}

	var geneFilter *regexp.Regexp

	if *geneFiltersFile != "" {
		geneFilter, err = readFiltersFile(*geneFiltersFile)

		if err != nil {
			return
		}
	}

	options = regions.Options{
		BedPath:             *bedPath,
		NormalsFileListPath: *normalsFileListPath,
		TumorsFileListPath:  *tumorsFileListPath,

		MinOverlap:    uint16(*minOverlap),
		MaxQValue:     float32(*maxQValue),
		SingleRunMode: *singleFlag,

		BedGeneIndex: uint(*bedGeneColumnIndex),

		BufferSize: uint(*bufferSize),

		GeneRegexFilter: geneFilter,
	}

	writerOptions = regions.WriterOptions{
		Path:           *outputPath,
		Separator:      *outputSeparator,
		FieldSeparator: *outputFieldSeparator,
	}

	dirWriterOptions = regions.WriterOptions{
		Path:           *outputDirPath,
		Separator:      *outputSeparator,
		FieldSeparator: *outputFieldSeparator,
	}

	return
}

func readFiltersFile(path string) (result *regexp.Regexp, err error) {
	file, err := os.Open(path)

	if err != nil {
		return
	}

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			err = e
		}
	}(file)

	scanner := bufio.NewScanner(file)

	var filters []string

	for scanner.Scan() {
		str := scanner.Text()

		if str == "" {
			continue
		}

		filters = append(filters, fmt.Sprintf("(?:%s)", str))
	}

	if err = scanner.Err(); err != nil {
		return
	}

	if len(filters) == 0 {
		return
	}

	result, err = regexp.Compile(fmt.Sprintf("^(?:%s)", strings.Join(filters, "|")))

	return
}

func run(options regions.Options, writerOptions regions.WriterOptions, dirWriterOptions regions.WriterOptions) (err error) {
	var bedFile *bed.File

	err = options.Validate()

	if err != nil {
		return
	}

	err = writerOptions.Validate()

	if err != nil {
		return
	}

	if options.BedPath != "" {
		fmt.Println(fmt.Sprintf("\nReading BED file '%s'...", options.BedPath))

		bedFile, err = loadBEDFile(options)

		if err != nil {
			return
		}

		fmt.Println(fmt.Sprintf(" \n%d lines loaded", len(bedFile.Items)))

		println()
	}

	overlaps, separateOverlaps, err := calculateOverlaps(
		bedFile,
		options.NormalsFileListPath,
		options.TumorsFileListPath,
		options,
		writerOptions,
		dirWriterOptions.Path != "",
	)

	if err != nil {
		return
	}

	fmt.Println(fmt.Sprintf(" \nWriting regions to '%s'...\n", writerOptions.Path))

	writer := regions.CreateCSVWriter(writerOptions)

	err = writer.WriteRegions(overlaps)

	if err != nil {
		return
	}

	if dirWriterOptions.Path != "" {
		fmt.Println(fmt.Sprintf(" \nWriting separate regions to directory '%s'...\n", dirWriterOptions.Path))

		multiWriter := regions.CreateMultiCSVWriter(dirWriterOptions)

		err = multiWriter.WriteSamples(separateOverlaps)

		if err != nil {
			return
		}
	}

	println("\nDone!\n")

	return
}
