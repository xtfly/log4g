package log

import (
	"io"
	"os"
	"strconv"
	"strings"
)

type rollingOutput struct {
	Output
}

// NewRollingOutput return a output instance that it print message to stdio
func NewRollingOutput(arg *CfgOutput) (o Output, err error) {
	r := &rollingOutput{}
	o = r

	fpath := arg.Properties["file"]
	rw := newRollingFileWriter(fpath, "")

	rw.archiveType = rollingArchiveNone
	switch arg.Properties["archive"] {
	case "zip":
		rw.archiveType = rollingArchiveZip
	case "gzip":
		rw.archiveType = rollingArchiveGzip
	}

	rw.maxRolls = getMaxRolls(arg.Properties["backups"])
	rw.dirPerm = getFileMode(arg.Properties["dir_perm"], defaultDirectoryPermissions)
	rw.filePerm = getFileMode(arg.Properties["file_perm"], defaultFilePermissions)
	rw.backPerm = getFileMode(arg.Properties["back_perm"], defaultBackupPermissions)

	rw.nameMode = rollingNameModePostfix
	switch arg.Properties["name_mode"] {
	case "prefix":
		rw.nameMode = rollingNameModePrefix
	case "postfix":
		rw.nameMode = rollingNameModePostfix
	}

	var w io.Writer
	if arg.Name == "size_rolling_file" {
		maxSize := getMaxSize(arg.Properties["size"])
		rws := &rollingFileWriterSize{rw, maxSize}
		rws.self = rws
		w = rws
	} else if arg.Name == "time_rolling_file" {
		timePattern := arg.Properties["pattern"]
		rws := &rollingFileWriterTime{rw, timePattern, ""}
		rws.self = rws
		w = rws
	}

	if arg.Properties["async"] == "true" {
		r.Output = NewAynscOutput(w,
			GetQueueSize(arg.Properties["queue_size"]), GetBatchNum(arg.Properties["batch_num"]))
	} else {
		r.Output = NewBaseOutput(w)
	}
	return r, nil
}

func getMaxSize(str string) int64 {
	if strings.HasSuffix(str, "K") {
		size, _ := strconv.Atoi(str)
		if size <= 0 {
			size = 10 * 1024
		}
		return int64(size * 1024)
	} else if strings.HasSuffix(str, "M") {
		size, _ := strconv.Atoi(str)
		if size <= 0 {
			size = 10
		}
		return int64(size * 1024 * 1024)
	} else if strings.HasSuffix(str, "G") {
		size, _ := strconv.Atoi(str)
		if size <= 0 {
			size = 1
		}
		return int64(size * 1024 * 1024 * 1024)
	}

	size, _ := strconv.Atoi(str)
	if size <= 0 {
		size = 10 * 1024 * 1024
	}
	return int64(size)
}

func getMaxRolls(str string) int {
	mr, _ := strconv.Atoi(str)
	if mr <= 0 {
		mr = 5
	}
	if mr > 20 {
		mr = 20
	}
	return mr
}

func getFileMode(str string, mode os.FileMode) os.FileMode {
	fm, _ := strconv.ParseInt(str, 0, 32)
	ret := os.FileMode(fm)
	if fm <= 0 {
		ret = mode
	}
	if fm > 0777 {
		ret = mode
	}
	return ret
}
