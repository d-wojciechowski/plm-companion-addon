package util

import (
	"dominikw.pl/wnc_plugin/proto/files"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	regexp2 "regexp"
	"runtime"
	"sort"
	"strconv"
	"time"

	"github.com/google/logger"
)

func FindLogFile(directory *files.LogFileLocation) (logFileName string, e error) {
	logger.Infof("LogViewer request received, looking for log files in: %s", directory.FileLocation)
	if directory.LogType == files.LogFileLocation_CUSTOM {
		_, e := ioutil.ReadFile(directory.FileLocation)
		return directory.FileLocation, e
	}
	infos, e := ioutil.ReadDir(directory.FileLocation)
	if e != nil {
		return "", e
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime().After(infos[j].ModTime())
	})
	logFileName, e = findLogFileOfTypeInDir(infos, directory.LogType)
	if e != nil {
		return "", e
	}
	return
}

/*
GetWindowsDrives returns all windows drives mounted in system.
If running system is not windows, empty array is returned.
*/
func GetWindowsDrives() (r []string) {
	if runtime.GOOS != "windows" {
		return []string{}
	}
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		f, err := os.Open(string(drive) + ":\\")
		if err == nil {
			r = append(r, string(drive))
			f.Close()
		}
	}
	return
}

func findLogFileOfTypeInDir(infos []os.FileInfo, logTypeEnum files.LogFileLocation_Source) (string, error) {
	logFileName := ""
	for _, info := range infos {
		if checkFileName(info.Name(), logTypeEnum) {
			logFileName = info.Name()
			logger.Infof("Log files chosen: %s", logFileName)
			break
		}
		if checkFileNameOmittingDate(info.Name(), logTypeEnum) {
			logFileName = info.Name()
			logger.Infof("Log files chosen: %s", logFileName)
			break
		}
	}
	if logFileName == "" {
		logger.Errorf("Log files not found for type: %s!", logTypeEnum.String())
		return "", errors.New(fmt.Sprintf("Log files not found for type: %s!", logTypeEnum.String()))
	}
	return logFileName, nil
}

func checkFileName(fileName string, logTypeEnum files.LogFileLocation_Source) (matched bool) {
	now := time.Now()
	t := getTypeString(logTypeEnum)

	regexp := fmt.Sprintf(`(?m)^%s\d{0,2}-%s%s%s\d{1,99}-\d{1,99}-log4j.log`,
		t,
		strconv.Itoa(now.Year())[:2],
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()))
	matched, _ = regexp2.MatchString(regexp, fileName)
	return
}

func checkFileNameOmittingDate(fileName string, logTypeEnum files.LogFileLocation_Source) (matched bool) {
	t := getTypeString(logTypeEnum)
	matched, _ = regexp2.MatchString(fmt.Sprintf(`(?m)^%s\d{0,2}-\d{1,99}-\d{1,99}-log4j.log`, t), fileName)
	return
}

func getTypeString(logType files.LogFileLocation_Source) (t string) {
	switch logType {
	case files.LogFileLocation_METHOD_SERVER:
		t = "MethodServer"
	case files.LogFileLocation_BACKGROUND_METHOD_SERVER:
		t = "BackgroundMethodServer"
	default:
		t = "MethodServer"
	}
	return
}

func GetPath(logFileDirectory string, logFileName string, logFile *files.LogFileLocation) (path string) {
	path = filepath.Join(logFileDirectory, logFileName)
	if logFile.LogType == files.LogFileLocation_CUSTOM {
		path = logFileName
	}
	return
}
