package apple

import "errors"

var (
	// ErrInvalidCertificate returns when parse a receipt
	// with invalid certificate from given root certificate.
	ErrInvalidCertificate = errors.New("invalid certificate in receipt")
	// ErrInvalidSignature returns when parse a receipt
	// which improperly signed.
	ErrInvalidSignature = errors.New("invalid signature of receipt")



)
