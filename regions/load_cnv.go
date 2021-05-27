package regions

import (
	"errors"
	"fmt"
	"github.com/overerd/clincnv_overlapper/models/clincnv"
)

func (s *Overlapper) loadCNVFile(name, path string) (file *clincnv.TableFile, err error) {
	file = clincnv.LoadClinCNVTable(clincnv.TableFileParserOptions{
		Name: name,
		Path: path,

		BufferSize: s.options.BufferSize,
	})

	err = file.Load()

	if err != nil {
		fmt.Printf("\n [!] %s", path)
		return
	}

	fmt.Printf("\n [*] %s (%d)", path, len(file.Items))

	return
}

func (s *Overlapper) loadCNVFiles(paths map[string]string) (err error) {
	if len(paths) < 2 && !s.options.SingleRunMode {
		err = errors.New("will not calculate overlap for just 1 file")
		return
	}

	s.files = make([]*clincnv.TableFile, len(paths))

	i := 0

	for name, path := range paths {
		file, err := s.loadCNVFile(name, path)

		if err != nil {
			return err
		}

		s.files[i] = file

		i++
	}

	return
}
