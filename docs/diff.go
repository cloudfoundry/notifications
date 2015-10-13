package docs

import "regexp"

const (
	date      = "Date: [A-Z]{1}[a-z]{2}, [0-3]{1}[0-9]{1} [A-Z]{1}[a-z]{2} 20[0-9]{2} [0-2]{1}[0-9]{1}:[0-5]{1}[0-9]{1}:[0-5]{1}[0-9]{1} [A-Z]{3}"
	auth      = "Authorization: Bearer .*"
	guid      = "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
	timestamp = "[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z"
)

func Diff(left, right string) bool {
	regexes := []*regexp.Regexp{
		regexp.MustCompile(date),
		regexp.MustCompile(auth),
		regexp.MustCompile(guid),
		regexp.MustCompile(timestamp),
	}

	for _, regex := range regexes {
		left = regex.ReplaceAllString(left, "")
		right = regex.ReplaceAllString(right, "")
	}

	return left != right
}
