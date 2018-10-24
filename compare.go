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

func FilterComparisonResult(ignoredPackages []string, compareResult CompareResult) CompareResult {
	toRemove := getIntersectingElements(compareResult.Removed, ignoredPackages)
	toAdd := getIntersectingElements(compareResult.Added, ignoredPackages)
	return CompareResult{toAdd, toRemove, compareResult.Unchanged}
}

func getIntersectingElements(a []string, b []string) []string {
	var result []string
	for _, value := range a {
		if contains(b, value) {
			continue
		} else {
			result = append(result, value)
		}
	}
	return result
}
