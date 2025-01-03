package reverseindex

import (
	"fmt"
	"hw_3/reverse_index/processing"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/RoaringBitmap/roaring"
	lsm "github.com/krasun/lsmtree"
)

type NormalizeType uint8

const (
	Lemming NormalizeType = iota
	Stemming
)

const DefPath = "./reverse_index_data"

type InvertedIndex struct {
	tree         *lsm.LSMTree
	processor    processing.Processing
	method       NormalizeType
	mutex        *sync.Mutex
	positionTree *lsm.LSMTree
}

type Params struct {
	Directory string
	Processor processing.Processing
	Method    NormalizeType
	Mutex     *sync.Mutex
}

func MakeTree(p *Params) (*lsm.LSMTree, error) {
	if _, err := os.Stat(p.Directory); os.IsNotExist(err) {
		errMk := os.Mkdir(p.Directory, 0755)
		if errMk != nil && !os.IsExist(errMk) {
			return nil, fmt.Errorf("os.Mkdir: %w", errMk)
		}

	}

	tree, err := lsm.Open(p.Directory)
	if err != nil {
		return nil, fmt.Errorf("lsm.Open: %w", err)
	}

	return tree, nil
}

func NewInvertedIndex(p *Params) (*InvertedIndex, error) {
	if p.Directory == "" {
		p.Directory = DefPath
	}

	tree, err := MakeTree(p)
	if err != nil {
		return nil, fmt.Errorf("MakeTree: %w", err)
	}

	return &InvertedIndex{
		tree:      tree,
		processor: p.Processor,
		method:    p.Method,
		mutex:     p.Mutex,
	}, nil
}

func (i *InvertedIndex) ProcessingText(text string, lang processing.Lang) ([]string, error) {
	if i.processor == nil {
		return []string{}, fmt.Errorf("processor is nil")
	}

	var newText []string
	var err error

	switch i.method {
	case Stemming:
		newText, err = i.processor.Stemming(text)
	case Lemming:
		newText, err = i.processor.Lemming(text, lang)
	default:
		log.Println("processing method are not chose...")
	}
	return newText, err
}

func (i *InvertedIndex) AddDocument(text string, index int, lang processing.Lang) error {
	textForDocument, err := i.ProcessingText(text, lang)
	if err != nil {
		return fmt.Errorf("i.ProcessingText: %w", err)
	}

	for _, word := range textForDocument {
		if len(word) <= 2 {
			continue
		}
		runes := []rune(word)
		prefStr := string(runes[0])
		sufStr := string(runes[len(runes)-1])

		for ind, sym := range runes[1 : len(runes)-1] {
			err := i.WriteWord(prefStr+"*", index)
			if err != nil {
				return fmt.Errorf("i.WriteWord: %w", err)
			}
			prefStr += string(sym)

			err = i.WriteWord("*"+sufStr, index)

			if err != nil {
				return fmt.Errorf("i.WriteWord: %w", err)
			}
			sufStr = string(runes[len(runes)-2-ind]) + sufStr
		}

		err := i.WriteWord(prefStr+"*", index)
		if err != nil {
			return fmt.Errorf("i.WriteWord: %w", err)
		}

		err = i.WriteWord("*"+sufStr, index)
		if err != nil {
			return fmt.Errorf("i.WriteWord: %w", err)
		}

		err = i.WriteWord(word, index)
		if err != nil {
			return fmt.Errorf("i.WriteWord: %w", err)
		}
	}
	return nil
}

func (i *InvertedIndex) WriteWord(word string, index int) error {
	docs, contains, err := i.ReadSafe(word)
	if err != nil {
		return fmt.Errorf("i.ReadSafe: %w", err)
	}

	bitmap := roaring.NewBitmap()
	if contains {
		err = bitmap.UnmarshalBinary(docs)
		if err != nil {
			return fmt.Errorf("bitmap.UnmarshalBinary: %w", err)
		}
	}

	bitmap.Add(uint32(index))
	docs, err = bitmap.MarshalBinary()
	if err != nil {
		return fmt.Errorf("bitmap.MarshalBinary: %w", err)
	}

	err = i.WriteSafe([]byte(word), docs)
	if err != nil {
		return fmt.Errorf("i.WriteSafe: %w", err)
	}

	return nil
}

func (i *InvertedIndex) GetBitmapDocuments(word string) (*roaring.Bitmap, error) {
	val, contains, err := i.ReadSafe(word)
	if err != nil {
		return nil, fmt.Errorf("i.ReadSafe: %w", err)
	}

	bitmap := roaring.NewBitmap()
	if !contains {
		return bitmap, nil
	}
	err = bitmap.UnmarshalBinary(val)
	if err != nil {
		return nil, err
	}

	return bitmap, nil
}

func (i *InvertedIndex) GetBitmapDocumentsOnBytes(bytes []byte) (*roaring.Bitmap, error) {
	val, contains, err := i.ReadSafeBytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("i.ReadSafeBytes: %w", err)
	}

	bitmap := roaring.NewBitmap()
	if !contains {
		return bitmap, nil
	}
	err = bitmap.UnmarshalBinary(val)
	if err != nil {
		return nil, err
	}

	return bitmap, nil
}

// GetBitmapsOnWord позволяет получить битмапу документов включая wildcard
func (i *InvertedIndex) GetBitmapsOnWord(word string) (*roaring.Bitmap, error) {
	var words []string

	if strings.Contains(word, "*") {
		words = strings.Split(word, "*")
	}

	if len(words) > 1 {
		if words[0] != "" && words[1] != "" {
			bitmapPrefix, err := i.GetBitmapDocuments(words[0] + "*")
			if err != nil {
				return nil, err
			}

			bitmapSuffix, err := i.GetBitmapDocuments("*" + words[1])
			if err != nil {
				return nil, err
			}

			bitmapPrefix.Intersects(bitmapSuffix)

			return bitmapPrefix, nil
		}
		if words[0] == "" {
			bitmapSuffix, err := i.GetBitmapDocuments("*" + words[1])
			if err != nil {
				return nil, err
			}
			return bitmapSuffix, nil
		} else {
			bitmap, err := i.GetBitmapDocuments(words[0] + "*")
			if err != nil {
				return nil, err
			}
			return bitmap, nil
		}
	}
	bitmap, err := i.GetBitmapDocuments(word)
	if err != nil {
		return nil, err
	}

	return bitmap, nil
}

// GetListDocumentsOnWord позволяет получить документы по слову, включая wildcard
func (i *InvertedIndex) GetListDocumentsOnWord(word string) ([]uint32, error) {
	bitmap, err := i.GetBitmapsOnWord(word)
	if err != nil {
		return nil, fmt.Errorf("i.GetBitmapsOnWord: %w", err)
	}
	return bitmap.ToArray(), nil
}

func (i *InvertedIndex) WriteSafe(bytesWord, docs []byte) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	err := i.tree.Put(bytesWord, docs)
	if err != nil {
		return fmt.Errorf("i.tree.Put: %w", err)
	}
	return nil
}

func (i *InvertedIndex) ReadSafe(word string) ([]byte, bool, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	val, contains, err := i.tree.Get([]byte(word))
	return val, contains, err
}

func (i *InvertedIndex) ReadSafeBytes(word []byte) ([]byte, bool, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	val, contains, err := i.tree.Get(word)
	return val, contains, err
}

func (i *InvertedIndex) WriteBytes(bytes []byte, index int) error {
	docs, contains, err := i.ReadSafeBytes(bytes)
	if err != nil {
		return fmt.Errorf("i.ReadSafe: %w", err)
	}

	bitmap := roaring.NewBitmap()
	if contains {
		err = bitmap.UnmarshalBinary(docs)
		if err != nil {
			return fmt.Errorf("bitmap.UnmarshalBinary: %w", err)
		}
	}

	bitmap.Add(uint32(index))
	docs, err = bitmap.MarshalBinary()
	if err != nil {
		return fmt.Errorf("bitmap.MarshalBinary: %w", err)
	}

	err = i.WriteSafe(bytes, docs)
	if err != nil {
		return fmt.Errorf("i.WriteSafe: %w", err)
	}

	return nil
}
