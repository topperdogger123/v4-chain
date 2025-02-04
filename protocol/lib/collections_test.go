package lib_test

import (
	"sort"
	"testing"

	"github.com/dydxprotocol/v4-chain/protocol/lib"
	"github.com/dydxprotocol/v4-chain/protocol/x/clob/types"

	"github.com/stretchr/testify/require"
)

func TestContainsDuplicates(t *testing.T) {
	// Empty case.
	require.False(t, lib.ContainsDuplicates([]types.OrderId{}))

	// Unique uint32 case.
	allUniqueUint32s := []uint32{1, 2, 3, 4}
	require.False(t, lib.ContainsDuplicates(allUniqueUint32s))

	// Duplicate uint32 case.
	containsDuplicateUint32 := append(allUniqueUint32s, 3)
	require.True(t, lib.ContainsDuplicates(containsDuplicateUint32))

	// Unique string case.
	allUniqueStrings := []string{"hello", "world", "h", "w"}
	require.False(t, lib.ContainsDuplicates(allUniqueStrings))

	// Duplicate string case.
	containsDuplicateString := append(allUniqueStrings, "world")
	require.True(t, lib.ContainsDuplicates(containsDuplicateString))
}

func TestGetSortedKeys(t *testing.T) {
	tests := map[string]struct {
		inputMap       map[string]string
		expectedResult []string
	}{
		"Nil input": {
			inputMap:       nil,
			expectedResult: []string{},
		},
		"Empty map": {
			inputMap:       map[string]string{},
			expectedResult: []string{},
		},
		"Non-empty map": {
			inputMap: map[string]string{
				"d": "4", "b": "2", "a": "1", "c": "3",
			},
			expectedResult: []string{"a", "b", "c", "d"},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			actualResult := lib.GetSortedKeys[sort.StringSlice](tc.inputMap)
			require.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestMapSlice(t *testing.T) {
	// Can increment all numbers in a slice by 1, and change type to `uint64`.
	require.Equal(
		t,
		[]uint64{2, 3, 4, 5},
		lib.MapSlice(
			[]uint32{1, 2, 3, 4},
			func(a uint32) uint64 {
				return uint64(a + 1)
			},
		),
	)

	// Can return the length of all strings in a slice.
	require.Equal(
		t,
		[]int{1, 2, 3, 5, 0},
		lib.MapSlice(
			[]string{"1", "22", "333", "hello", ""},
			func(a string) int {
				return len(a)
			},
		),
	)

	// Works properly on empty slice.
	require.Equal(
		t,
		[]int{},
		lib.MapSlice(
			[]string{},
			func(a string) int {
				return 1000
			},
		),
	)

	// Works properly on constant function.
	require.Equal(
		t,
		[]bool{true, true, true},
		lib.MapSlice(
			[]string{"hello", "world", "hello"},
			func(a string) bool {
				return true
			},
		),
	)
}

func TestFilterSlice(t *testing.T) {
	// Can filter out all numbers less than 3.
	require.Equal(
		t,
		[]uint32{1, 2},
		lib.FilterSlice(
			[]uint32{1, 2, 3, 4},
			func(a uint32) bool {
				return a < 3
			},
		),
	)

	// Can filter out all strings that have length greater than 3.
	require.Equal(
		t,
		[]string{"hello"},
		lib.FilterSlice(
			[]string{"1", "22", "333", "hello"},
			func(a string) bool {
				return len(a) > 3
			},
		),
	)

	// Works properly on empty slice.
	require.Equal(
		t,
		[]string{},
		lib.FilterSlice(
			[]string{},
			func(a string) bool {
				return true
			},
		),
	)

	// Works properly on constant function that always returns true.
	require.Equal(
		t,
		[]string{"hello", "world", "hello"},
		lib.FilterSlice(
			[]string{"hello", "world", "hello"},
			func(a string) bool {
				return true
			},
		),
	)

	// Works properly on constant function that always returns false.
	require.Equal(
		t,
		[]string{},
		lib.FilterSlice(
			[]string{"hello", "world", "hello"},
			func(a string) bool {
				return false
			},
		),
	)
}

func TestSliceToSet(t *testing.T) {
	slice := make([]int, 0)
	for i := 0; i < 3; i++ {
		slice = append(slice, i)
	}
	set := lib.SliceToSet(slice)
	require.Equal(
		t,
		map[int]struct{}{
			0: {},
			1: {},
			2: {},
		},
		set,
	)
	stringSlice := []string{
		"one",
		"two",
	}
	stringSet := lib.SliceToSet(stringSlice)
	require.Equal(
		t,
		map[string]struct{}{
			"one": {},
			"two": {},
		},
		stringSet,
	)

	emptySlice := []types.OrderId{}
	emptySet := lib.SliceToSet(emptySlice)
	require.Equal(
		t,
		map[types.OrderId]struct{}{},
		emptySet,
	)
}

func TestSliceToSet_PanicOnDuplicate(t *testing.T) {
	stringSlice := []string{
		"one",
		"two",
		"one",
	}
	require.PanicsWithValue(
		t,
		"SliceToSet: duplicate value: one",
		func() {
			lib.SliceToSet(stringSlice)
		},
	)
}

func TestMergeAllMapsWithDistinctKeys(t *testing.T) {
	tests := map[string]struct {
		inputMaps []map[string]string

		expectedMap map[string]string
		expectedErr bool
	}{
		"Success: nil input": {
			inputMaps:   nil,
			expectedMap: map[string]string{},
		},
		"Success: single map": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2",
			},
		},
		"Success: single map, empty": {
			inputMaps:   []map[string]string{},
			expectedMap: map[string]string{},
		},
		"Success: multiple maps, all empty or nil": {
			inputMaps: []map[string]string{
				{}, nil,
			},
			expectedMap: map[string]string{},
		},
		"Success: multiple maps, some empty": {
			inputMaps: []map[string]string{
				{}, nil, {"a": "1", "b": "2"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2",
			},
		},
		"Success: multiple maps, no empty": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
				{"c": "3", "d": "4"},
			},
			expectedMap: map[string]string{
				"a": "1", "b": "2", "c": "3", "d": "4",
			},
		},
		"Error: duplicate keys": {
			inputMaps: []map[string]string{
				{"a": "1", "b": "2"},
				{"c": "3", "d": "4"},
				{"a": "5"}, // duplicate key
			},
			expectedErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.expectedErr {
				require.PanicsWithValue(
					t,
					"duplicate key: a",
					func() { lib.MergeAllMapsMustHaveDistinctKeys(tc.inputMaps...) })
			} else {
				actualMap := lib.MergeAllMapsMustHaveDistinctKeys(tc.inputMaps...)
				require.Equal(t, tc.expectedMap, actualMap)
			}
		})
	}
}
