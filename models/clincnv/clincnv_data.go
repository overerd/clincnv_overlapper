package clincnv

import (
	"errors"
	"fmt"
	"github.com/overerd/clincnv_overlapper/models"
	"strconv"
	"strings"
)

const reLineGermline = "^([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)?\\t([^\\t]+)$"

const reLineSomatic = "^([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)?\\t([^\\t]+)?\\t([^\\t]+)?\\t([^\\t]+)?\\t([^\\t]+)?\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t([^\\t]+)\\t(.+)$"

type Item struct {
	models.ChromosomeRegionData
	ChromosomeData
}

type ChromosomeData struct {
	CNChange            byte    `json:"cn_change"`
	NormCNChange        string  `json:"norm_cn_change"`
	State               State   `json:"state"`
	MajorCNAllele       uint32  `json:"major_cn_allele"`
	MinorCNAllele       uint32  `json:"minor_cn_allele"`
	TumorClonality      float32 `json:"tumor_clonality"`
	MajorCNAllele2      uint32  `json:"major_cn_allele2"`
	MinorCNAllele2      uint32  `json:"minor_cn_allele2"`
	TumorClonality2     float32 `json:"tumor_clonality2"`
	LogLikelihood       float64 `json:"log_likelihood"`
	MedianLogLikelihood float64 `json:"median_log_likelihood"`
	NoOfRegions         uint32  `json:"no_of_regions"`
	Length              float32 `json:"length"`
	Potential           float32 `json:"potential"`
	Genes               string  `json:"genes"`
	OnTargetRDCILower   float32 `json:"Ontarget_RD_CI_lower"`
	OnTargetRDCIUpper   float32 `json:"Ontarget_RD_CI_upper"`
	OffTargetRDCILower  float32 `json:"Offtarget_RD_CI_lower"`
	OffTargetRDCIUpper  float32 `json:"Offtarget_RD_CI_upper"`
	LowMedTumorBAF      float32 `json:"Lowmed_tumor_BAF"`
	HighMedTumorBAF     float32 `json:"Highmed_tumor_BAF"`
	BAFQValueFDR        float32 `json:"BAF_qval_fdr"`
	QValue              float32 `json:"q_value"`
}

func (p *ChromosomeData) Fill(groups [][]byte, isSomatic bool) (err error) {
	if isSomatic {
		var CNChange uint64
		var logLikelihood float64
		var medianLogLikelihood float64
		var noOfRegions uint64
		var majorCNAllele uint64
		var minorCNAllele uint64
		var tumorClonality float64
		var majorCNAllele2 uint64
		var minorCNAllele2 uint64
		var tumorClonality2 float64
		var overallQValue float64
		var state State
		var onTargetRDCILower float64
		var onTargetRDCIUpper float64
		var offTargetRDCILower float64
		var offTargetRDCIUpper float64
		var lowMedTumorBAF float64
		var highMedTumorBAF float64
		var bafQValueFDR float64

		stateStr := string(groups[PosSomaticState])

		switch stateStr {
		case "DEL":
			state = StateDEL
		case "AMP":
			state = StateAMP
		case "LOH":
			state = StateLOH
		default:
			return errors.New(fmt.Sprintf("unknown somatic region state: %s", stateStr))
		}

		CNChange, err = strconv.ParseUint(string(groups[PosSomaticTumorCNChange]), 10, 8)

		if err != nil {
			return
		}

		normCNChange := string(groups[PosSomaticCNChange])

		logLikelihood, err = strconv.ParseFloat(string(groups[PosSomaticLogLikelihood]), 32)

		if err != nil {
			return
		}

		medianLogLikelihood, err = strconv.ParseFloat(string(groups[PosSomaticMedianLogLikelihood]), 32)

		if err != nil {
			return
		}

		noOfRegions, err = strconv.ParseUint(string(groups[PosSomaticNoOfRegions]), 10, 32)

		if err != nil {
			return
		}

		majorCNAllele, err = strconv.ParseUint(string(groups[PosSomaticMajorCNAllele]), 10, 32)

		if err != nil {
			return
		}

		minorCNAllele, err = strconv.ParseUint(string(groups[PosSomaticMinorCNAllele]), 10, 32)

		if err != nil {
			return
		}

		tumorClonality, err = strconv.ParseFloat(string(groups[PosSomaticTumorClonality]), 32)

		if err != nil {
			return
		}

		majorCNAllele2Str := string(groups[PosSomaticMajorCNAllele2])

		if majorCNAllele2Str != "" {
			majorCNAllele2, err = strconv.ParseUint(majorCNAllele2Str, 10, 32)

			if err != nil {
				return
			}
		}

		minorCNAllele2Str := string(groups[PosSomaticMinorCNAllele2])

		if minorCNAllele2Str != "" {
			minorCNAllele2, err = strconv.ParseUint(minorCNAllele2Str, 10, 32)

			if err != nil {
				return
			}
		}

		tumorClonality2Str := string(groups[PosSomaticTumorClonality2])

		if tumorClonality2Str != "" {
			tumorClonality2, err = strconv.ParseFloat(tumorClonality2Str, 32)

			if err != nil {
				return
			}
		}

		onTargetRDCILowerStr := string(groups[PosSomaticOnTargetRDCILower])

		if strings.ToLower(onTargetRDCILowerStr) != "na" && onTargetRDCILowerStr != "" {
			onTargetRDCILower, err = strconv.ParseFloat(onTargetRDCILowerStr, 32)

			if err != nil {
				return
			}
		}

		onTargetRDCIUpperStr := string(groups[PosSomaticOnTargetRDCIUpper])

		if strings.ToLower(onTargetRDCIUpperStr) != "na" && onTargetRDCIUpperStr != "" {
			onTargetRDCIUpper, err = strconv.ParseFloat(onTargetRDCIUpperStr, 32)

			if err != nil {
				return
			}
		}

		offTargetRDCILowerStr := string(groups[PosSomaticOffTargetRDCILower])

		if strings.ToLower(offTargetRDCILowerStr) != "na" && offTargetRDCILowerStr != "" {
			offTargetRDCILower, err = strconv.ParseFloat(offTargetRDCILowerStr, 32)

			if err != nil {
				return
			}
		}

		offTargetRDCIUpperStr := string(groups[PosSomaticOffTargetRDCIUpper])

		if strings.ToLower(offTargetRDCIUpperStr) != "na" && offTargetRDCIUpperStr != "" {
			offTargetRDCIUpper, err = strconv.ParseFloat(offTargetRDCIUpperStr, 32)

			if err != nil {
				return
			}
		}

		lowMedTumorBAFStr := string(groups[PosSomaticLowMedTumorBAF])

		if strings.ToLower(lowMedTumorBAFStr) != "na" && lowMedTumorBAFStr != "" {
			lowMedTumorBAF, err = strconv.ParseFloat(lowMedTumorBAFStr, 32)

			if err != nil {
				return
			}
		}

		highMedTumorBAFStr := string(groups[PosSomaticHighMedTumorBAF])

		if strings.ToLower(highMedTumorBAFStr) != "na" && highMedTumorBAFStr != "" {
			highMedTumorBAF, err = strconv.ParseFloat(highMedTumorBAFStr, 32)

			if err != nil {
				return
			}
		}

		bafQValueFDRStr := string(groups[PosSomaticBAFQValueFDR])

		if strings.ToLower(bafQValueFDRStr) != "na" && bafQValueFDRStr != "" {
			bafQValueFDR, err = strconv.ParseFloat(bafQValueFDRStr, 32)

			if err != nil {
				return
			}
		}

		overallQValue, err = strconv.ParseFloat(strings.ReplaceAll(string(groups[PosSomaticOverallQValue]), " ", ""), 32)

		if err != nil {
			return
		}

		p.CNChange = uint8(CNChange)
		p.NormCNChange = normCNChange
		p.State = state
		p.MajorCNAllele = uint32(majorCNAllele)
		p.MinorCNAllele = uint32(minorCNAllele)
		p.TumorClonality = float32(tumorClonality)
		p.MajorCNAllele2 = uint32(majorCNAllele2)
		p.MinorCNAllele2 = uint32(minorCNAllele2)
		p.TumorClonality2 = float32(tumorClonality2)
		p.LogLikelihood = logLikelihood
		p.MedianLogLikelihood = medianLogLikelihood
		p.NoOfRegions = uint32(noOfRegions)
		p.Genes = string(groups[PosSomaticGenes])
		p.OnTargetRDCILower = float32(onTargetRDCILower)
		p.OnTargetRDCIUpper = float32(onTargetRDCIUpper)
		p.OffTargetRDCILower = float32(offTargetRDCILower)
		p.OffTargetRDCIUpper = float32(offTargetRDCIUpper)
		p.LowMedTumorBAF = float32(lowMedTumorBAF)
		p.HighMedTumorBAF = float32(highMedTumorBAF)
		p.BAFQValueFDR = float32(bafQValueFDR)
		p.QValue = float32(overallQValue)

		if strings.ToLower(p.Genes) == "na" {
			p.Genes = ""
		}
	} else {
		var cNChange uint64
		var logLikelihood float64
		var noOfRegions uint64
		var length float64
		var potential float64
		var qValue float64

		cNChange, err = strconv.ParseUint(string(groups[PosGermlineCNChange]), 10, 8)

		if err != nil {
			return
		}

		logLikelihood, err = strconv.ParseFloat(string(groups[PosGermlineLogLikelihood]), 32)

		if err != nil {
			return
		}

		noOfRegions, err = strconv.ParseUint(string(groups[PosGermlineNoOfRegions]), 10, 32)

		if err != nil {
			return
		}

		length, err = strconv.ParseFloat(strings.ReplaceAll(string(groups[PosGermlineLength]), " ", ""), 32)

		if err != nil {
			return
		}

		potential, err = strconv.ParseFloat(strings.ReplaceAll(string(groups[PosGermlinePotential]), " ", ""), 32)

		if err != nil {
			return
		}

		qValue, err = strconv.ParseFloat(strings.ReplaceAll(string(groups[PosGermlineQValue]), " ", ""), 32)

		if err != nil {
			return
		}

		p.CNChange = uint8(cNChange)
		p.LogLikelihood = logLikelihood
		p.NoOfRegions = uint32(noOfRegions)
		p.Length = float32(length)
		p.Potential = float32(potential)
		p.Genes = string(groups[PosGermlineGenes])
		p.QValue = float32(qValue)
	}

	return
}
