// Package log contains helper for logrus
package log

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// WithFields creates a logrus.Entry with specified fields as pairs (key, value, key, value, ...)
// panics len of fields is not even
func WithFields(fields ...interface{}) *logrus.Entry {
	return logrus.WithFields(pairs(fields...))
}

// Debug prints debug msgs
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

func pairs(kv ...interface{}) map[string]interface{} {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("odd number of pairs?!: %d", len(kv)))
	}

	v := map[string]interface{}{}
	var key string
	for i, s := range kv {
		if i%2 == 0 {
			key = fmt.Sprint(s)
			continue
		}

		v[key] = s
	}
	return v
}
