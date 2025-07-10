package utils

import (
	"fmt"
	"slices"
	"store/constant"
	"strings"
)

func SortChecker(definedSort []string, sort string) (res string, err error) {
	var order = "ASC"

	if desc := strings.HasPrefix(sort, "-"); desc {
		order = "DESC"
		sort = strings.TrimPrefix(sort, "-")
	}

	if !slices.Contains(constant.ProductSort, sort) {
		return "", fmt.Errorf("invalid sort: %s", sort)
	}

	return sort + " " + order, nil
}
