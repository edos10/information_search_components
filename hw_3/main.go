package main

import (
	"fmt"
	"math"

	"github.com/RoaringBitmap/roaring"
)

// TRangePredicate struct
type TRangePredicate struct {
	MaxBitCount int
}

func NewTRangePredicate(maxBitCount int) *TRangePredicate {
	return &TRangePredicate{MaxBitCount: maxBitCount}
}

func (trp *TRangePredicate) GetPredicates(l, r uint32) [][]bool {
	var predicates [][]bool
	trp.getPredicatesImpl(0, math.MaxUint32, l, r, []bool{}, &predicates)
	return predicates
}

func (trp *TRangePredicate) getPredicatesImpl(curLeft, curRight, requestedLeft, requestedRight uint32, path []bool, result *[][]bool) {
	if curLeft > curRight || requestedLeft > requestedRight {
		return
	}
	if curLeft == requestedLeft && curRight == requestedRight {
		tempPath := make([]bool, len(path))
		copy(tempPath, path)
		*result = append(*result, tempPath)
		return
	}
	mid := (curLeft + curRight) / 2

	path = append(path, false)
	trp.getPredicatesImpl(curLeft, mid, requestedLeft, min(requestedRight, mid), path, result)
	path = path[:len(path)-1]

	path = append(path, true)
	trp.getPredicatesImpl(mid+1, curRight, max(requestedLeft, mid+1), requestedRight, path, result)
	path = path[:len(path)-1]

}

// TInvertedDateIntervalIndex struct
type TInvertedDateIntervalIndex struct {
	AddedDocs        *roaring.Bitmap
	DocIDsByBitStart []*roaring.Bitmap
	DocIDsByBitEnd   []*roaring.Bitmap
}

func NewTInvertedDateIntervalIndex() *TInvertedDateIntervalIndex {
	index := &TInvertedDateIntervalIndex{
		AddedDocs:        roaring.NewBitmap(),
		DocIDsByBitStart: make([]*roaring.Bitmap, 32),
		DocIDsByBitEnd:   make([]*roaring.Bitmap, 32),
	}

	for i := 0; i < 32; i++ {
		index.DocIDsByBitStart[i] = roaring.NewBitmap()
		index.DocIDsByBitEnd[i] = roaring.NewBitmap()
	}

	return index
}

func (tdii *TInvertedDateIntervalIndex) AddDocument(doc TDocument, intervalBegin, intervalEnd uint32) {
	tdii.AddedDocs.Add(doc.ID)
	for i := 0; i < 32; i++ {
		if (intervalBegin>>(31-i))&1 != 0 {
			tdii.DocIDsByBitStart[i].Add(doc.ID)
		}
		if (intervalEnd>>(31-i))&1 != 0 {
			tdii.DocIDsByBitEnd[i].Add(doc.ID)
		}
	}
}

func (tdii *TInvertedDateIntervalIndex) FindDocsByInterval(intervalBegin, intervalEnd uint32) *roaring.Bitmap {
	predicates1 := NewTRangePredicate(32).GetPredicates(0, intervalEnd)
	intervalSuitableDocsStarts := tdii.evaluatePredicates(predicates1, tdii.DocIDsByBitStart)

	predicates2 := NewTRangePredicate(32).GetPredicates(intervalBegin, math.MaxUint32)
	intervalSuitableDocsEnds := tdii.evaluatePredicates(predicates2, tdii.DocIDsByBitEnd)

	intervalSuitableDocsStarts.And(intervalSuitableDocsEnds)
	return intervalSuitableDocsStarts
}

func (tdii *TInvertedDateIntervalIndex) FindDocsByTimePoint(timestamp uint32) *roaring.Bitmap {
	return tdii.FindDocsByInterval(timestamp, timestamp)
}

func (tdii *TInvertedDateIntervalIndex) evaluatePredicates(predicates [][]bool, bitSliceIndex []*roaring.Bitmap) *roaring.Bitmap {
	docs := roaring.NewBitmap()
	for _, predicate := range predicates {
		docs.Or(tdii.evaluatePredicate(predicate, bitSliceIndex))
	}
	return docs
}

func (tdii *TInvertedDateIntervalIndex) evaluatePredicate(predicate []bool, bitSliceIndex []*roaring.Bitmap) *roaring.Bitmap {
	docs := tdii.AddedDocs.Clone()
	for i := 0; i < len(predicate); i++ {
		if predicate[i] {
			docs.And(bitSliceIndex[i])
		} else {
			docs.AndNot(bitSliceIndex[i])
		}
	}
	return docs
}

// Helper functions (you might have these elsewhere)
func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

// TDocument struct
type TDocument struct {
	ID uint32
}

func main() {
	index := NewTInvertedDateIntervalIndex()

	doc1 := TDocument{ID: 1}
	doc2 := TDocument{ID: 2}
	doc3 := TDocument{ID: 3}

	index.AddDocument(doc1, 30, 40)
	index.AddDocument(doc2, 10, 20)
	index.AddDocument(doc3, 1, 25)

	docs := index.FindDocsByInterval(7, 30)

	fmt.Println("Documents found:", docs)

	docsByTimePoint := index.FindDocsByTimePoint(10)
	fmt.Println("Documents by time point:", docsByTimePoint)
}
