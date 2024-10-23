package errs

import (
	"net/http"
)

var (
	// OK indicates the operation was successful.
	OK = ErrCode{value: 0}
	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument = ErrCode{value: 1}

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound = ErrCode{value: 2}

	// Internal errors. Means some invariants expected by underlying
	// system has been broken. If you see one of these errors,
	// something is very broken.
	Internal = ErrCode{value: 3}
)

var codeNames = map[ErrCode]string{
	OK:              "ok",
	InvalidArgument: "invalid_argument",
	NotFound:        "not_found",
	Internal:        "internal",
}

var httpStatus = map[ErrCode]int{
	OK:              http.StatusOK,
	InvalidArgument: http.StatusBadRequest,
	NotFound:        http.StatusNotFound,
	Internal:        http.StatusInternalServerError,
}
