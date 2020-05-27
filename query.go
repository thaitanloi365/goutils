package goutils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseQueryToArrayString parse array
func ParseQueryToArrayString(value string, delimiter ...string) []string {
	if value == "" {
		return []string{}
	}
	var c = ","
	if len(delimiter) > 0 {
		c = delimiter[0]
	}

	var segments = strings.Split(value, c)
	var listStatus = []string{}

	for _, s := range segments {
		listStatus = append(listStatus, s)
	}

	return listStatus
}

// ParseQueryToArrayInt64 parse array
func ParseQueryToArrayInt64(value string, delimiter ...string) []int64 {
	if value == "" {
		return []int64{}
	}
	var c = ","
	if len(delimiter) > 0 {
		c = delimiter[0]
	}

	var segments = strings.Split(value, c)
	var listStatus = []int64{}

	for _, s := range segments {
		v, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			listStatus = append(listStatus, v)

		}
	}

	return listStatus
}

// ParseQueryToArrayInt parse array
func ParseQueryToArrayInt(value string, delimiter ...string) []int {
	if value == "" {
		return []int{}
	}
	var c = ","
	if len(delimiter) > 0 {
		c = delimiter[0]
	}

	var segments = strings.Split(value, c)
	var listStatus = []int{}

	for _, s := range segments {
		v, err := strconv.ParseInt(s, 10, 64)
		if err == nil {
			listStatus = append(listStatus, int(v))

		}
	}

	return listStatus
}

// ParseQueryToOrderBy get order by
func ParseQueryToOrderBy(query string, keys ...string) []string {
	var orderBy = []string{}
	if query != "" {
		var segments = strings.Split(query, "|")
		for _, segment := range segments {
			var keyOps = strings.Split(segment, ",")
			if len(keyOps) == 2 {
				var key = keyOps[0]
				var op = keyOps[1]
				if len(keys) > 0 {
					if ContainStr(keys, key) {
						orderBy = append(orderBy, fmt.Sprintf("%s %s", key, op))
					}
				} else {
					orderBy = append(orderBy, fmt.Sprintf("%s %s", key, op))
				}
			}
		}
	}

	return orderBy
}

// ContainStr contains
func ContainStr(list []string, item string) bool {
	for _, i := range list {
		return item == i
	}

	return false
}

// ContainInt contains
func ContainInt(list []int, item int) bool {
	for _, i := range list {
		return item == i
	}

	return false
}
