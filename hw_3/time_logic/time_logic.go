package timelogic

import (
	"errors"
	"fmt"
	reverseindex "hw_3/reverse_index"
	"hw_3/reverse_index/processing"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/RoaringBitmap/roaring"
)

// Начнем с 2024-12-01 12:27 - 1733045312
//
//
//
//

const (
	DefPathFirst    = "./time_index_1"
	DefPathSecond   = "./time_index_2"
	layoutParseTime = "2006-05-16 21:30:02"
)

// TRangePredicate struct
type RangePredicate struct {
	MaxBitCount int
}

type TimeIndex struct {
	treeTimeStart    *reverseindex.InvertedIndex
	treeTimeEnd      *reverseindex.InvertedIndex
	mutex            *sync.Mutex
	AddedDocs        *roaring.Bitmap
	DocIDsByBitStart []*roaring.Bitmap
	DocIDsByBitEnd   []*roaring.Bitmap
}

type Params struct {
	Paths []string
}

func NewTRangePredicate(maxBitCount int) *RangePredicate {
	return &RangePredicate{MaxBitCount: maxBitCount}
}

func (trp *RangePredicate) GetPredicates(l, r uint32) [][]bool {
	var predicates [][]bool
	trp.getPredicatesImpl(0, math.MaxUint32, l, r, []bool{}, &predicates)
	return predicates
}

func (trp *RangePredicate) getPredicatesImpl(curLeft, curRight, requestedLeft, requestedRight uint32, path []bool, result *[][]bool) {
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

func NewTimeIndex(p *Params) (*TimeIndex, error) {
	if len(p.Paths) < 2 {
		return nil, fmt.Errorf("length of paths: %d < 2, exiting", len(p.Paths))
	}

	firstPath := p.Paths[0]
	secondPath := p.Paths[1]

	if firstPath == "" {
		firstPath = DefPathFirst
	}

	if secondPath == "" {
		secondPath = DefPathSecond
	}

	firstTree, err := reverseindex.NewInvertedIndex(&reverseindex.Params{
		Directory: firstPath,
		Processor: &processing.MyProcessing{},
	})
	if err != nil {
		return nil, fmt.Errorf("lsm.Open: %w", err)
	}

	secondTree, err := reverseindex.NewInvertedIndex(&reverseindex.Params{
		Directory: secondPath,
		Processor: &processing.MyProcessing{},
	})
	if err != nil {
		return nil, fmt.Errorf("lsm.Open: %w", err)
	}

	btStart := make([]*roaring.Bitmap, 32)
	btEnd := make([]*roaring.Bitmap, 32)

	for i := 0; i < 32; i++ {
		btStart[i] = roaring.NewBitmap()
		btEnd[i] = roaring.NewBitmap()
	}

	return &TimeIndex{
		treeTimeStart:    firstTree,
		treeTimeEnd:      secondTree,
		AddedDocs:        roaring.NewBitmap(),
		DocIDsByBitStart: btStart,
		DocIDsByBitEnd:   btEnd,
	}, nil
}

// AddDocumentOnStringTime добавляет документ в определенный временной диапазон, конец диапазона может быть не задан, время задается в формате 2006-05-12 12:23:20
func (t *TimeIndex) AddDocumentOnStringTime(document int, startTime, endTime string) error {
	if startTime == "" {
		return errors.New("start time is empty")
	}

	leftTime, err := time.Parse(time.RFC3339[:len(layoutParseTime)], startTime)
	if err != nil {
		return fmt.Errorf("invalid start time: %s, impossible to parse", startTime)
	}

	unixLeft := leftTime.Unix()
	unixRight := int64(1<<63 - 1)

	if endTime != "" {
		rightTime, err := time.Parse(time.RFC3339[:len(layoutParseTime)], startTime)
		if err != nil {
			return fmt.Errorf("invalid start time: %s, impossible to parse", startTime)
		}

		unixRight = rightTime.Unix()
	}

	err = t.treeTimeStart.WriteSafe([]byte(strconv.FormatInt(unixLeft, 10)), []byte(strconv.Itoa(int(document))))
	if err != nil {
		return fmt.Errorf("t.treeTime.WriteSafe: %w", err)
	}

	err = t.treeTimeEnd.WriteSafe([]byte(strconv.FormatInt(unixRight, 10)), []byte(strconv.Itoa(int(document))))
	if err != nil {
		return fmt.Errorf("t.treeTime.WriteSafe: %w", err)
	}

	return nil
}

func (t *TimeIndex) AddDocumentOnTimestamp(document int, startTime, endTime int64) error {
	t.AddedDocs.Add(uint32(document))
	for i := 0; i < 32; i++ {
		if (startTime>>(31-i))&1 != 0 {
			t.DocIDsByBitStart[i].Add(uint32(document))
		}
		if (endTime>>(31-i))&1 != 0 {
			t.DocIDsByBitEnd[i].Add(uint32(document))
		}
	}
	return nil
}

func (t *TimeIndex) FindDocsByInterval(intervalBegin, intervalEnd uint32) *roaring.Bitmap {
	predicates1 := NewTRangePredicate(32).GetPredicates(0, intervalEnd)
	intervalSuitableDocsStarts := t.evaluatePredicates(predicates1, t.DocIDsByBitStart)

	predicates2 := NewTRangePredicate(32).GetPredicates(intervalBegin, math.MaxUint32)
	intervalSuitableDocsEnds := t.evaluatePredicates(predicates2, t.DocIDsByBitEnd)

	intervalSuitableDocsStarts.And(intervalSuitableDocsEnds)
	return intervalSuitableDocsStarts
}

func (t *TimeIndex) FindDocsByTimePoint(timestamp uint32) *roaring.Bitmap {
	return t.FindDocsByInterval(timestamp, timestamp)
}

func (t *TimeIndex) evaluatePredicates(predicates [][]bool, bitSliceIndex []*roaring.Bitmap) *roaring.Bitmap {
	docs := roaring.NewBitmap()
	for _, predicate := range predicates {
		docs.Or(t.evaluatePredicate(predicate, bitSliceIndex))
	}
	return docs
}

func (t *TimeIndex) evaluatePredicate(predicate []bool, bitSliceIndex []*roaring.Bitmap) *roaring.Bitmap {
	docs := t.AddedDocs.Clone()
	for i := 0; i < len(predicate); i++ {
		if predicate[i] {
			docs.And(bitSliceIndex[i])
		} else {
			docs.AndNot(bitSliceIndex[i])
		}
	}
	return docs
}
