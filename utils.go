package main

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true iff slice contains element
func contains(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func deleteElementByIndex(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func getElementIndex(slice []string, element string) int {
	for i := range slice {
		if slice[i] == element {
			return i
		}
	}
	return -1
}

func copy(a []string) []string {
	var b []string
	for _, v := range a {
		b = append(b, v)

	}
	return b
}
