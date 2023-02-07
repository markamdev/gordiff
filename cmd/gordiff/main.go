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
	forcedRun := false

	flag.BoolVar(&helpRequest, "help", false, "print help message")
	flag.BoolVar(&verboseRun, "verbose", false, "launch verbose messages")
	flag.BoolVar(&forcedRun, "force", false, "overwrite output file if exists")
	flag.Parse()

	if verboseRun {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.ErrorLevel)
	}

	if helpRequest {
		printUsage()
		// save weak sum
		os.Exit(exitSuccess)
	}

	params := flag.Args()
	if len(params) == 0 {
		logrus.Errorln("no mandatory params given to application")
		printUsage()
		os.Exit(exitInvalidOptions)
	}

	opts := &appOptions{forceMode: forcedRun}

	var err error
	switch params[0] {
	case modeDelta:
		logrus.Debug("delta mode selected")
		err = runDelta(params, opts)
	case modeSingature:
		logrus.Debug("signature mode selected")
		err = runSignature(params, opts)
	default:
		printUsage()
		// save weak sum
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
GoRDIFF - rdiff-like tool written in Go(lang)

Usage:
	%s [OPTIONS] signature [input_file [signature_file]]
	%s [OPTIONS] delta signature_file [updated_file [delta_file]]

Options:
	-verbose 	enable printing logs for each execution step
	-force		overwrite output file if exists

	-help		print this message

Notes:
	* For all '*_file' parameters the '--' string can be used instead of file paths. GoRDIFF will then use
	standard input or standard output (depending on file role) instead of files.
	* It is not allowed to use '--' for both input files of 'delta' mode
	* It is not allowed to use already existing file as output file (overwriting not possible)

`,
		os.Args[0], os.Args[0])
}

func runDelta(params []string, opts *appOptions) error {
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

	deltaFile, err := files.GetOutputFile(params[3], opts.forceMode)
	if err != nil {
		return fmt.Errorf("failed to open file for delta: %w", err)
	}
	defer deltaFile.Close()

	return hashing.GenerateDelta(sigFile, newFile, deltaFile)
}

func runSignature(params []string, opts *appOptions) error {
	if len(params) != 3 {
		printUsage()
		return common.ErrInvalidParams
	}

	inFile, err := files.GetInputFile(params[1])
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	sigFile, err := files.GetOutputFile(params[2], opts.forceMode)
	if err != nil {
		return fmt.Errorf("failed to open file for signature: %w", err)
	}
	defer sigFile.Close()

	signer := hashing.NewSigGen(2048)

	return signer.Compute(inFile, sigFile)
}
