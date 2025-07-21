package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

func DecodeJson(body []byte, result any) error {
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	return nil
}

func GetProjectPath() string {
	projectPath := os.Getenv("PROJECT_PATH")
	if projectPath != "" {
		return projectPath
	}

	currentDir, _ := os.Getwd()
	return currentDir
}

func BuildLookupMap(slice []string) map[string]bool {
	mappedSlice := make(map[string]bool, len(slice))

	for _, key := range slice {
		mappedSlice[key] = true
	}

	return mappedSlice
}

func WriteStdOutErr(logger *zerolog.Logger, err error, msg string) {
	logger.Error().Interface("exchangeInfo", err).Msg(msg)
}

// func pivot(slice []uint, low uint, high uint) uint {
// 	mid := low + (high - low)

// 	if slice[low] > slice[mid] {
// 		swap(&slice[low], &slice[mid])
// 	}
// 	if slice[low] > slice[high] {
// 		swap(&slice[low], &slice[high])
// 	}
// 	if slice[mid] > slice[high] {
// 		swap(&slice[mid], &slice[high])
// 	}

// 	return slice[mid]
// }

// func partition(slice []uint, pivot uint, low uint, high uint) (uint, uint) {
// 	i := low
// 	j := low
// 	k := high

// 	for j <= k {
// 		if slice[j] < pivot {
// 			swap(&slice[i], &slice[j])
// 			i = i + 1
// 			j = j + 1
// 		} else if slice[j] > pivot {
// 			swap(&slice[j], &slice[k])
// 			k = k - 1
// 		} else {
// 			j = j + 1
// 		}
// 	}

// 	return i - 1, k + 1
// }

// func Quicksort(slice []uint, low uint, high uint) {
// 	if low < high {
// 		pivot := pivot(slice, low, high)
// 		left, right := partition(slice, pivot, low, high)
// 		Quicksort(slice, low, left)
// 		Quicksort(slice, right, high)
// 	}
// }

// func swap(a, b *uint) {
// 	*a, *b = *b, *a
// }
