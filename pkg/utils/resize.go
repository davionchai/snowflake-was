package utils

import (
	"fmt"
	"strings"
)

var availableSize []string = []string{
	"xsmall",
	"small",
	"medium",
	"large",
	"xlarge",
	"xxlarge",
	"xxxlarge",
	"xxxxlarge",
}

type WarehouseCenter struct {
	Size     string
	minSize  string
	maxSize  string
	minIndex int
	maxIndex int
}

func NewWarehouseCenter(size, minSize, maxSize string) (*WarehouseCenter, error) {
	warehouseCenter := &WarehouseCenter{
		Size:     sizeCaster(size),
		minSize:  sizeCaster(minSize),
		maxSize:  sizeCaster(maxSize),
		minIndex: -1,
		maxIndex: -1,
	}

	for idx, val := range availableSize {
		if val == warehouseCenter.minSize {
			warehouseCenter.minIndex = idx
		}
		if val == warehouseCenter.maxSize {
			warehouseCenter.maxIndex = idx
		}
	}

	if warehouseCenter.minIndex == -1 || warehouseCenter.maxIndex == -1 {
		return nil, fmt.Errorf("Invalid warehouse size")
	}

	if warehouseCenter.minIndex > warehouseCenter.maxIndex {
		return nil, fmt.Errorf("Min warehouse size cannot be larger than max warehouse size")
	}

	return warehouseCenter, nil
}

func (warehouseCenter *WarehouseCenter) Upsize() bool {
	index := getIndexOfSize(availableSize, warehouseCenter.Size)
	// if current_size > max_size, go to nearest max_size
	// 	this is to ensure the wh stay max capped size
	if index > warehouseCenter.maxIndex {
		warehouseCenter.Size = availableSize[warehouseCenter.maxIndex]
		return true
	} else if index < warehouseCenter.maxIndex {
		warehouseCenter.Size = availableSize[index+1]
		return true
	}

	return false
}

func (warehouseCenter *WarehouseCenter) Downsize() bool {
	index := getIndexOfSize(availableSize, warehouseCenter.Size)
	// if current_size < min_size, go to nearest min_size
	// 	this is to ensure the wh stay min floored size
	if index < warehouseCenter.minIndex {
		warehouseCenter.Size = availableSize[warehouseCenter.minIndex]
		return true
	} else if index > warehouseCenter.minIndex {
		warehouseCenter.Size = availableSize[index-1]
		return true
	}
	return false
}

func sizeCaster(size string) (castedSize string) {
	size = strings.ToLower(size)
	switch size {
	case "x-small", "xsmall":
		castedSize = "xsmall"
	case "small":
		castedSize = "small"
	case "medium":
		castedSize = "medium"
	case "large":
		castedSize = "large"
	case "x-large", "xlarge":
		castedSize = "xlarge"
	case "2x-large", "xxlarge":
		castedSize = "xxlarge"
	case "3x-large", "xxxlarge":
		castedSize = "xxxlarge"
	case "4x-large", "xxxxlarge":
		castedSize = "xxxxlarge"
	default:
		castedSize = "xsmall"
	}

	return castedSize
}

func getIndexOfSize(availableSize []string, size string) int {
	for idx, val := range availableSize {
		if val == size {
			return idx
		}
	}
	return -1
}
