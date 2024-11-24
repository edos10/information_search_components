package boollogic

import (
	reverseindex "hw_3/reverse_index"
	"hw_3/reverse_index/processing"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReverseIndexWithBaseBoolLogic(t *testing.T) {
	reverseindex.CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := reverseindex.NewInvertedIndex(&reverseindex.Params{
		Processor: processing.NewMyProcessing([]string{"The"}),
		Method:    reverseindex.Lemming,
		Mutex:     &commonMutex,
	})

	defer reverseindex.CleanupDb()
	t.Parallel()

	t.Run("Базовое добавление и получение по and", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(And, []string{"Language", "Programming", "Types"}, nil)

		assert.Nil(t, err)
		newIndex.AddDocument("Language programming types", 0, processing.EN)
		newIndex.AddDocument("Language programming business", 1, processing.EN)

		newIndex.AddDocument("Language Python", 2, processing.EN)
		newIndex.AddDocument("Language programming broken type", 3, processing.EN)

		docs, err := newBoolIndex.Search(newIndex, processing.EN)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 3}, docs)
	})

	t.Run("Базовое добавление и получение по or", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(Or, []string{"something", "strange", "never"}, nil)

		assert.Nil(t, err)
		newIndex.AddDocument("Something strange", 0, processing.EN)
		newIndex.AddDocument("The never", 1, processing.EN)

		newIndex.AddDocument("The question", 2, processing.EN)
		newIndex.AddDocument("One or two", 3, processing.EN)

		docs, err := newBoolIndex.Search(newIndex, processing.EN)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1}, docs)
	})

	t.Run("Поиск по булевой формуле глубины 1", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(And, []string{"Language", "Programming", "Today"}, []*Node{
			{
				Operation: Or,
				Words:     []string{"Golang", "Python"},
			},
		})

		assert.Nil(t, err)
		newIndex.AddDocument("The programming languages today are very interesting. For example, Python", 0, processing.EN)
		newIndex.AddDocument("The programming language Python is one of the most popular", 1, processing.EN)

		docs, err := newBoolIndex.Search(newIndex, processing.EN)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0}, docs)
	})
}

func Test_ReverseIndexWithMediumBoolLogic(t *testing.T) {

	reverseindex.CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := reverseindex.NewInvertedIndex(&reverseindex.Params{
		Processor: processing.NewMyProcessing([]string{"The"}),
		Method:    reverseindex.Lemming,
		Mutex:     &commonMutex,
	})

	defer reverseindex.CleanupDb()
	t.Parallel()

	t.Run("Поиск по булевой формуле глубины 2", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(Or, []string{"Яйца", "Молоко", "Творог"}, []*Node{
			{
				Operation: And,
				Words:     []string{"Огурцы", "Помидоры", "Мясо"},
				Nodes: []*Node{
					{
						Operation: Or,
						Words:     []string{"Перец", "Баклажаны"},
					},
				},
			},
		})

		assert.Nil(t, err)
		newIndex.AddDocument("Для приготовления нужны яйца, масло и творог 2 пачки, а также мука и соль", 0, processing.RU)
		newIndex.AddDocument("В этом рецепте нужно мясо, молоко, перец, яйца, оливки и соль", 1, processing.RU)
		newIndex.AddDocument("Мы возьмем творог, огурцы. Добавим помидоры к салату. Начнем готовить мясо и добавим к ним баклажаны", 2, processing.RU)

		docs, err := newBoolIndex.Search(newIndex, processing.RU)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 2}, docs)
	})

	t.Run("Поиск по булевой формуле глубины 2 с AND в начале", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(And, []string{"Творог", "Молоко"}, []*Node{
			{
				Operation: Or,
				Words:     []string{"Огурцы", "Помидоры", "Мясо"},
				Nodes: []*Node{
					{
						Operation: Or,
						Words:     []string{"Перец", "Баклажаны"},
					},
				},
			},
		})

		assert.Nil(t, err)
		newIndex.AddDocument("Для приготовления нужны яйца, масло и творог 2 пачки, а также мука и соль", 0, processing.RU)
		newIndex.AddDocument("В этом рецепте нужно мясо, молоко, перец, яйца, оливки и соль", 1, processing.RU)
		newIndex.AddDocument("Мы возьмем творог и молоко. Добавим помидоры к салату. Начнем готовить мясо и добавим к ним баклажаны", 2, processing.RU)

		docs, err := newBoolIndex.Search(newIndex, processing.RU)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{2}, docs)
	})
}

func Test_ReverseIndexCornerCases(t *testing.T) {

	reverseindex.CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := reverseindex.NewInvertedIndex(&reverseindex.Params{
		Processor: processing.NewMyProcessing([]string{"The"}),
		Method:    reverseindex.Lemming,
		Mutex:     &commonMutex,
	})

	defer reverseindex.CleanupDb()
	t.Parallel()

	t.Run("Обработка русского текста английским лемматизатором", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(Or, []string{"Яйца", "Молоко", "Творог"}, []*Node{
			{
				Operation: And,
				Words:     []string{"Огурцы", "Помидоры", "Мясо"},
				Nodes: []*Node{
					{
						Operation: Or,
						Words:     []string{"Перец", "Баклажаны"},
					},
				},
			},
		})

		assert.Nil(t, err)
		newIndex.AddDocument("Для приготовления нужны яйца, масло и творог 2 пачки, а также мука и соль", 0, processing.EN)
		newIndex.AddDocument("В этом рецепте нужно мясо, молоко, перец, яйца, оливки и соль", 1, processing.RU)
		newIndex.AddDocument("Мы возьмем творог, огурцы. Добавим помидоры к салату. Начнем готовить мясо и добавим к ним баклажаны", 2, processing.RU)

		docs, err := newBoolIndex.Search(newIndex, processing.RU)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 2}, docs)
	})
}

func Test_ReverseIndexWildcardWithBool(t *testing.T) {

	reverseindex.CleanupDb()

	commonMutex := sync.Mutex{}
	newIndex, err := reverseindex.NewInvertedIndex(&reverseindex.Params{
		Processor: processing.NewMyProcessing([]string{"Под"}),
		Method:    reverseindex.Lemming,
		Mutex:     &commonMutex,
	})

	defer reverseindex.CleanupDb()
	t.Parallel()

	t.Run("Wildcard с булевой формулой", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(Or, []string{"*но", "Стог"}, nil)

		assert.Nil(t, err)
		newIndex.AddDocument("Окно", 0, processing.RU)
		newIndex.AddDocument("Сено", 1, processing.RU)
		newIndex.AddDocument("Стог сена", 2, processing.RU)

		docs, err := newBoolIndex.Search(newIndex, processing.RU)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 1, 2}, docs)
	})

	t.Run("Wildcard с булевой формулой по AND", func(t *testing.T) {
		t.Parallel()

		newBoolIndex := New(And, []string{"*но", "нужен"}, nil)

		assert.Nil(t, err)
		newIndex.AddDocument("Окно нужно", 0, processing.RU)
		newIndex.AddDocument("Сено", 1, processing.RU)
		newIndex.AddDocument("Стог сена нужен", 2, processing.RU)

		docs, err := newBoolIndex.Search(newIndex, processing.RU)
		assert.Nil(t, err)
		assert.Equal(t, []uint32{0, 2}, docs)
	})
}
