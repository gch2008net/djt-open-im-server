package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/magefile/mage/sh"
	"github.com/openimsdk/gomake/mageutil"
)

func main() {
	fmt.Println("mian.................11111111111111")
	mageutil.InitForSSC()
	BuildAndStartAndStopSingleCmd()
	BuildLinux()
	// mageutil.Build()
	// mageutil.StartToolsAndServices()
	// mageutil.StopAndCheckBinaries()
	// mageutil.CheckAndReportBinariesStatus()

	//单独build/start/stop某个模块
	// BuildAndStartAndStopSingleCmd()

	// stop
	// path := "D:\\OpenIM\\open-im-server\\_output\\bin\\platforms\\windows\\amd64\\openim-api.exe"
	// path := "D:\\OpenIM\\open-im-server\\_output_bin\\openim-api.exe~"
	// mageutil.BatchKillExistBinaries([]string{path})

	// if err := cmd.NewApiCmd().Exec(); err != nil {
	// 	program.ExitWithError(err)
	// }
}

var Default = Build

func Build() {
	mageutil.Build()
}

func Start() {
	mageutil.InitForSSC()
	// err := setMaxOpenFiles()
	// if err != nil {
	// 	mageutil.PrintRed("setMaxOpenFiles failed " + err.Error())
	// 	os.Exit(1)
	// }
	mageutil.StartToolsAndServices()
}

func Stop() {
	mageutil.StopAndCheckBinaries()
}

func Check() {
	mageutil.CheckAndReportBinariesStatus()
}

func BuildLinux() {
	if _, err := os.Stat("start-config.yml"); err == nil {
		mageutil.InitForSSC()
		mageutil.KillExistBinaries()
	}

	//制定系统架构进行build
	platforms := "linux_amd64"
	/*  */ for _, platform := range strings.Split(platforms, " ") {
		mageutil.CompileForPlatform(platform)
	}
	mageutil.PrintGreen("All binaries under cmd and tools were successfully compiled.")
}

func BuildSingleCmd(dirFix string, filename string, target string) (string, string, error) {

	mageutil.PrintGreen("start BuildSingleCmd ...")

	//env
	targetOS, targetArch := runtime.GOOS, runtime.GOARCH
	_output_bin := "_output_bin/windows"
	if target == "linux" {
		targetOS = "linux"
		targetArch = "amd64"
		_output_bin = "_output_bin/linux"
	}

	//outdir
	//build后的输出路径
	outputDir := filepath.Join(dirFix, _output_bin)

	//create dir
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Println("failed to  create directory .")
		return "", "", err
	}

	outputFileName := filename
	if targetOS == "windows" {
		outputFileName += ".exe"
	}

	dir := filepath.Join(dirFix, "cmd", filename)

	err := sh.RunWith(map[string]string{"GOOS": targetOS, "GOARCH": targetArch}, "go", "build", "-o", filepath.Join(outputDir, outputFileName), filepath.Join(dir, "main.go"))
	if err != nil {
		fmt.Println("err:" + err.Error())
		return "", "", err
	}

	fmt.Printf("%v successfully compiled. \n", filename)

	return outputDir, outputFileName, nil
}

func StartSingleCmd(outputDir string, outputFileName string) {
	mageutil.PrintGreen("start StartSingleCmd ...")

	serviceBinaries := make(map[string]int)
	serviceBinaries[outputFileName] = 1
	for binary, count := range serviceBinaries {
		binFullPath := filepath.Join(outputDir, binary)
		for i := 0; i < count; i++ {
			args := []string{"-i", strconv.Itoa(i), "-c", mageutil.OpenIMOutputConfig}
			cmd := exec.Command(binFullPath, args...)
			fmt.Printf("Starting %s\n", cmd.String())
			cmd.Dir = outputDir
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				fmt.Errorf("failed to start %s with args %v: %v", binFullPath, args, err)
			}
		}
	}
}

func BuildAndStartAndStopSingleCmd() {

	mageutil.PrintGreen("start BuildAndStartAndStopSingleCmd ...")

	mageutil.InitForSSC()

	//项目路径
	var dirFix string = "D:\\OpenIM\\v3.8.0\\open-im-server"

	//模块名称
	// var filename string = "openim-api"
	// var filename string =  "openim-cmdutils"
	// var filename string =  "openim-crontask"
	// var filename string =  "openim-msggateway"
	// var filename string =  "openim-msgtransfer"
	var filename string = "openim-rpc\\openim-rpc-friend"
	// var filename string =  "openim-rpc"

	//单独build某个模块
	outputDir, outputFileName, err := BuildSingleCmd(dirFix, filename, "linux")
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}

	//单独start某个模块
	StartSingleCmd(outputDir, outputFileName)

	//单独stop某个模块
	mageutil.PrintGreen("start BatchKillExistBinaries ...")

	path := filepath.Join(outputDir, outputFileName)
	var binaryPaths []string = []string{path}
	mageutil.BatchKillExistBinaries(binaryPaths)

}
