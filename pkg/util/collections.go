package util

// https://gobyexample.com/collection-functions
func StringSliceToSet(vs []string) map[string]bool {
	vsm := map[string]bool{}
	for _, v := range vs {
		vsm[v] = true
	}
	return vsm
}
