package wmiext

import "github.com/pkg/errors"

// VM Lookup errors

var (
	NotFound = errors.New("not found")

	// VM Lookup errors
	// NoResults = errors.New("no results found")
	Failed   = errors.New("failed")
	TimedOut = errors.New("timed out")

	InvalidInput   error = errors.New("Invalid Input")
	InvalidType    error = errors.New("Invalid Type")
	NotSupported   error = errors.New("Not Supported")
	AlreadyExists  error = errors.New("Already Exists")
	InvalidFilter  error = errors.New("Invalid Filter")
	NotImplemented error = errors.New("Not Implemented")
	Unknown        error = errors.New("Unknown Reason")

	PermissionDenied error = errors.New("Permission Denied")
)
