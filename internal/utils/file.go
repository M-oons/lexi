package utils

import "os"

func IsFile(input string) bool {
	info, err := os.Stat(input)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}
