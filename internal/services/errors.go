package services

import "errors"

var (
	ErrSendingWalletNotFound   = errors.New("sending wallet not found")
	ErrReceivingWalletNotFound = errors.New("receiving wallet not found")
)
