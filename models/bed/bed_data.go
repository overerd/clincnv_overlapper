package bed

import (
	"github.com/overerd/clincnv_overlapper/models"
	"strings"
)

const reLine = "^([^\t]+)\t([^\t]+)\t([^\t]+)\t([^\t]+)(?:\t([^\t]+))*"

type Item struct {
	models.ChromosomeRegionData
	ItemData
}

type ItemData struct {
	Name  string   `json:"name"`
	Genes []string `json:"-"`
}

func (p *ItemData) Fill(groups [][]byte, options *FileParserOptions) (err error) {
	p.Name = string(groups[posName])
	p.Genes = strings.Split(string(groups[options.GeneColumnIndex]), ",")

	return
}
