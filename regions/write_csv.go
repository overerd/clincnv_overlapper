package regions

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type WriterOptions struct {
	Path           string
	Separator      string
	FieldSeparator string
}

type Writer struct {
	options WriterOptions
}

func (o *WriterOptions) Validate() (err error) {
	return
}

func (w *Writer) WriteRegions(regions *map[string][]AnnotatedRegionData) (err error) {
	fmt.Println(fmt.Sprintf("\nWriting output to file '%s'...", w.options.Path))
	f, err := os.Create(w.options.Path)

	defer func(f *os.File) {
		e := f.Close()

		if e != nil {
			err = e
		}
	}(f)

	if err != nil {
		return
	}

	var bytes int
	var writtenBytes = 0

	writer := bufio.NewWriter(f)

	bytes, err = writer.WriteString("#chr	start	end	genes	samples\n")

	if err != nil {
		return
	}

	writtenBytes += bytes

	i := 0

	chromosomes := make([]string, len(*regions))

	for chr := range *regions {
		chromosomes[i] = chr

		i++
	}

	digits := strings.Split("1234567890", "")
	digitsMap := make(map[string]bool, 10)

	for _, digit := range digits {
		digitsMap[digit] = true
	}

	naturalComparer := func(a, b string) bool {
		aValue, err := strconv.ParseInt(a, 10, 64)

		if err != nil {
			return a < b
		}

		bValue, err := strconv.ParseInt(b, 10, 64)

		if err != nil {
			return a < b
		}

		return aValue < bValue
	}

	sort.Slice(chromosomes, func(i, j int) bool {
		word := chromosomes[i]

		if word[:3] == "chr" {
			return naturalComparer(chromosomes[i][3:], chromosomes[j][3:])
		}

		if word[:1] == "c" {
			if _, ok := digitsMap[word[1:2]]; ok {
				return naturalComparer(chromosomes[i][1:], chromosomes[j][1:])
			}
		}

		return chromosomes[i] < chromosomes[j]
	})

	for _, chr := range chromosomes {
		regions := (*regions)[chr]

		sort.Slice(regions, func(i, j int) bool {
			if regions[i].Start < regions[j].Start {
				return true
			}

			if regions[i].End > regions[j].End && regions[i].Start == regions[j].Start {
				return true
			}

			return false
		})

		for i, item := range regions {
			var samples string

			if len(item.Samples) == 0 {
				samples = "-"
			} else {
				samples = strings.Join(item.Samples, w.options.FieldSeparator)
			}

			bytes, err = writer.WriteString(fmt.Sprintf(
				"%s	%d	%d	%s	%s\n",
				chr, item.Start, item.End, item.Genes, samples,
			))

			if err != nil {
				fmt.Errorf("    [!] region %d had an error", i)

				return
			}

			writtenBytes += bytes
		}
	}

	fmt.Println(fmt.Sprintf("\n %d bytes written\n", writtenBytes))

	err = writer.Flush()

	return
}

func CreateCSVWriter(options WriterOptions) (writer *Writer) {
	writer = &Writer{
		options: options,
	}

	return
}
