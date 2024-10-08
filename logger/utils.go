package logger

import (
	"path"
	"runtime"
	"strings"
)

func GetCurrentFuncInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Clean(file)
	return trimFuncName(funcName), trimFileName(fileName), lineNo
}

func trimFuncName(funcName string) string {
	return strings.TrimPrefix(funcName, "github.com/gnanasuryateja/")
}

func trimFileName(fileName string) string {
	fileNamePieces := strings.Split(fileName, "/")
	return fileNamePieces[len(fileNamePieces)-2] + "/" + fileNamePieces[len(fileNamePieces)-1]
}
