package regions

import (
	"fmt"
	"regexp"
)

type Options struct {
	BedPath             string
	NormalsFileListPath string
	TumorsFileListPath  string

	Debug bool

	MinOverlap uint16
	MaxQValue  float64

	MinLogLikelihood       float64
	MinMedianLogLikelihood float64

	UseBonferroniCorrection bool

	BedGeneIndex uint

	BufferSize uint

	GeneRegexFilter *regexp.Regexp

	SingleRunMode bool
}

func (o *Options) Validate() (err error) {
	if o.MinOverlap > 1 && o.SingleRunMode {
		o.MinOverlap = 1

		fmt.Println(fmt.Sprintf("   [!] min-overlap has been set to 1 due to run-single mode"))
	}

	return
}
