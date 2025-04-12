package op

import (
	"fmt"
	"runtime"
)

// Wrap оборачивает ошибку с информацией об операции
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	pc, _, _, ok := runtime.Caller(2) // Пропускаем Wrap и вызывающую функцию
	if !ok {
		return fmt.Errorf("unknown operation: %w", err)
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return fmt.Errorf("unknown: %w", err)
	}

	funcName := getBaseFuncName(fn.Name())
	pkg := getCallerPackage(pc)

	return fmt.Errorf("%s (%s): \n%w", funcName, pkg, err)
}
