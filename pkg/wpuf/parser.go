package wpuf

import (
	"regexp"
	"strings"
)

// parse the response received when looking up the username from id
func parseAuthorResponse(bodyString string) string {
	re, _ := regexp.Compile("title.*&#8211")
	match := re.FindString(bodyString)
	// first try
	if len(match) != 0 {
		user := strings.ReplaceAll(match, "title>", "")
		user = strings.ReplaceAll(user, " &#8211", "")
		if len(user) < 500 {
			return user
		}
	}

	re, _ = regexp.Compile("/author/.*/")
	match = re.FindString(bodyString)
	// first try
	if len(match) != 0 {
		user := strings.Split(match, "author/")[1]
		user = strings.Split(user, "/")[0]
		return user
	}
	return ""
}
