package main

const (
	exitSuccess          = 0
	exitInvalidOptions   = 1
	exitOperationFailure = 2
)

const (
	modeDelta     = "delta"
	modeSingature = "signature"
)

type appOptions struct {
	forceMode bool
	blocSize  int
}
