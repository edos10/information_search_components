package reverseindex

import (
	"hw_3/reverse_index/processing"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReverseIndexWithoutStopWords(t *testing.T) {
	CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := NewInvertedIndex(&Params{
		Processor: processing.NewMyProcessing([]string{}),
		Method:    Lemming,
		Mutex:     &commonMutex,
	})

	defer CleanupDb()
	t.Parallel()

	t.Run("Базовое добавление и получение документа", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)

		docs, err := newIndex.GetListDocuments("my")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0}, docs)
	})

	t.Run("Базовое добавление и получение 2 документов", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)
		newIndex.AddDocument("Indexes are very pretty!", 1)

		docs, err := newIndex.GetListDocuments("index")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1}, docs)
	})

	t.Run("Отсутствие документа", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)
		newIndex.AddDocument("Indexes are very pretty!", 1)

		docs, err := newIndex.GetListDocuments("cuty")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{}, docs)
	})
}

func Test_ReverseIndexWithStopWords(t *testing.T) {
	CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := NewInvertedIndex(&Params{
		Processor: processing.NewMyProcessing([]string{"The"}),
		Method:    Lemming,
		Mutex:     &commonMutex,
	})

	defer CleanupDb()
	t.Parallel()

	t.Run("Базовое добавление и получение документа со стоп-словами", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)

		docs, err := newIndex.GetListDocuments("my")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0}, docs)
	})

	t.Run("Базовое добавление и получение 2 документов со стоп-словами", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("The simple sentence", 2)
		newIndex.AddDocument("Simple indexes are very pretty!", 3)

		docs, err := newIndex.GetListDocuments("simple")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{2, 3}, docs)
	})

	t.Run("Отсутствие документа со стоп-словами", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("Information search is very cool", 4)
		newIndex.AddDocument("The my index are very pretty!", 5)

		docs, err := newIndex.GetListDocuments("The")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{}, docs)
	})
}

func Test_OvermanyRequests(t *testing.T) {

}
