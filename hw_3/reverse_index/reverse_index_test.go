package reverseindex

import (
	boollogic "hw_3/bool_logic"
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

		docs, err := newIndex.GetListDocumentsOnWord("new")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0}, docs)
	})

	t.Run("Базовое добавление и получение 2 документов", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)
		newIndex.AddDocument("Indexes are very pretty!", 1)

		docs, err := newIndex.GetListDocumentsOnWord("index")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1}, docs)
	})

	t.Run("Отсутствие документа", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)
		newIndex.AddDocument("Indexes are very pretty!", 1)

		docs, err := newIndex.GetListDocumentsOnWord("cuty")
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
		docs, err := newIndex.GetListDocumentsOnWord("new")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0}, docs)
	})

	t.Run("Базовое добавление и получение 2 документов со стоп-словами", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("The simple sentence", 2)
		newIndex.AddDocument("Simple indexes are very pretty!", 3)

		docs, err := newIndex.GetListDocumentsOnWord("simple")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{2, 3}, docs)
	})

	t.Run("Отсутствие документа со стоп-словами", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("Information search is very cool", 4)
		newIndex.AddDocument("The my index are very pretty!", 5)

		docs, err := newIndex.GetListDocumentsOnWord("The")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{}, docs)
	})
}

func Test_ReverseIndexWithWildcard(t *testing.T) {
	CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := NewInvertedIndex(&Params{
		Processor: processing.NewMyProcessing([]string{"The"}),
		Method:    Lemming,
		Mutex:     &commonMutex,
	})

	defer CleanupDb()
	t.Parallel()

	t.Run("Базовое добавление и получение документа с крайним левым wildcard", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)
		newIndex.AddDocument("Never", 1)

		docs, err := newIndex.GetListDocumentsOnWord("ne*")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1}, docs)
	})

	t.Run("Базовое добавление и получение документа с крайним правым wildcard", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("My new index is so nice", 0)
		newIndex.AddDocument("delice", 1)

		docs, err := newIndex.GetListDocumentsOnWord("*ec")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1}, docs)
	})

	t.Run("Базовое добавление и получение документа с full wildcard", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("Many time many time", 0)
		newIndex.AddDocument("There is many requests", 1)
		newIndex.AddDocument("The", 2)

		docs, err := newIndex.GetListDocumentsOnWord("t*e")
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1}, docs)
	})
}

func Test_ReverseIndexWithBaseBoolLogic(t *testing.T) {
	CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := NewInvertedIndex(&Params{
		Processor: processing.NewMyProcessing([]string{"The"}),
		Method:    Lemming,
		Mutex:     &commonMutex,
	})

	defer CleanupDb()
	t.Parallel()

	t.Run("Базовое добавление и получение по and", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("The my new index is so nice", 0)
		newIndex.AddDocument("The never", 1)

		newIndex.AddDocument("The my new index is so nice", 2)
		newIndex.AddDocument("The never", 3)

		docs, err := newIndex.GetListDocumentsOnBoolLogic(boollogic.New(boollogic.And, []string{"new", "nice"}, nil))
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 2}, docs)
	})

	t.Run("Базовое добавление и получение по or", func(t *testing.T) {
		t.Parallel()

		assert.Nil(t, err)
		newIndex.AddDocument("Something strange", 0)
		newIndex.AddDocument("The never", 1)

		newIndex.AddDocument("The question", 2)
		newIndex.AddDocument("One or two", 3)

		docs, err := newIndex.GetListDocumentsOnBoolLogic(boollogic.New(boollogic.Or, []string{"the", "never", "one", "two"}, nil))
		assert.Nil(t, err)
		assert.Equal(t, []uint32{1, 3}, docs)
	})
}
