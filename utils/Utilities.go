package utils

import (
	"os"
	"regexp"
	"strings"
)

// PrepareDate is a method to modify SalesForce date string to a go time string
func PrepareDate(src string) string {
	regZone := regexp.MustCompile("^(.*)(\\+||-)(\\d{2})(\\d{2})$")
	dest := regZone.ReplaceAllString(src, "$1$2$3:$4")
	return dest
}

// GetEnvOrElse is a method to get the value of environment variables and set a default if empty or not present
func GetEnvOrElse(envVar string, defaultValue string) string {
	str := os.Getenv(envVar)
	if strings.Compare(str, "") == 0 {
		str = defaultValue
	}
	return str
}
