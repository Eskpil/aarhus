package contracts

import (
	"fmt"
	"regexp"
	"strings"
)

type Template struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`

	// container tag, for container registry. Eg openjdk:17
	Tag string `json:"tag" bson:"tag"`

	// example: java -Xmx%MEMORY%G -Xmc%CPU% -jar fabricmc-1.20.1.jar nogui
	StartupCommand string `json:"startup_command" bson:"startup_command"`
}

func (t *Template) Parse() []string {
	var keys []string
	var currentKey strings.Builder
	inKey := false

	for i := 0; i < len(t.StartupCommand); i++ {
		if t.StartupCommand[i] == '%' {
			if inKey {
				keys = append(keys, currentKey.String())
				currentKey.Reset()
				inKey = false
			} else {
				inKey = true
			}
		} else if inKey {
			currentKey.WriteByte(t.StartupCommand[i])
		}
	}

	if inKey {
		keys = append(keys, currentKey.String())
	}

	return keys
}

func (t *Template) Apply(arguments map[string]string) (string, error) {
	output := t.StartupCommand

	required := t.Parse()

	for _, value := range required {
		if _, ok := arguments[value]; !ok {
			return "", fmt.Errorf("missing required value: %s", value)
		}
	}

	for key, value := range arguments {
		pattern := regexp.QuoteMeta(fmt.Sprintf("%%%s%%", key))
		re := regexp.MustCompile(pattern)
		output = re.ReplaceAllString(output, value)
	}

	return output, nil
}
