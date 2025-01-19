package roundrobin_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cyberhck/roundguard/pkg/roundrobin"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("it returns items in order", func(t *testing.T) {
		rb := roundrobin.New([]int{1, 2, 3, 4})
		value, err := rb.GetNext()
		assert.NoError(t, err)
		assert.Equal(t, 1, *value)
		value, err = rb.GetNext()
		assert.NoError(t, err)
		assert.Equal(t, 2, *value)
		value, err = rb.GetNext()
		assert.NoError(t, err)
		assert.Equal(t, 3, *value)
		value, err = rb.GetNext()
		assert.NoError(t, err)
		assert.Equal(t, 4, *value)
		value, err = rb.GetNext()
		assert.NoError(t, err)
		assert.Equal(t, 1, *value)
	})
	t.Run("race condition, can read while items are being constantly replaced", func(t *testing.T) {
		rb := roundrobin.New([]int{1, 2, 3, 4})
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			for i := 0; i < 100; i++ {
				rb.ResetWithNewItems([]int{5, 6, 7, 8})
				time.Sleep(time.Microsecond * 10)
			}
			wg.Done()
		}()
		go func() {
			val, err := rb.GetNext()
			assert.NoError(t, err)
			assert.NotEmpty(t, 5, *val)
			wg.Done()
		}()
	})
	t.Run("returns an error if no items are available", func(t *testing.T) {
		rb := roundrobin.New([]int{})
		_, err := rb.GetNext()
		assert.EqualError(t, err, "no remaining items")
	})
	t.Run("returns list of all available items", func(t *testing.T) {
		rb := roundrobin.New([]int{1, 2, 3, 4})
		assert.Equal(t, []int{1, 2, 3, 4}, rb.GetAllItems())
	})
}
