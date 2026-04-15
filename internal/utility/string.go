package utility

// MaskString Token (只保留前6位和后4位，中间隐藏)
func MaskString(s string, prefixLen, suffixLen int) string {
	if len(s) <= prefixLen+suffixLen {
		return s
	}
	return s[:prefixLen] + "****" + s[len(s)-suffixLen:]
}
