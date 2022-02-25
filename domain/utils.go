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
