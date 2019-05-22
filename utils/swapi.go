package utils

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// ResourceId returns resource id from url
// eg: https://swapi.co/api/films/1/
func ResourceId(url string) (int, error) {
	if !strings.Contains(url, "http") {
		return 0, errors.New("Invalid URL!")
	}

	s := strings.Split(url, "/")

	id, err := strconv.Atoi(s[len(s) - 2])
	if err != nil {
		return 0, err
	}

	return id, err
}

// ID returns graphql ID field value for a url
func ID(url string) (string, error) {
	id, err := ResourceId(url)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(id), err
}
