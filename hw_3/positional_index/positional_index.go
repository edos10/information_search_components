package positional_index

import (
	"encoding/json"
	"fmt"
	reverseindex "hw_3/reverse_index"
	"hw_3/reverse_index/processing"
	"log"
	"strings"
	"sync"

	"github.com/krasun/lsmtree"
)

const defPathPos = "./pos_index_data"

type PositionalIndex struct {
	treePositions  *lsmtree.LSMTree
	mutex          *sync.Mutex
	processor      processing.Processing
	method         reverseindex.NormalizeType
	lastCountWords map[int]int
}

type Params struct {
	Directory string
	Processor processing.Processing
	Method    reverseindex.NormalizeType
	Mutex     *sync.Mutex
}

type QueryResult struct {
	Documents [][]uint32 `json:"documents"`
}

func NewPosIndex(p *Params) (*PositionalIndex, error) {
	tree, err := reverseindex.MakeTree(&reverseindex.Params{
		Directory: defPathPos,
	})
	if err != nil {
		return nil, fmt.Errorf("reverseindex.MakeTree: %w", err)
	}
	return &PositionalIndex{
		treePositions:  tree,
		processor:      p.Processor,
		method:         p.Method,
		mutex:          p.Mutex,
		lastCountWords: make(map[int]int),
	}, nil
}

func (i *PositionalIndex) ProcessingText(text string, lang processing.Lang) ([]string, error) {
	if i.processor == nil {
		return []string{}, fmt.Errorf("processor is nil")
	}

	for _, symbol := range "/?!*.)({}[]:," {
		if strings.Contains("?!.,", string(symbol)) {
			text = strings.ReplaceAll(text, string(symbol), " ")
		} else {
			text = strings.ReplaceAll(text, string(symbol), "")
		}
	}

	var newText []string
	var err error

	switch i.method {
	case reverseindex.Stemming:
		newText, err = i.processor.Stemming(text)
	case reverseindex.Lemming:
		newText, err = i.processor.Lemming(text, lang)
	default:
		newText = strings.Split(text, " ")
		log.Println("processing method are not chose...")
	}
	return newText, err
}

func (i *PositionalIndex) AddDocument(text string, index int, lang processing.Lang) error {
	textForDocument, err := i.ProcessingText(text, lang)
	if err != nil {
		return fmt.Errorf("i.ProcessingText: %w", err)
	}

	count := 0
	if val, ok := i.lastCountWords[index]; ok {
		count = val
	}

	for _, word := range textForDocument {
		err := i.WriteWord(word, index, uint32(count))
		if err != nil {
			return fmt.Errorf("i.GetPositionsWord: %w", err)
		}
		count++
	}
	i.lastCountWords[index] = count

	return nil
}

func (i *PositionalIndex) WriteWord(word string, index int, pos uint32) error {
	docs, contains, err := i.ReadSafe(word)
	if err != nil {
		return fmt.Errorf("i.ReadSafe: %w", err)
	}

	res := new(QueryResult)
	if contains {
		err = json.Unmarshal(docs, res)
		if err != nil {
			return fmt.Errorf("json.Unmarshal: %w", err)
		}

		fmt.Println("TRUE", word, index, pos, res)

		flag := false

		for i, doc := range res.Documents {
			if len(doc) < 2 {
				return fmt.Errorf("invalid storage data on '%s' word", word)
			} else {
				if uint32(index) == doc[0] {
					res.Documents[i] = append(res.Documents[i], pos)
					flag = true
				}
			}
		}

		if !flag {
			arr := []uint32{uint32(index), pos}
			res.Documents = append(res.Documents, arr)
		}
	} else {
		arr := []uint32{uint32(index), pos}
		res.Documents = append(res.Documents, arr)
	}

	docs, err = json.Marshal(res)
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}

	err = i.WriteSafe([]byte(word), docs)
	if err != nil {
		return fmt.Errorf("i.WriteSafe: %w", err)
	}

	return nil
}

func (i *PositionalIndex) GetPositionsWord(word string) (*QueryResult, error) {
	val, contains, err := i.ReadSafe(word)
	if err != nil {
		return nil, fmt.Errorf("i.ReadSafe: %w", err)
	}

	res := new(QueryResult)
	if !contains {
		return &QueryResult{}, nil
	}
	err = json.Unmarshal(val, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (i *PositionalIndex) WriteSafe(bytesWord, docs []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	err := i.treePositions.Put(bytesWord, docs)
	if err != nil {
		return fmt.Errorf("i.tree.Put: %w", err)
	}
	return nil
}

func (i *PositionalIndex) ReadSafe(word string) ([]byte, bool, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	val, contains, err := i.treePositions.Get([]byte(word))
	return val, contains, err
}
