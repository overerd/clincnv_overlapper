package bed

import (
	"github.com/overerd/clincnv_overlapper/models"
)

const (
	posName = iota + models.PosEnd + 1
)

const PosGenes = posName

type FileParserOptions struct {
	Path string `json:"path"`

	GeneColumnIndex uint `json:"gene_column_index"`
}
