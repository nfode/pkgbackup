package main

type CompareResult struct {
	Added     []string
	Removed   []string
	Unchanged []string
}

func ComparePackages(a []string, b []string) CompareResult {
	first := copy(a)
	second := copy(b)
	var compareResult = CompareResult{}
	for _, v := range first {
		res := contains(second, v)
		if res {
			compareResult.Unchanged = append(compareResult.Unchanged, v)
			second = deleteElementByIndex(second, getElementIndex(second, v))
		} else {
			compareResult.Added = append(compareResult.Added, v)
		}
	}
	compareResult.Removed = second
	return compareResult
}
