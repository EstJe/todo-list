package op

import (
	"fmt"
	"runtime"
)

func Wrap(err error) error {
	if err == nil {
		return nil
	}

	// Получаем информацию о вызывающем коде (пропускаем 2 кадра: сам Wrap и вызывающую его функцию)
	pc, file, _, ok := runtime.Caller(2)
	if !ok {
		return fmt.Errorf("unknown operation: %w", err)
	}

	op := formatOperation(pc, file)
	return fmt.Errorf("%s: %w", op, err)
}

// Вынесенная функция для форматирования операции
func formatOperation(pc uintptr, file string) string {
	projectRoot, err := findProjectRoot(file)
	if err != nil {
		return "unknown"
	}

	moduleName, err := getModuleName(projectRoot)
	if err != nil {
		return "unknown"
	}

	fullFuncName := runtime.FuncForPC(pc).Name()
	callerName := formatCallerName(fullFuncName, moduleName)
	pkgPath := getPackagePath(file)

	return fmt.Sprintf("%s (%s)", callerName, pkgPath)
}
