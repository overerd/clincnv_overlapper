package clincnv

import "github.com/overerd/clincnv_overlapper/models"

type State byte

const (
	StateDEL State = iota
	StateAMP
	StateLOH
)

const (
	PosGermlineCNChange = iota + models.PosEnd + 1
	PosGermlineLogLikelihood
	PosGermlineNoOfRegions
	PosGermlineLength
	PosGermlinePotential
	PosGermlineGenes
	PosGermlineQValue
)

const (
	PosSomaticTumorCNChange = iota + models.PosEnd + 1
	PosSomaticState
	PosSomaticMajorCNAllele
	PosSomaticMinorCNAllele
	PosSomaticTumorClonality
	PosSomaticCNChange
	PosSomaticLogLikelihood
	PosSomaticMedianLogLikelihood
	PosSomaticNoOfRegions
	PosSomaticMajorCNAllele2
	PosSomaticMinorCNAllele2
	PosSomaticTumorClonality2
	PosSomaticGenes
	PosSomaticOnTargetRDCILower
	PosSomaticOnTargetRDCIUpper
	PosSomaticOffTargetRDCILower
	PosSomaticOffTargetRDCIUpper
	PosSomaticLowMedTumorBAF
	PosSomaticHighMedTumorBAF
	PosSomaticBAFQValueFDR
	PosSomaticOverallQValue
)

type TableFileParserOptions struct {
	Name string `json:"name"`
	Path string `json:"path"`

	BufferSize uint
}
