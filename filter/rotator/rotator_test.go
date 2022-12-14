package rotator

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/x0rworld/go-bloomfilter/bitmap"
	"github.com/x0rworld/go-bloomfilter/config"
	"github.com/x0rworld/go-bloomfilter/filter"
	"testing"
	"time"
)

func TestRotator_doRotation(t *testing.T) {
	rotator := genRotator(t, genDefaultRotatorConfig())
	data := "hello"
	// check filter is empty
	exist, err := rotator.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, false, exist)
	// add element to filter
	err = rotator.Add(data)
	assert.NoError(t, err)
	// check element is in filter
	exist, err = rotator.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, true, exist)

	_ = rotator.rotate()

	// validate data in current and next filter
	cExist, err := rotator.pair.Load().(*filterPair).current.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, true, cExist)
	nExist, err := rotator.pair.Load().(*filterPair).next.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, false, nExist)
}

func TestRotator_Exist(t *testing.T) {
	data := "hello"
	// scenario 1: current & next don't have data, expect to get non-existing
	rotator := genRotator(t, genDefaultRotatorConfig())
	exist, err := rotator.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, false, exist)

	// scenario 2: current does have data but next doesn't, expect to get existing
	rotator = genRotator(t, genDefaultRotatorConfig())
	err = rotator.pair.Load().(*filterPair).current.Add(data)
	assert.NoError(t, err)
	exist, err = rotator.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, true, exist)

	// scenario 3: next does have data but current doesn't, expect to get non-existing (this case should not happen)
	rotator = genRotator(t, genDefaultRotatorConfig())
	err = rotator.pair.Load().(*filterPair).next.Add(data)
	assert.NoError(t, err)
	exist, err = rotator.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, false, exist)

	// scenario 4: current & next do have data, expect to get existing
	rotator = genRotator(t, genDefaultRotatorConfig())
	err = rotator.pair.Load().(*filterPair).current.Add(data)
	assert.NoError(t, err)
	err = rotator.pair.Load().(*filterPair).next.Add(data)
	assert.NoError(t, err)
	exist, err = rotator.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, true, exist)
}

func TestRotator_Add(t *testing.T) {
	rotator := genRotator(t, genDefaultRotatorConfig())
	data := "hello"
	err := rotator.Add(data)
	assert.NoError(t, err)
	cExist, err := rotator.pair.Load().(*filterPair).current.Exist(data)
	assert.NoError(t, err)
	nExist, err := rotator.pair.Load().(*filterPair).next.Exist(data)
	assert.NoError(t, err)
	assert.Equal(t, true, cExist && nExist)
}

func genDefaultRotatorConfig() config.RotatorConfig {
	return config.RotatorConfig{
		Enable: true,
		Freq:   3 * time.Second,
	}
}

func genRotator(t *testing.T, cfg config.RotatorConfig) *Rotator {
	rotator, err := NewRotator(context.Background(), cfg, newFilter)
	if err != nil {
		t.Error("failed to new rotator")
	}
	return rotator
}

func newFilter(_ context.Context) (filter.Filter, error) {
	return filter.NewBloomFilter(bitmap.NewInMemory(100), 100, 3), nil
}
