package reverseindex

import (
	"errors"
	"fmt"
	"hw_3/reverse_index/processing"
	"os"
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
	tree      *lsm.LSMTree
	processor processing.Processing
	method    NormalizeType
	mutex     *sync.Mutex
}

type Params struct {
	Directory string
	Processor processing.Processing
	Method    NormalizeType
	Mutex     *sync.Mutex
}

func NewInvertedIndex(p *Params) (*InvertedIndex, error) {
	if p.Directory == "" {
		p.Directory = DefPath
	}

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

	return &InvertedIndex{
		tree:      tree,
		processor: p.Processor,
		method:    p.Method,
		mutex:     p.Mutex,
	}, nil
}

func (i *InvertedIndex) ProcessingText(text string) ([]string, error) {
	if i.processor == nil {
		return []string{}, fmt.Errorf("processor is nil")
	}

	var newText []string
	var err error

	switch i.method {
	case Stemming:
		newText, err = i.processor.Stemming(text)
	case Lemming:
		newText, err = i.processor.Lemming(text)
	default:
		err = errors.New("unknown type of normalizing")
	}
	return newText, err
}

func (i *InvertedIndex) AddDocument(text string, index int) error {
	textForDocument, err := i.ProcessingText(text)
	if err != nil {
		return fmt.Errorf("i.ProcessingText: %w", err)
	}

	for _, word := range textForDocument {
		bytesWord := []byte(word)
		docs, contains, err := i.tree.Get(bytesWord)
		if err != nil {
			return fmt.Errorf("i.tree.Get: %w", err)
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

		err = i.tree.Put(bytesWord, docs)
		if err != nil {
			return fmt.Errorf("i.tree.Put: %w", err)
		}
	}
	return nil
}

func (i *InvertedIndex) GetBitmapDocuments(word string) (*roaring.Bitmap, error) {
	val, contains, err := i.ReadSafe(word)
	if err != nil {
		return nil, fmt.Errorf("i.tree.Get: %w", err)
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

func (i *InvertedIndex) GetListDocuments(word string) ([]uint32, error) {
	bitmap, err := i.GetBitmapDocuments(word)
	if err != nil {
		return []uint32{}, err
	}
	uint32Array := bitmap.ToArray()
	return uint32Array, nil
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
