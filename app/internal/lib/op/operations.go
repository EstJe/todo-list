package op

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Operation возвращает полное представление операции в формате "CallerName (PackagePath)"
func Operation() string {
	pc, file, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	callerName := getCallerName(pc, file)
	pkgPath := getPackagePath(file)

	return fmt.Sprintf("%s (%s)", callerName, pkgPath)
}

// GetPackagePath возвращает относительный путь пакета вызывающей функции
func GetPackagePath() string {
	_, file, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return getPackagePath(file)
}

// GetCallerName возвращает имя вызывающей функции/метода
func GetCallerName() string {
	pc, file, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return getCallerName(pc, file)
}

// Общая логика для получения пути пакета
func getPackagePath(file string) string {
	projectRoot, err := findProjectRoot(file)
	if err != nil {
		return ""
	}

	relPath, _ := filepath.Rel(projectRoot, file)
	return filepath.ToSlash(filepath.Dir(relPath))
}

// Общая логика для получения имени вызывающего
func getCallerName(pc uintptr, file string) string {
	projectRoot, err := findProjectRoot(file)
	if err != nil {
		return ""
	}

	moduleName, err := getModuleName(projectRoot)
	if err != nil {
		return ""
	}

	fullFuncName := runtime.FuncForPC(pc).Name()
	return formatCallerName(fullFuncName, moduleName)
}

// Остальные вспомогательные функции остаются без изменений
// (findProjectRoot, getModuleName, formatCallerName)

func findProjectRoot(startPath string) (string, error) {
	currentDir := filepath.Dir(startPath)
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			return "", fmt.Errorf("go.mod not found")
		}
		currentDir = parentDir
	}
}

func getModuleName(projectRoot string) (string, error) {
	data, err := os.ReadFile(filepath.Join(projectRoot, "go.mod"))
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(line[len("module "):]), nil
		}
	}
	return "", fmt.Errorf("module declaration not found")
}

func formatCallerName(fullName, moduleName string) string {
	cleanName := strings.TrimPrefix(fullName, moduleName+"/")
	parts := strings.Split(cleanName, "/")
	if len(parts) == 0 {
		return ""
	}

	lastPart := parts[len(parts)-1]
	subParts := strings.Split(lastPart, ".")

	switch len(subParts) {
	case 2:
		return fmt.Sprintf("%s.%s()", subParts[0], subParts[1])
	case 3:
		typeName := strings.TrimPrefix(subParts[1], "(")
		typeName = strings.TrimSuffix(typeName, ")")
		typeName = strings.TrimPrefix(typeName, "*")
		return fmt.Sprintf("%s.%s{}.%s()", subParts[0], typeName, subParts[2])
	default:
		return lastPart
	}
}
