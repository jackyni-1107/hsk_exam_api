package utility

import (
	"sort"
	"strconv"
	"strings"
)

// JoinSortedInt64IDs 将 id 排序后拼接，用于审计中多对多关联摘要对比。
func JoinSortedInt64IDs(ids []int64) string {
	if len(ids) == 0 {
		return ""
	}
	s := append([]int64(nil), ids...)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	parts := make([]string, len(s))
	for i, v := range s {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(parts, ",")
}
