package boollogic

import (
	"fmt"
	reverseindex "hw_3/reverse_index"
	"hw_3/reverse_index/processing"
	"strings"

	"github.com/RoaringBitmap/roaring"
)

type Operation uint8

const (
	And Operation = iota
	Or
)

type Node struct {
	Words     []string
	Nodes     []*Node
	Operation Operation
	Value     uint8
}

func New(operation Operation, words []string, node []*Node) *Node {
	return &Node{
		Operation: operation,
		Words:     words,
		Nodes:     node,
	}
}

func (n *Node) Search(index *reverseindex.InvertedIndex, lang processing.Lang) ([]uint32, error) {
	bitmapDocs, err := n.SearchBitmaps(index, lang)
	if err != nil {
		return nil, fmt.Errorf("n.SearchBitmaps: %w", err)
	}

	return bitmapDocs.ToArray(), nil
}

func (n *Node) SearchBitmaps(index *reverseindex.InvertedIndex, lang processing.Lang) (*roaring.Bitmap, error) {
	ans := roaring.NewBitmap()

	for i, word := range n.Words {
		word, err := processing.LemmingWord(strings.ToLower(word), lang)
		if err != nil {
			return nil, fmt.Errorf("processing.LemmingWord: %w", err)
		}
		n.Words[i] = word
	}

	switch n.Operation {
	case Or:
		for _, word := range n.Words {
			bitmaps, err := index.GetBitmapsOnWord(word)
			if err != nil {
				return nil, fmt.Errorf("index.GetBitmapsOnWord: %w", err)
			}
			ans.Or(bitmaps)
		}

		for _, node := range n.Nodes {
			bitmaps, err := node.SearchBitmaps(index, lang)
			if err != nil {
				return nil, fmt.Errorf("node.SearchBitmaps: %w", err)
			}
			ans.Or(bitmaps)
		}
	case And:
		if len(n.Words) > 0 {
			bitmaps, err := index.GetBitmapDocuments(n.Words[0])
			if err != nil {
				return nil, fmt.Errorf("index.GetBitmapsOnWord: %w", err)
			}
			ans.Or(bitmaps)
		}

		for _, word := range n.Words {
			bitmaps, err := index.GetBitmapsOnWord(word)
			if err != nil {
				return nil, fmt.Errorf("index.GetBitmapsOnWord: %w", err)
			}
			ans.And(bitmaps)
		}

		if len(ans.ToArray()) == 0 && len(n.Nodes) > 0 {
			bitmaps, err := n.Nodes[0].SearchBitmaps(index, lang)
			if err != nil {
				return nil, fmt.Errorf("node.SearchBitmaps: %w", err)
			}
			ans.Or(bitmaps)
		}

		for _, node := range n.Nodes {
			bitmaps, err := node.SearchBitmaps(index, lang)
			if err != nil {
				return nil, fmt.Errorf("node.SearchBitmaps: %w", err)
			}
			ans.And(bitmaps)
		}
	}

	return ans, nil
}
