package main

type CompareResult struct {
	Added     []string
	Removed   []string
	Unchanged []string
}

func ComparePackages(a []string, b []string) CompareResult {
	var compareResult = CompareResult{}
	for _, v := range a {
		res := contains(b, v)
		if res {
			compareResult.Unchanged = append(compareResult.Unchanged, v)
			b = deleteElementByIndex(b, getElementIndex(b,v))
		} else {
			compareResult.Added = append(compareResult.Added, v)
		}
	}
	compareResult.Removed = b
	return compareResult
}


