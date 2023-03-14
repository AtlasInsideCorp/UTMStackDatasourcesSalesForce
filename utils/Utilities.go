package utils

import "regexp"

func PrepareDate(src string) string {
	regZone := regexp.MustCompile("^(.*)(\\+||-)(\\d{2})(\\d{2})$")
	dest := regZone.ReplaceAllString(src, "$1$2$3:$4")
	return dest
}
