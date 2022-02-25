package domain

// Contains check uint slice contains element
func Contains(s []uint, x uint) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}

	return false
}

// Remove delete uint value from slice. It returns (newSlice, removedID)
// 0 means invalid value
func Remove(s []uint, val uint) ([]uint, uint) {
	for i, v := range s {
		if v == val {
			return append(s[:i], s[i+1:]...), v
		}
	}
	return s, 0
}
