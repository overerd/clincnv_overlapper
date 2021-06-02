package clincnv

import (
	"bufio"
	"os"
	"regexp"
)

type TableFile struct {
	Options     TableFileParserOptions `json:"options"`
	Chromosomes map[string][]*Item     `json:"chromosomes"`
	Items       []Item                 `json:"items"`

	Map map[string]map[uint64]map[uint64]*Item `json:"-"`

	ReversedMap map[string]map[uint64]map[uint64]*Item `json:"-"`

	reGermline *regexp.Regexp
	reSomatic  *regexp.Regexp
}

func (b *TableFile) FindClosestItem(chr string, start, end uint64) (found bool, item *Item) {
	items := b.Map[chr]

	var right *Item

	for l, left := range items {
		if l <= start {
			for r := range left {
				if r >= end || r >= start {
					right = left[r]
				}
			}
		}
	}

	if right == nil {
		items = b.ReversedMap[chr]

		for l, left := range items {
			if l >= end {
				for r := range left {
					if r <= start || r <= end {
						right = left[r]
					}
				}
			}
		}
	}

	if right != nil {
		found = true
		item = right
	}

	return
}

func (b *TableFile) FindIntervalItems(chr string, start, end uint64) (results []*Item) {
	items := b.Map[chr]

	for l, left := range items {
		if l <= start {
			for r := range left {
				if r >= end || r >= start {
					results = append(results, left[r])
				}
			}
		}
	}

	return
}

func (b *TableFile) Load() (err error) {
	b.Items = []Item{}
	b.Chromosomes = make(map[string][]*Item)

	b.reGermline, err = regexp.Compile(reLineGermline)

	if err != nil {
		return err
	}

	b.reSomatic, err = regexp.Compile(reLineSomatic)

	if err != nil {
		return err
	}

	err = b.parseFile(b.Options.BufferSize)

	return err
}

func (b *TableFile) readLine(line []byte) (err error) {
	if line[0] == byte('#') {
		return
	}

	matches := b.reSomatic.FindSubmatch(line)

	isSomatic := true

	if matches == nil {
		isSomatic = false

		matches = b.reGermline.FindSubmatch(line)
	}

	if matches == nil {
		return
	}
	if len(matches) <= 1 {
		return
	}

	item := Item{}

	err = item.ChromosomeRegionData.Fill(matches)

	if err != nil {
		return
	}

	err = item.ChromosomeData.Fill(matches, isSomatic)

	if err != nil {
		return
	}

	b.Items = append(b.Items, item)

	if _, ok := b.Chromosomes[item.Chr]; ok {
		b.Chromosomes[item.Chr] = append(b.Chromosomes[item.Chr], &item)
	} else {
		b.Chromosomes[item.Chr] = []*Item{
			&item,
		}
	}

	if _, ok := b.Map[item.Chr]; !ok {
		b.Map[item.Chr] = make(map[uint64]map[uint64]*Item)
	}

	if _, ok := b.ReversedMap[item.Chr]; !ok {
		b.ReversedMap[item.Chr] = make(map[uint64]map[uint64]*Item)
	}

	if _, ok := b.Map[item.Chr][item.Start]; !ok {
		b.Map[item.Chr][item.Start] = make(map[uint64]*Item)
	}

	if _, ok := b.ReversedMap[item.Chr][item.End]; !ok {
		b.ReversedMap[item.Chr][item.End] = make(map[uint64]*Item)
	}

	b.Map[item.Chr][item.Start][item.End] = &item
	b.ReversedMap[item.Chr][item.End][item.Start] = &item

	return
}

func (b *TableFile) parseFile(bufferSize uint) (err error) {
	file, err := os.Open(b.Options.Path)

	if err != nil {
		return err
	}

	defer func(file *os.File) {
		e := file.Close()
		if e != nil {
			err = e
		}
	}(file)

	scanner := bufio.NewScanner(file)

	buffer := make([]byte, bufferSize)

	scanner.Buffer(buffer, int(bufferSize)*10)

	b.Map = make(map[string]map[uint64]map[uint64]*Item)
	b.ReversedMap = make(map[string]map[uint64]map[uint64]*Item)

	for scanner.Scan() {
		err = b.readLine(scanner.Bytes())

		if err != nil {
			return err
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return
}

func LoadClinCNVTable(options TableFileParserOptions) (bed *TableFile) {
	bed = &TableFile{
		Options: options,
	}

	return
}
