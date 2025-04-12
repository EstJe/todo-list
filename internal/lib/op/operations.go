package op

import (
	"fmt"
	"runtime"
	"strings"
)

var (
	moduleName string
	cachedPkg  string
)

// Operation возвращает представление операции в формате "CallerName (Package)"
func Operation() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	// Получаем базовое имя функции
	funcName := getBaseFuncName(fn.Name())

	// Получаем базовое имя пакета
	pkg := getCallerPackage(pc)

	return fmt.Sprintf("%s (%s)", funcName, pkg)
}

// getBaseFuncName извлекает читаемое имя функции
func getBaseFuncName(fullName string) string {
	// Примеры входных данных:
	// "github.com/user/project/pkg.(*Type).Method"
	// "github.com/user/project/pkg.Function"

	parts := strings.Split(fullName, ".")
	if len(parts) < 2 {
		return fullName
	}

	lastPart := parts[len(parts)-1]

	// Обработка методов типа (*Type).Method
	if strings.Contains(lastPart, ")") {
		methodParts := strings.Split(lastPart, ".")
		if len(methodParts) == 2 {
			return methodParts[1] + "()"
		}
	}

	return lastPart + "()"
}

// getCallerPackage возвращает имя пакета вызывающей функции
func getCallerPackage(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	fullName := fn.Name()

	// Разбираем полное имя функции
	parts := strings.Split(fullName, ".")
	if len(parts) < 2 {
		return "unknown"
	}

	// Извлекаем часть с пакетом (предпоследний элемент)
	pkgPart := parts[len(parts)-2]

	// Обрабатываем случаи методов:
	// 1. Для обычных функций: "pkg.Func"
	// 2. Для методов: "pkg.(*Type).Method"
	if strings.HasPrefix(pkgPart, "(") && strings.HasSuffix(pkgPart, ")") {
		// Это метод - извлекаем тип
		typePart := strings.TrimPrefix(pkgPart, "(")
		typePart = strings.TrimSuffix(typePart, ")")
		typePart = strings.TrimPrefix(typePart, "*")

		// Ищем пакет типа (предполагаем формат package.Type)
		typeParts := strings.Split(typePart, ".")
		if len(typeParts) > 1 {
			return typeParts[0]
		}
		return typePart
	}

	return pkgPart
}
