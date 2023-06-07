package response

var (
	ErrBadRequest   = NewResponseError(400, "bad_request")
	ErrUnauthorized = NewResponseError(401, "unauthorized")
	ErrForbidden    = NewResponseError(403, "forbidden")
	ErrNotFound     = NewResponseError(404, "not_found")
	ErrConflict     = NewResponseError(409, "conflict")
	ErrInternal     = NewResponseError(500, "internal_server_error")
)

var (
	OK = NewResponseOk(200)
)
