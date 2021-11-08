package base

import (
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
)

func Check(a interface{}) {
	if a == nil {
		_, f, l, _ := runtime.Caller(1)
		strs := strings.Split(f, "/")
		size := len(strs)
		if size > 4 {
			size = 4
		}
		f = filepath.Join(strs[len(strs)-size:]...)
		log.Panicf("object can't be nil, cause by: %s(%d)", f, l)
	}
}
