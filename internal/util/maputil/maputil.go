package maputil

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/samber/lo"

	"github.com/hanselrd/domino/internal/util/sliceutil"
)

func SortedKeys[K comparable, V any](in ...map[K]V) ([]K, error) {
	ks := lo.Keys(in...)
	if ks, ok := sliceutil.Convert[K, string](ks); ok {
		slices.Sort(ks)
		ks, _ := sliceutil.Convert[string, K](ks)
		return ks, nil
	}
	return nil, fmt.Errorf("failed to sort %s keys: %v", reflect.TypeOf(ks[0]).Name(), ks)
}
