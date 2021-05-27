package bed

import (
	"bufio"
	"os"
	"regexp"
)

type File struct {
	Options     FileParserOptions  `json:"options"`
	Chromosomes map[string][]*Item `json:"chromosomes"`
	Items       []Item             `json:"items"`

	Map map[string]map[uint64]map[uint64]*Item `json:"-"`

	ReversedMap map[string]map[uint64]map[uint64]*Item `json:"-"`

	re *regexp.Regexp

	geneColumnIndex uint
}

func populateGeneMap(geneMap *map[string]bool, item *Item) {
	for _, gene := range (*item).Genes {
		(*geneMap)[gene] = true
	}
}

func gatherForwardMappedGenes(items *map[uint64]map[uint64]*Item, geneMap *map[string]bool, start, end uint64) {
	for s, left := range *items {
		if s > end {
			continue
		}

		for e, item := range left {
			if e > start {
				populateGeneMap(geneMap, item)
			}
		}
	}
}

func gatherReverseMappedGenes(items *map[uint64]map[uint64]*Item, geneMap *map[string]bool, start, end uint64) {
	for e, right := range *items {
		if e < start {
			continue
		}

		for s, item := range right {
			if s < end {
				populateGeneMap(geneMap, item)
			}
		}
	}
}

func (b *File) SelectAllGenes(chr string, start, end uint64) (genes []string) {
	geneMap := make(map[string]bool)

	items := b.Map[chr]
	reversedItems := b.Map[chr]

	gatherForwardMappedGenes(&items, &geneMap, start, end)
	gatherReverseMappedGenes(&reversedItems, &geneMap, start, end)

	for gene := range geneMap {
		genes = append(genes, gene)
	}

	return
}

func (b *File) lookForClosestItem(sourceMap *map[uint64]map[uint64]*Item, start, end uint64) (found bool, item *Item) {
	var left map[uint64]*Item
	var right *Item

	for l := range *sourceMap {
		if l <= start {
			left = (*sourceMap)[l]

			for r := range left {
				if r >= end {
					right = left[r]
				}
			}
		}
	}

	if right == nil {
		return
	}

	found = true
	item = right

	return
}

func (b *File) FindClosestItem(chr string, start, end uint64) (found bool, item *Item) {
	items := b.Map[chr]
	reversedItems := b.ReversedMap[chr]

	found, item = b.lookForClosestItem(&items, start, end)

	if !found {
		found, item = b.lookForClosestItem(&reversedItems, start, end)
	}

	return
}

func (b *File) Load() (err error) {
	b.Items = []Item{}
	b.Chromosomes = make(map[string][]*Item)

	b.re, err = regexp.Compile(reLine)

	if err != nil {
		return err
	}

	err = b.parseFile()

	return err
}

func (b *File) extendMap(item *Item) {
	if _, ok := b.Map[item.Chr]; !ok {
		b.Map[item.Chr] = make(map[uint64]map[uint64]*Item)
	}

	if _, ok := b.Map[item.Chr][item.Start]; !ok {
		b.Map[item.Chr][item.Start] = make(map[uint64]*Item)
	}

	b.Map[item.Chr][item.Start][item.End] = item
}

func (b *File) extendReverseMap(item *Item) {
	if _, ok := b.ReversedMap[item.Chr]; !ok {
		b.ReversedMap[item.Chr] = make(map[uint64]map[uint64]*Item)
	}

	if _, ok := b.ReversedMap[item.Chr][item.End]; !ok {
		b.ReversedMap[item.Chr][item.End] = make(map[uint64]*Item)
	}

	b.ReversedMap[item.Chr][item.End][item.Start] = item
}

func (b *File) readLine(line []byte) (err error) {
	if line[0] == byte('#') {
		return
	}

	matches := b.re.FindSubmatch(line)

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

	err = item.ItemData.Fill(matches, &b.Options)

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

	b.extendMap(&item)
	b.extendReverseMap(&item)

	return
}

func (b *File) parseFile() (err error) {
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

	b.Map = make(map[string]map[uint64]map[uint64]*Item)

	b.ReversedMap = make(map[string]map[uint64]map[uint64]*Item)

	for scanner.Scan() {
		err = b.readLine(scanner.Bytes())

		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return
}

func LoadBEDFile(options FileParserOptions) (bed *File) {
	bed = &File{
		Options: options,
	}

	bed.geneColumnIndex = options.GeneColumnIndex

	return
}
