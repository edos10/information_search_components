package positional_index

import (
	reverseindex "hw_3/reverse_index"
	"hw_3/reverse_index/processing"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_BasePositionalIndex(t *testing.T) {
	t.Parallel()

	defer reverseindex.CleanupDb()

	t.Run("Базовое добавление и получение списка позиций", func(t *testing.T) {
		newIndex, err := NewPosIndex(&Params{
			Directory: defPathPos,
			Processor: &processing.MyProcessing{},
			Method:    3,
			Mutex:     &sync.Mutex{},
		})
		assert.Nil(t, err)
		err = newIndex.AddDocument("text a b c d text d text", 0, processing.EN)
		assert.Nil(t, err)

		q, err := newIndex.GetPositionsWord("text")
		assert.Nil(t, err)
		assert.Equal(t, [][]uint32{
			{0, 0, 5, 7},
		}, q.Documents)

	})

	t.Run("Повторное добавление и получение списка позиций", func(t *testing.T) {
		newIndex, err := NewPosIndex(&Params{
			Directory: defPathPos,
			Processor: &processing.MyProcessing{},
			Method:    3,
			Mutex:     &sync.Mutex{},
		})
		assert.Nil(t, err)
		err = newIndex.AddDocument("text a b c d text d text", 1, processing.EN)
		assert.Nil(t, err)

		q, err := newIndex.GetPositionsWord("text")
		assert.Nil(t, err)
		assert.Equal(t, [][]uint32{
			{0, 0, 5, 7},
			{1, 0, 5, 7},
		}, q.Documents)

		err = newIndex.AddDocument("rand text", 1, processing.EN)
		assert.Nil(t, err)

		q, err = newIndex.GetPositionsWord("text")
		assert.Nil(t, err)
		assert.Equal(t, [][]uint32{
			{0, 0, 5, 7},
			{1, 0, 5, 7, 9},
		}, q.Documents)
	})
}
