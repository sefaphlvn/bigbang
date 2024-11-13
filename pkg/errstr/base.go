package errstr

import "errors"

var (
	ErrNotAuthorized         = errors.New("you are not authorized to perform this action")
	ErrNameAlreadyExists     = errors.New("name already exists")
	ErrUnknownDBError        = errors.New("unknown db error")
	ErrUserIDEmpty           = errors.New("userID cannot be empty")
	ErrNoDocumentsUpdate     = errors.New("document not found or not permission to update")
	ErrNoDocumentsDelete     = errors.New("document not found or not permission to delete")
	ErrNoDocuments           = errors.New("document not found")
	ErrListenerNotFound      = errors.New("listener not found")
	ErrInvalidIndexName      = errors.New("invalid index name")
	ErrInvalidVersion        = errors.New("invalid version format")
	ErrFailedToUpdateVersion = errors.New("failed to update resource version")
	ErrUnexpectedResource    = errors.New("unexpected resource format")
	ErrValidationFailed      = errors.New("validation failed")
	ErrInvalidRefreshToken   = errors.New("invalid refresh token")
)
