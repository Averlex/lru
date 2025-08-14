package lru

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("key1", 100)
		c.Set("key2", 200)
		c.Set("key3", 300)

		c.Clear()

		notInCacheChecks(t, &c, "key1")
		notInCacheChecks(t, &c, "key2")
		notInCacheChecks(t, &c, "key3")
	})

	t.Run("incorrect capacity", incorrectCapacity)
	t.Run("single element cache", cacheSingleItemSuite)
	t.Run("multi element cache", cacheMultiItemSuite)
	t.Run("eviction", cacheEvictionSuite)
	t.Run("stress", cacheStressSuite)
}

//nolint:revive
func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func notInCacheChecks(t *testing.T, c *Cache, key Key) {
	t.Helper()

	v, ok := (*c).Get(key)
	require.False(t, ok)
	require.Nil(t, v)
}

func incorrectCapacity(t *testing.T) {
	t.Helper()

	t.Run("zero capacity cache is unusable", func(t *testing.T) {
		c := NewCache(0)
		require.Nil(t, c)
	})

	t.Run("negative capacity cache is unusable", func(t *testing.T) {
		c := NewCache(-1)
		require.Nil(t, c)
	})
}

type CacheTestHelper struct {
	suite.Suite
	cache Cache
}

func (s *CacheTestHelper) isNotInCache(k Key) {
	v, ok := s.cache.Get(k)
	s.False(ok)
	s.Nil(v)
}

func (s *CacheTestHelper) isInCache(k Key, val any) {
	v, ok := s.cache.Get(k)
	s.True(ok)
	s.Equal(val, v)
}

func (s *CacheTestHelper) setExisting(k Key, val any) {
	wasInCache := s.cache.Set(k, val)
	s.True(wasInCache)
}

func (s *CacheTestHelper) setNew(k Key, val any) {
	wasInCache := s.cache.Set(k, val)
	s.False(wasInCache)
}

type SingleItemCacheSuite struct {
	CacheTestHelper
}

func (s *SingleItemCacheSuite) SetupTest() {
	s.cache = NewCache(1)
}

func (s *SingleItemCacheSuite) TestSetToEmpty() {
	s.setNew("key1", 100)
}

func (s *SingleItemCacheSuite) TestSetWithUpdate() {
	s.setNew("key1", 100)
	s.setExisting("key1", 200)
}

func (s *SingleItemCacheSuite) TestSetToFull() {
	s.setNew("key1", 100)
	s.setNew("key2", 200)
}

func (s *SingleItemCacheSuite) TestGetFromEmpty() {
	s.isNotInCache("key1")
}

func (s *SingleItemCacheSuite) TestGetFromFilled() {
	s.cache.Set("key1", 100)
	s.isInCache("key1", 100)
}

func (s *SingleItemCacheSuite) TestGetNonExistent() {
	s.cache.Set("key1", 100)
	s.isNotInCache("key2")
}

func (s *SingleItemCacheSuite) TestClearEmpty() {
	s.cache.Clear()
	s.isNotInCache("key1")
}

func (s *SingleItemCacheSuite) TestClearFilled() {
	s.cache.Set("key1", 100)
	s.cache.Clear()
	s.isNotInCache("key1")
}

func cacheSingleItemSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(SingleItemCacheSuite))
}

type MultiItemCacheSuite struct {
	CacheTestHelper
}

func (s *MultiItemCacheSuite) SetupTest() {
	s.cache = NewCache(3)
}

func (s *MultiItemCacheSuite) TestSetToEmpty() {
	s.setNew("key1", 100)
	s.isInCache("key1", 100)
}

func (s *MultiItemCacheSuite) TestSetToPartiallyFilled() {
	s.cache.Set("key1", 100)
	s.setNew("key2", 200)
	s.isInCache("key1", 100)
	s.isInCache("key2", 200)
}

func (s *MultiItemCacheSuite) TestSetWithUpdate() {
	s.cache.Set("key1", 100)
	s.cache.Set("key2", 200)

	s.setExisting("key1", 101)
	s.setExisting("key2", 201)

	s.isInCache("key1", 101)
	s.isInCache("key2", 201)
}

func (s *MultiItemCacheSuite) TestGetFromEmpty() {
	s.isNotInCache("key1")
}

func (s *MultiItemCacheSuite) TestGetFromPartiallyFilled() {
	s.cache.Set("key1", 100)
	s.isInCache("key1", 100)
	s.isNotInCache("key2")
}

func (s *MultiItemCacheSuite) TestGetFromFull() {
	s.cache.Set("key1", 100)
	s.cache.Set("key2", 200)
	s.cache.Set("key3", 300)

	s.isInCache("key1", 100)
	s.isInCache("key2", 200)
	s.isInCache("key3", 300)
}

func (s *MultiItemCacheSuite) TestGetNonExistent() {
	s.cache.Set("key1", 100)
	s.cache.Set("key2", 200)
	s.cache.Set("key3", 300)
	s.isNotInCache("key4")
}

func (s *MultiItemCacheSuite) TestClearEmptyCache() {
	s.cache.Clear()
	s.isNotInCache("key1")
}

func (s *MultiItemCacheSuite) TestClearPartiallyFilled() {
	s.cache.Set("key1", 100)
	s.cache.Set("key2", 200)

	s.cache.Clear()

	s.isNotInCache("key1")
	s.isNotInCache("key2")
}

func (s *MultiItemCacheSuite) TestClearFull() {
	s.cache.Set("key1", 100)
	s.cache.Set("key2", 200)
	s.cache.Set("key3", 300)

	s.cache.Clear()

	s.isNotInCache("key1")
	s.isNotInCache("key2")
	s.isNotInCache("key3")
}

func cacheMultiItemSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(MultiItemCacheSuite))
}

type CacheEvictionSuite struct {
	CacheTestHelper
}

func (s *CacheEvictionSuite) SetupTest() {
	s.cache = NewCache(3)
}

func cacheEvictionSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(CacheEvictionSuite))
}

func (s *CacheEvictionSuite) TestQueueSizeEviction() {
	s.cache.Set("key1", 100)
	s.cache.Set("key2", 200)
	s.cache.Set("key3", 300)

	s.setNew("key4", 400)

	s.isNotInCache("key1")
	s.isInCache("key2", 200)
	s.isInCache("key3", 300)
	s.isInCache("key4", 400)
}

func (s *CacheEvictionSuite) TestUnusedEviction() {
	s.cache.Set("key1", 100) // [100 nil nil]
	s.cache.Set("key2", 200) // [200 100 nil]
	s.cache.Set("key3", 300) // [300 200 100]

	s.cache.Get("key1") // [100 300 200]
	s.cache.Get("key3") // [300 100 200]

	s.cache.Set("key1", 101) // [101 300 200]
	s.cache.Set("key2", 201) // [201 101 300]

	s.cache.Set("key4", 400) // [400 201 101] -> key3 is evicted

	s.isNotInCache("key3")
	s.isInCache("key1", 101)
	s.isInCache("key2", 201)
	s.isInCache("key4", 400)
}

func (s *CacheEvictionSuite) TestSingleItemCacheEviction() {
	c := NewCache(1)
	c.Set("key1", 100)
	c.Set("key2", 200)

	v, ok := c.Get("key1")
	s.False(ok)
	s.Nil(v)

	v, ok = c.Get("key2")
	s.True(ok)
	s.Equal(200, v)
}

type CacheStressSuite struct {
	CacheTestHelper
}

func (s *CacheStressSuite) SetupTest() {
	s.cache = NewCache(1000)
}

func (s *CacheStressSuite) TestHighFrequencySets() {
	iterations := 1_000_000

	for i := 0; i < iterations; i++ {
		key := Key(strconv.Itoa(i))
		s.setNew(key, i)

		// To reduce test execution time.
		if i%1000 == 0 {
			s.isInCache(key, i)
		}
	}
}

func (s *CacheStressSuite) TestHeavyItems() {
	const items = 100_000
	bigValue := make([]byte, 1024) // 1KB

	for i := 0; i < items; i++ {
		key := Key(strconv.Itoa(i))
		s.cache.Set(key, bigValue)

		// To reduce test execution time.
		if i%10_000 == 0 {
			s.isInCache(key, bigValue)
		}
	}

	// Old items should be evicted.
	oldKey := Key("0")
	s.isNotInCache(oldKey)
}

func (s *CacheStressSuite) TestComplexMultithreading() {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	defer wg.Wait()

	getsDone, setsDone, clearsDone := make(chan int), make(chan int), make(chan int)

	go func(res chan<- int) {
		defer wg.Done()
		defer close(res)
		setsDone := 0
		for i := range 100_000 {
			key := Key(strconv.Itoa(i))
			s.cache.Set(key, i)
			setsDone++
		}
		res <- setsDone
	}(getsDone)

	go func(res chan<- int) {
		defer wg.Done()
		defer close(res)
		getsDone := 0
		for i := range 100_000 {
			key := Key(strconv.Itoa(i))
			s.cache.Get(key)
			getsDone++
		}
		res <- getsDone
	}(setsDone)

	go func(res chan<- int) {
		defer wg.Done()
		defer close(res)
		clearsDone := 0
		for range 100_000 {
			s.cache.Clear()
			clearsDone++
		}
		res <- clearsDone
	}(clearsDone)

	total := <-getsDone + <-setsDone + <-clearsDone
	s.Require().Equal(300_000, total)
}

func cacheStressSuite(t *testing.T) {
	t.Helper()
	suite.Run(t, new(CacheStressSuite))
}
