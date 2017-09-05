package log

import (
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	typeRollingSize = "size_rolling_file"
	typeRollingTime = "time_rolling_file"
)

type rollingOutput struct {
	Output
}

// NewRollingOutput return a output instance that it print message to stdio
func NewRollingOutput(cfg CfgOutput) (o Output, err error) {
	r := &rollingOutput{}
	o = r

	fpath := cfg["file"]
	rw := newRollingFileWriter(fpath, "")

	rw.archiveType = rollingArchiveNone
	switch cfg["archive"] {
	case "zip":
		rw.archiveType = rollingArchiveZip
	case "gzip":
		rw.archiveType = rollingArchiveGzip
	}

	rw.maxRolls = getMaxRolls(cfg["backups"])
	rw.dirPerm = getFileMode(cfg["dir_perm"], defaultDirectoryPermissions)
	rw.filePerm = getFileMode(cfg["file_perm"], defaultFilePermissions)
	rw.backPerm = getFileMode(cfg["back_perm"], defaultBackupPermissions)

	rw.nameMode = rollingNameModePostfix
	switch cfg["name_mode"] {
	case "prefix":
		rw.nameMode = rollingNameModePrefix
	case "postfix":
		rw.nameMode = rollingNameModePostfix
	}

	var w io.Writer
	if cfg.Name() == typeRollingSize {
		maxSize := getMaxSize(cfg["size"])
		rws := &rollingFileWriterSize{rw, maxSize}
		rws.self = rws
		w = rws
	} else if cfg.Name() == typeRollingTime {
		timePattern := cfg["pattern"]
		rws := &rollingFileWriterTime{rw, timePattern, ""}
		rws.self = rws
		w = rws
	}

	if cfg["async"] == "true" {
		r.Output = NewAsyncOutput(w, GetThresholdLvl(cfg["threshold"]),
			GetQueueSize(cfg["queue_size"]), GetBatchNum(cfg["batch_num"]))
	} else {
		r.Output = NewBaseOutput(w, GetThresholdLvl(cfg["threshold"]))
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
