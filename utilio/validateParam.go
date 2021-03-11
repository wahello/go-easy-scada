package utilio

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateConnStr(dst *map[string]*string, srcStr string) error {

	paramSlice := strings.Split(srcStr, ";")
	if len(paramSlice) < len(*dst) {
		return errors.New(fmt.Sprintf("传入的参数不足%d个!", len(*dst)))
	}

	for _, v := range paramSlice {
		tempSlice := strings.Split(v, "=")
		if len(tempSlice) < 2 {
			return errors.New(fmt.Sprintf("非 key=value格式! %s", tempSlice))
		}

		if _, ok := (*dst)[tempSlice[0]]; ok {
			*(*dst)[tempSlice[0]] = tempSlice[1]
		}

	}

	// 空值验证判断
	for k, v := range *dst {
		if *v == "" {
			return errors.New(fmt.Sprintf("mq组件:所需参数[%s]没有传入！", k))
		}
	}

	return nil
}
