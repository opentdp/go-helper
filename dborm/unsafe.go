package dborm

import (
	"errors"
	"regexp"
	"strings"
)

func OrderSafe(data string) error {

	var expr = regexp.MustCompile(`^(\w+)( DESC)?$`)

	for _, v := range strings.Split(data, ",") {
		if !expr.MatchString(v) {
			return errors.New("unsafe sorting field")
		}
	}

	return nil

}
