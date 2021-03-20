package database

import "strings"

func strip(s string) string {
	return strings.ReplaceAll(s, " ", "")
}
