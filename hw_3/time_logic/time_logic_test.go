package timelogic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHuinya(t *testing.T) {

	start := int64(21) // 10101
	for i := 0; i < 5; i++ {
		fmt.Println(i, 1<<i&start)
	}
}

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
