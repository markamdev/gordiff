package main

import (
	"flag"
	"fmt"
	"githubcom/markamdev/gordiff/pkg/common"
	"githubcom/markamdev/gordiff/pkg/files"
	"githubcom/markamdev/gordiff/pkg/hashing"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	// log everything to stderr as stdout can be used for singature/delta output
	logrus.SetOutput(os.Stderr)

	helpRequest := false
	verboseRun := false

	flag.BoolVar(&helpRequest, "help", false, "print help message")
	flag.BoolVar(&verboseRun, "verbose", false, "launch verbose messages")
	flag.Parse()

	if verboseRun {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	if helpRequest {
		printUsage()
		os.Exit(exitSuccess)
	}

	params := flag.Args()
	if len(params) == 0 {
		logrus.Errorln("no mandatory params given to application")
		printUsage()
		os.Exit(exitInvalidOptions)
	}

	var err error
	switch params[0] {
	case modeDelta:
		logrus.Debug("delta mode selected")
		err = runDelta(params)
	case modeSingature:
		logrus.Debug("signature mode selected")
		err = runSignature(params)
	default:
		printUsage()
		os.Exit(exitInvalidOptions)
	}

	switch err {
	case nil:
		logrus.Debug("exiting without any error")
		os.Exit(exitSuccess)
	case common.ErrInvalidParams:
		logrus.Errorf("funcrion '%s' failed to execute due to invalid params")
		os.Exit(exitInvalidOptions)
	default:
		logrus.Errorf("function '%s' failed to execute with error: %s", params[0], err.Error())
		os.Exit(exitOperationFailure)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `
GoRDIFF - rdiff-like tool written in Golang

Usage:
	%s [OPTIONS] signature [input_file [signature_file]]
	%s [OPTIONS] delta signature_file [updated_file [delta_file]]

Options:
	-verbose 	enable printing logs for each execution step
	-help		print this message

Notes:
	* For all '*_file' parameters the '--' string can be used instead of file paths. GoRDIFF will then use
	standard input or standard output (depending on file role) instead of files.
	* It is not allowed to use '--' for both input files of 'delta' mode
	* It is not allowed to use already existing file as output file (overwriting not possible)

`,
		os.Args[0], os.Args[0])
}

func runDelta(params []string) error {
	if len(params) != 4 {
		printUsage()
		return common.ErrInvalidParams
	}

	if params[1] == "--" && params[2] == "--" {
		return common.ErrInvalidParams
	}

	sigFile, err := files.GetInputFile(params[1])
	if err != nil {
		return fmt.Errorf("failed to open signature file: %w", err)
	}
	defer sigFile.Close()

	newFile, err := files.GetInputFile(params[2])
	if err != nil {
		return fmt.Errorf("failed to open modified file: %w", err)
	}
	defer newFile.Close()

	deltaFile, err := files.GetOutputFile(params[3])
	if err != nil {
		return fmt.Errorf("failed to open file for delta: %w", err)
	}
	defer deltaFile.Close()

	return hashing.GenerateDelta(sigFile, newFile, deltaFile)
}

func runSignature(params []string) error {
	if len(params) != 3 {
		printUsage()
		return common.ErrInvalidParams
	}

	inFile, err := files.GetInputFile(params[1])
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	logrus.Debugf("file data: %+v", *inFile)
	defer inFile.Close()

	sigFile, err := files.GetOutputFile(params[2])
	if err != nil {
		return fmt.Errorf("failed to open file for signature: %w", err)
	}
	logrus.Debugf("file data: %+v", *sigFile)
	defer sigFile.Close()

	return hashing.GenerateSignature(inFile, sigFile)
}
