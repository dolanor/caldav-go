package utils

func StringPath(paths []string) string {
	for _, test := range paths {
		if test != "" {
			return test
		}
	}
	return "/"
}
