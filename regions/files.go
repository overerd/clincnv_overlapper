package regions

import "github.com/overerd/clincnv_overlapper/models/clincnv"

func (s *Overlapper) GetFiles() (files *[]*clincnv.TableFile) {
	return &s.files
}
