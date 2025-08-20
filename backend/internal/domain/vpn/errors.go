package vpn

import "errors"

var (
	ErrReqParsing          = errors.New("failed to parse req to JSON")
	ErrReqCreation         = errors.New("failed to create http POST request")
	ErrReq                 = errors.New("failed to make request")
	ErrReadResponse        = errors.New("failed to read response")
	ErrParsingResponse     = errors.New("failed to parse response")
	ErrCreatingCertificate = errors.New("failed to create vpn certificate")
)
