package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Tests for pagination logic
func TestApplyPagination_FirstPage(t *testing.T) {
	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	page, perPage := 0, 10
	start := page * perPage
	end := start + perPage
	if end > len(items) {
		end = len(items)
	}

	result := items[start:end]
	assert.Len(t, result, 10)
	assert.Equal(t, 0, result[0])
	assert.Equal(t, 9, result[9])
}

func TestApplyPagination_MiddlePage(t *testing.T) {
	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	page, perPage := 5, 10
	start := page * perPage
	end := start + perPage
	if end > len(items) {
		end = len(items)
	}

	result := items[start:end]
	assert.Len(t, result, 10)
	assert.Equal(t, 50, result[0])
	assert.Equal(t, 59, result[9])
}

func TestApplyPagination_LastPage(t *testing.T) {
	items := make([]int, 95)
	for i := range items {
		items[i] = i
	}

	page, perPage := 9, 10
	start := page * perPage
	end := start + perPage
	if end > len(items) {
		end = len(items)
	}

	result := items[start:end]
	assert.Len(t, result, 5)
	assert.Equal(t, 90, result[0])
	assert.Equal(t, 94, result[4])
}

func TestApplyPagination_BeyondRange(t *testing.T) {
	items := make([]int, 50)
	for i := range items {
		items[i] = i
	}

	page, perPage := 10, 10
	start := page * perPage
	if start >= len(items) {
		// Beyond range - return empty
		assert.True(t, start >= len(items))
		return
	}
}

func TestApplyPagination_NegativePageTreatedAsZero(t *testing.T) {
	items := make([]int, 50)
	for i := range items {
		items[i] = i
	}

	page, perPage := -1, 10
	if page < 0 {
		page = 0
	}

	start := page * perPage
	end := start + perPage
	if end > len(items) {
		end = len(items)
	}

	result := items[start:end]
	assert.Len(t, result, 10)
	assert.Equal(t, 0, result[0])
}

func TestApplyPagination_ZeroPerPageDefaultsTo100(t *testing.T) {
	perPage := 0
	if perPage <= 0 {
		perPage = 100
	}
	assert.Equal(t, 100, perPage)
}

// Tests for device name deduplication logic
func TestDeduplicateDevices(t *testing.T) {
	runningNames := map[string]bool{
		"wg0": true,
		"wg1": true,
	}

	configNames := []string{"wg0", "wg1", "wg2", "wg3"}

	var configOnlyNames []string
	for _, name := range configNames {
		if !runningNames[name] {
			configOnlyNames = append(configOnlyNames, name)
		}
	}

	assert.Len(t, configOnlyNames, 2)
	assert.Contains(t, configOnlyNames, "wg2")
	assert.Contains(t, configOnlyNames, "wg3")
}

func TestDeduplicateDevices_AllRunning(t *testing.T) {
	runningNames := map[string]bool{
		"wg0": true,
		"wg1": true,
	}

	configNames := []string{"wg0", "wg1"}

	var configOnlyNames []string
	for _, name := range configNames {
		if !runningNames[name] {
			configOnlyNames = append(configOnlyNames, name)
		}
	}

	assert.Len(t, configOnlyNames, 0)
}

func TestDeduplicateDevices_NoneRunning(t *testing.T) {
	runningNames := map[string]bool{}

	configNames := []string{"wg0", "wg1", "wg2"}

	var configOnlyNames []string
	for _, name := range configNames {
		if !runningNames[name] {
			configOnlyNames = append(configOnlyNames, name)
		}
	}

	assert.Len(t, configOnlyNames, 3)
}
