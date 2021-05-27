package models

import "strconv"

const (
	_ = iota
	PosChr
	PosStart
	PosEnd
)

type ChromosomeRegion struct {
	Start uint64 `json:"start"`
	End   uint64 `json:"end"`
}

type ChromosomeRegionData struct {
	Chr string `json:"chr"`
	ChromosomeRegion
}

func (p *ChromosomeRegionData) Fill(groups [][]byte) (err error) {
	p.Chr = string(groups[PosChr])

	p.Start, err = strconv.ParseUint(string(groups[PosStart]), 10, 64)

	if err != nil {
		return
	}

	p.End, err = strconv.ParseUint(string(groups[PosEnd]), 10, 64)

	if err != nil {
		return
	}

	return
}
