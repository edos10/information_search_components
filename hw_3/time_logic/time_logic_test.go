package timelogic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BaseTimeIndexWithIntervals(t *testing.T) {
	index, err := NewTimeIndex(&Params{
		Paths: []string{"./one", "./two"},
	})
	require.NoError(t, err)

	index.AddDocumentOnTimestamp(3, 30, 40)
	index.AddDocumentOnTimestamp(5, 10, 20)
	index.AddDocumentOnTimestamp(6, 1, 25)

	docs := index.FindDocsByInterval(7, 29)

	assert.Equal(t, []uint32{5, 6}, docs.ToArray())
	docsByTimePoint := index.FindDocsByTimePoint(10)

	assert.Equal(t, []uint32{5, 6}, docsByTimePoint.ToArray())
}

func Test_SearchTimeIndexOnSeveralDocs(t *testing.T) {
	index, err := NewTimeIndex(&Params{
		Paths: []string{"./one", "./two"},
	})
	require.NoError(t, err)

	index.AddDocumentOnTimestamp(1, 1, 10)
	index.AddDocumentOnTimestamp(2, 10, 20)
	index.AddDocumentOnTimestamp(3, 30, 50)

	docs := index.FindDocsByInterval(1, 50)

	assert.Equal(t, []uint32{1, 2, 3}, docs.ToArray())

	docs = index.FindDocsByInterval(20, 50)

	assert.Equal(t, []uint32{2, 3}, docs.ToArray())
}
