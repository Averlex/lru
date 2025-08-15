package lru

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	testValue = 42
)

func verifyListStructure(t *testing.T, l List[any], length int, front, back *ListItem[any]) {
	t.Helper()
	require.Equal(t, length, l.Len())
	require.Equal(t, front, l.Front())
	require.Equal(t, back, l.Back())
}

func verifyMoveToFront(t *testing.T, l List[any], length int, elem *ListItem[any]) {
	t.Helper()

	require.Equal(t, length, l.Len())
	// Comparing values since for the cases where method may reallocate memory.
	require.Equal(t, l.Front().Value, elem.Value)
	require.Nil(t, l.Front().Prev)
	require.Nil(t, l.Back().Next)
}

func emptyList(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name string
		test func(*testing.T)
	}{
		{"push operations", emptyListPushOperations},
		{"remove", emptyListRemove},
		{"move to front", emptyListMoveToFront},
	}

	for _, tC := range testCases {
		t.Run(tC.name, tC.test)
	}
}

func emptyListPushOperations(t *testing.T) {
	t.Helper()

	tests := []struct {
		name   string
		method func(List[any], any) *ListItem[any]
	}{
		{"push front", List[any].PushFront},
		{"push back", List[any].PushBack},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewList[any]()
			elem := tt.method(l, testValue)
			verifyListStructure(t, l, 1, elem, elem)
		})
	}
}

func emptyListRemove(t *testing.T) {
	t.Helper()

	l := NewList[any]()
	l.Remove(nil)
	verifyListStructure(t, l, 0, nil, nil)
}

func emptyListMoveToFront(t *testing.T) {
	t.Helper()

	l := NewList[any]()
	elem := &ListItem[any]{Value: testValue}
	l.MoveToFront(elem)

	verifyMoveToFront(t, l, 1, elem)
	require.Equal(t, l.Front(), l.Back())
}

func singleElementList(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name string
		test func(*testing.T)
	}{
		{"data structure", singleElementStructure},
		{"push operations", singleElementPushOperations},
		{"remove", singleElementRemove},
		{"move to front", singleElementMoveToFront},
	}

	for _, tC := range testCases {
		t.Run(tC.name, tC.test)
	}
}

func singleElementStructure(t *testing.T) {
	t.Helper()

	l := NewList[any]()
	elem := l.PushFront(testValue)
	verifyListStructure(t, l, 1, elem, elem)
}

func singleElementPushOperations(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name                   string
		method                 func(List[any], any) *ListItem[any]
		headMethod, tailMethod func(List[any]) *ListItem[any]
	}{
		{"push front", List[any].PushFront, List[any].Front, List[any].Back},
		{"push back", List[any].PushBack, List[any].Back, List[any].Front},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			l := NewList[any]()
			first := l.PushFront(testValue)
			second := tC.method(l, testValue)

			// Data structure testing.
			require.Equal(t, 2, l.Len())
			require.Equal(t, tC.headMethod(l), second)
			require.Equal(t, tC.tailMethod(l), first)

			// Linkage testing.
			require.Equal(t, l.Front().Next, l.Back())
			require.Equal(t, l.Back().Prev, l.Front())
		})
	}
}

func singleElementRemove(t *testing.T) {
	t.Helper()

	l := NewList[any]()
	elem := l.PushFront(testValue)
	l.Remove(elem)
	verifyListStructure(t, l, 0, nil, nil)
}

func singleElementMoveToFront(t *testing.T) {
	t.Helper()

	l := NewList[any]()
	elem := l.PushFront(testValue)
	l.MoveToFront(elem)

	verifyMoveToFront(t, l, 1, elem)
	require.Equal(t, l.Front(), l.Back())
}

func twoElementList(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name string
		test func(*testing.T)
	}{
		{"data structure", twoElementStructure},
		{"push operations", twoElementPushOperations},
		{"remove", twoElementRemove},
		{"move to front", twoElementMoveToFront},
	}

	for _, tC := range testCases {
		t.Run(tC.name, tC.test)
	}
}

func twoElementStructure(t *testing.T) {
	t.Helper()

	l := NewList[any]()
	first := l.PushFront(testValue)
	second := l.PushBack(testValue)

	verifyListStructure(t, l, 2, first, second)
	require.Equal(t, l.Front().Next, l.Back())
	require.Equal(t, l.Back().Prev, l.Front())
	require.Nil(t, l.Front().Prev)
	require.Nil(t, l.Back().Next)
}

func twoElementPushOperations(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name   string
		method func(List[any], any) *ListItem[any]
	}{
		{"push front", List[any].PushFront},
		{"push back", List[any].PushBack},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			l := NewList[any]()
			first := l.PushFront(testValue)
			second := l.PushBack(testValue)
			newElem := tC.method(l, testValue)

			if tC.name == "push front" {
				verifyListStructure(t, l, 3, newElem, second)
			} else {
				verifyListStructure(t, l, 3, first, newElem)
			}
			// Linkage testing.
			require.Equal(t, l.Front().Next, l.Back().Prev)
			require.Equal(t, l.Front().Next.Next, l.Back())
			require.Equal(t, l.Back().Prev.Prev, l.Front())
		})
	}
}

func twoElementRemove(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name string
		test func(*testing.T, List[any])
	}{
		{"remove first", removeFirstElement},
		{"remove last", removeLastElement},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			l := NewList[any]()
			l.PushFront(testValue)
			l.PushBack(testValue)
			tC.test(t, l)
			require.Nil(t, l.Front().Prev)
			require.Nil(t, l.Back().Next)
		})
	}
}

func removeFirstElement(t *testing.T, l List[any]) {
	t.Helper()

	second := l.Back()
	l.Remove(l.Front())
	verifyListStructure(t, l, 1, second, second)
}

func removeLastElement(t *testing.T, l List[any]) {
	t.Helper()

	first := l.Front()
	l.Remove(l.Back())
	verifyListStructure(t, l, 1, first, first)
}

func twoElementMoveToFront(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name string
	}{
		{"move first"},
		{"move last"},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			l := NewList[any]()
			first := l.PushFront(testValue)
			second := l.PushBack(testValue)

			if tC.name == "move first" {
				l.MoveToFront(first)
				verifyMoveToFront(t, l, 2, first)
				require.Equal(t, second, l.Back())
			} else {
				l.MoveToFront(second)
				verifyMoveToFront(t, l, 2, second)
				require.Equal(t, first, l.Back())
			}

			require.Equal(t, l.Front().Next, l.Back())
			require.Equal(t, l.Back().Prev, l.Front())
		})
	}
}

func complexBehavior(t *testing.T) {
	t.Helper()

	suite.Run(t, new(BehaviorTestSuite))
}

type BehaviorTestSuite struct {
	suite.Suite
	l        List[any]
	cycleLen int
	expected []int
}

func (s *BehaviorTestSuite) SetupTest() {
	s.l = NewList[any]()

	s.cycleLen = 10
	s.expected = make([]int, 0, s.cycleLen)

	// expected: [0, 10, 20, 30, 40, 50, 60, 70, 80, 90]
	for i := range s.cycleLen {
		s.l.PushBack(i * 10)
		s.expected = append(s.expected, i*10)
	}
}

func (s *BehaviorTestSuite) TestClearList() {
	s.Require().NotNil(s.l)
	s.Require().Equal(s.cycleLen, s.l.Len())
	s.Require().Equal(s.expected, s.getList(s.l))

	// Clearing
	for range s.cycleLen {
		s.l.Remove(s.l.Back())
	}

	s.Require().Equal(0, s.l.Len())
}

func (s *BehaviorTestSuite) TestCyclicMoveToFront() {
	s.Require().NotNil(s.l)
	s.Require().Equal(s.cycleLen, s.l.Len())
	s.Require().Equal(s.expected, s.getList(s.l))

	// Cyclic shift.
	for range s.cycleLen {
		s.l.MoveToFront(s.l.Back())
	}

	s.Require().Equal(s.expected, s.getList(s.l))
}

func (s *BehaviorTestSuite) getList(l List[any]) []int {
	elems := make([]int, 0, s.cycleLen)
	for i := l.Front(); i != nil; i = i.Next {
		elems = append(elems, i.Value.(int))
	}
	return elems
}

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList[any]()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList[any]()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("custom empty list", func(t *testing.T) { emptyList(t) })
	t.Run("list with a single element", func(t *testing.T) { singleElementList(t) })
	t.Run("list with 2 elements", func(t *testing.T) { twoElementList(t) })
	t.Run("complex behavior", func(t *testing.T) { complexBehavior(t) })
}
