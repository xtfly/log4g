package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	messageLen = 10
)

var bytesFileTest = []byte(strings.Repeat("A", messageLen))

//===============================================================
type fileWriterTestCase struct {
	files           []string
	fileName        string
	rollingType     rollingType
	fileSize        int64
	maxRolls        int
	datePattern     string
	writeCount      int
	resFiles        []string
	nameMode        rollingNameMode
	archiveType     rollingArchiveType
	archiveExploded bool
	archivePath     string
}

func createRollingSizeFileWriterTestCase(
	files []string,
	fileName string,
	fileSize int64,
	maxRolls int,
	writeCount int,
	resFiles []string,
	nameMode rollingNameMode,
	archiveType rollingArchiveType,
	archiveExploded bool,
	archivePath string) *fileWriterTestCase {

	return &fileWriterTestCase{files, fileName, rollingTypeSize, fileSize, maxRolls, "", writeCount, resFiles, nameMode, archiveType, archiveExploded, archivePath}
}

type fileWriterTester struct {
	testCases    []*fileWriterTestCase
	writerGetter func(*fileWriterTestCase) (io.WriteCloser, error)
	t            *testing.T
}

func (tester *fileWriterTester) testCase(testCase *fileWriterTestCase, testNum int) {
	defer cleanupWriterTest(tester.t)

	tester.t.Logf("Start test  [%v]\n", testNum)

	for _, filePath := range testCase.files {
		dir, _ := filepath.Split(filePath)

		var err error

		if 0 != len(dir) {
			err = os.MkdirAll(dir, defaultDirectoryPermissions)
			if err != nil {
				tester.t.Error(err)
				return
			}
		}

		fi, err := os.Create(filePath)
		if err != nil {
			tester.t.Error(err)
			return
		}

		err = fi.Close()
		if err != nil {
			tester.t.Error(err)
			return
		}
	}

	fwc, err := tester.writerGetter(testCase)
	if err != nil {
		tester.t.Error(err)
		return
	}
	defer fwc.Close()

	tester.performWrite(fwc, testCase.writeCount)

	files, err := getWriterTestResultFiles()
	if err != nil {
		tester.t.Error(err)
		return
	}

	tester.checkRequiredFilesExist(testCase, files)
	tester.checkJustRequiredFilesExist(testCase, files)

}

func (tester *fileWriterTester) test() {
	for i, tc := range tester.testCases {
		cleanupWriterTest(tester.t)
		tester.testCase(tc, i)
	}
}

func (tester *fileWriterTester) performWrite(fileWriter io.Writer, count int) {
	for i := 0; i < count; i++ {
		_, err := fileWriter.Write(bytesFileTest)

		if err != nil {
			tester.t.Error(err)
			return
		}
	}
}

func (tester *fileWriterTester) checkRequiredFilesExist(testCase *fileWriterTestCase, files []string) {
	var found bool
	for _, expected := range testCase.resFiles {
		found = false
		exAbs, err := filepath.Abs(expected)
		if err != nil {
			tester.t.Errorf("filepath.Abs failed for %s", expected)
			continue
		}

		for _, f := range files {
			if af, e := filepath.Abs(f); e == nil {
				tester.t.Log(af)
				if exAbs == af {
					found = true
					break
				}
			} else {
				tester.t.Errorf("filepath.Abs failed for %s", f)
			}
		}

		if !found {
			tester.t.Errorf("expected file: %s doesn't exist. Got %v\n", exAbs, files)
		}
	}
}

func (tester *fileWriterTester) checkJustRequiredFilesExist(testCase *fileWriterTestCase, files []string) {
	for _, f := range files {
		found := false
		for _, expected := range testCase.resFiles {

			exAbs, err := filepath.Abs(expected)
			if err != nil {
				tester.t.Errorf("filepath.Abs failed for %s", expected)
			} else {
				if exAbs == f {
					found = true
					break
				}
			}
		}

		if !found {
			tester.t.Errorf("unexpected file: %v", f)
		}
	}
}

func getWriterTestResultFiles() ([]string, error) {
	var p []string

	visit := func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() && isWriterTestFile(path) {
			abs, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("filepath.Abs failed for %s", path)
			}

			p = append(p, abs)
		}

		return nil
	}

	err := filepath.Walk(".", visit)
	if nil != err {
		return nil, err
	}

	return p, nil
}

func isWriterTestFile(fn string) bool {
	return strings.Contains(fn, ".testlog") || strings.Contains(fn, ".zip") || strings.Contains(fn, ".gz")
}

func cleanupWriterTest(t *testing.T) {
	toDel, err := getDirFilePaths(".", isWriterTestFile, true)
	if nil != err {
		t.Fatal("Cannot list files in test directory!")
	}

	for _, p := range toDel {
		if err = tryRemoveFile(p); nil != err {
			t.Errorf("cannot remove file %s in test directory: %s", p, err.Error())
		}
	}

	if err = os.RemoveAll("dir"); nil != err {
		t.Errorf("cannot remove temp test directory: %s", err.Error())
	}
}

//===============================================================

func newFileWriterTester(
	testCases []*fileWriterTestCase,
	writerGetter func(*fileWriterTestCase) (io.WriteCloser, error), t *testing.T) *fileWriterTester {

	return &fileWriterTester{testCases, writerGetter, t}
}

func createRollingDatefileWriterTestCase(
	files []string,
	fileName string,
	datePattern string,
	writeCount int,
	resFiles []string,
	nameMode rollingNameMode,
	archiveType rollingArchiveType,
	archiveExploded bool,
	archivePath string) *fileWriterTestCase {

	return &fileWriterTestCase{files, fileName, rollingTypeTime, 0, 0, datePattern, writeCount, resFiles, nameMode, archiveType, archiveExploded, archivePath}
}

func TestShouldArchiveWithTar(t *testing.T) {
	compressionType := compressionTypes[rollingArchiveGzip]

	archiveName := compressionType.rollingArchiveTypeName("log", false)

	if archiveName != "log.tar.gz" {
		t.Fatalf("archive name should be log.tar.gz but got %v", archiveName)
	}
}

func TestRollingFileWriter(t *testing.T) {
	t.Logf("Starting rolling file writer tests")
	newFileWriterTester(rollingfileWriterTests, rollingFileWriterGetter, t).test()
}

//===============================================================

func rollingFileWriterGetter(testCase *fileWriterTestCase) (io.WriteCloser, error) {
	rw := newRollingFileWriter(testCase.fileName, testCase.archivePath)
	rw.archiveType = testCase.archiveType
	rw.maxRolls = testCase.maxRolls
	rw.nameMode = testCase.nameMode
	rw.archiveExploded = testCase.archiveExploded

	if testCase.rollingType == rollingTypeSize {
		rws := &rollingFileWriterSize{rw, testCase.fileSize}
		rws.self = rws
		return rws, nil
	} else if testCase.rollingType == rollingTypeTime {
		rws := &rollingFileWriterTime{rw, testCase.datePattern, ""}
		rws.self = rws
		return rws, nil
	}

	return nil, fmt.Errorf("incorrect rollingType")
}

//===============================================================
var rollingfileWriterTests = []*fileWriterTestCase{
	createRollingSizeFileWriterTestCase([]string{}, "log.testlog", 10, 10, 1, []string{"log.testlog"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, "log.testlog", 10, 10, 2, []string{"log.testlog", "log.testlog.1"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{"1.log.testlog"}, "log.testlog", 10, 10, 2, []string{"log.testlog", "1.log.testlog", "2.log.testlog"}, rollingNameModePrefix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{"log.testlog.1"}, "log.testlog", 10, 1, 2, []string{"log.testlog", "log.testlog.2"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, "log.testlog", 10, 1, 2, []string{"log.testlog", "log.testlog.1"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{"log.testlog.9"}, "log.testlog", 10, 1, 2, []string{"log.testlog", "log.testlog.10"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{"log.testlog.a", "log.testlog.1b"}, "log.testlog", 10, 1, 2, []string{"log.testlog", "log.testlog.1", "log.testlog.a", "log.testlog.1b"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, `dir/log.testlog`, 10, 10, 1, []string{`dir/log.testlog`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, `dir/log.testlog`, 10, 10, 2, []string{`dir/log.testlog`, `dir/1.log.testlog`}, rollingNameModePrefix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{`dir/dir/log.testlog.1`}, `dir/dir/log.testlog`, 10, 10, 2, []string{`dir/dir/log.testlog`, `dir/dir/log.testlog.1`, `dir/dir/log.testlog.2`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{`dir/dir/dir/log.testlog.1`}, `dir/dir/dir/log.testlog`, 10, 1, 2, []string{`dir/dir/dir/log.testlog`, `dir/dir/dir/log.testlog.2`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, `./log.testlog`, 10, 1, 2, []string{`log.testlog`, `log.testlog.1`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{`././././log.testlog.9`}, `log.testlog`, 10, 1, 2, []string{`log.testlog`, `log.testlog.10`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{"dir/dir/log.testlog.a", "dir/dir/log.testlog.1b"}, "dir/dir/log.testlog", 10, 1, 2, []string{"dir/dir/log.testlog", "dir/dir/log.testlog.1", "dir/dir/log.testlog.a", "dir/dir/log.testlog.1b"}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, `././dir/log.testlog`, 10, 10, 1, []string{`dir/log.testlog`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, `././dir/log.testlog`, 10, 10, 2, []string{`dir/log.testlog`, `dir/log.testlog.1`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{`././dir/dir/log.testlog.1`}, `dir/dir/log.testlog`, 10, 10, 2, []string{`dir/dir/log.testlog`, `dir/dir/log.testlog.1`, `dir/dir/log.testlog.2`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{`././dir/dir/dir/log.testlog.1`}, `dir/dir/dir/log.testlog`, 10, 1, 2, []string{`dir/dir/dir/log.testlog`, `dir/dir/dir/log.testlog.2`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{}, `././log.testlog`, 10, 1, 2, []string{`log.testlog`, `log.testlog.1`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{`././././log.testlog.9`}, `log.testlog`, 10, 1, 2, []string{`log.testlog`, `log.testlog.10`}, rollingNameModePostfix, rollingArchiveNone, false, ""),
	createRollingSizeFileWriterTestCase([]string{"././dir/dir/log.testlog.a", "././dir/dir/log.testlog.1b"}, "dir/dir/log.testlog", 10, 1, 2, []string{"dir/dir/log.testlog", "dir/dir/log.testlog.1", "dir/dir/log.testlog.a", "dir/dir/log.testlog.1b"}, rollingNameModePostfix, rollingArchiveNone, true, ""),
	//createRollingSizeFileWriterTestCase([]string{"log.testlog", "log.testlog.1"}, "log.testlog", 10, 1, 2, []string{"log.testlog", "log.testlog.2", "dir/log.testlog.1.zip"}, rollingNameModePostfix, rollingArchiveZip, true, "dir"),
	// ====================
}
