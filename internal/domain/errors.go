package domain

import "errors"

var (
	ErrContextNotFound = errors.New("context not found")

	ErrSaveAudit = errors.New("error save event to audit")

	ErrSaveContext = errors.New("error save message to context")
)
