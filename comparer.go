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
			b = delete(b, index(b,v))
		} else {
			compareResult.Added = append(compareResult.Added, v)
		}
	}
	compareResult.Removed = b
	return compareResult
}

func delete(slice []string, index int) []string {
	return append(slice[:index], slice[index+1:]...)
}

func index(slice []string, element string) int {
	for i := range slice {
		if slice[i] == element {
			return i
		}
	}
	return -1
}
func contains(elements []string, item string) bool {
	for _, v := range elements {
		if v == item {
			return true
		}
	}
	return false
}
