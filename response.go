package servermock

type Response struct {
	body    []byte
	status  int
	headers map[string]string
}

type Responses = []Response

// The list of all status codes is available at
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status

// Information responses
var (
	ContinueResponse           = NewResponse(WithResponseStatus(100))
	SwitchingProtocolsResponse = NewResponse(WithResponseStatus(101))
	ProcessingResponse         = NewResponse(WithResponseStatus(102))
	EarlyHintsResponse         = NewResponse(WithResponseStatus(103))
)

// Successful responses
var (
	OkResponse                          = NewResponse(WithResponseStatus(200))
	CreatedResponse                     = NewResponse(WithResponseStatus(201))
	AcceptedResponse                    = NewResponse(WithResponseStatus(202))
	NonAuthoritativeInformationResponse = NewResponse(WithResponseStatus(203))
	NoContentResponse                   = NewResponse(WithResponseStatus(204))
	ResetContentResponse                = NewResponse(WithResponseStatus(205))
	PartialContentResponse              = NewResponse(WithResponseStatus(206))
)

// Redirection messages
var (
	MultipleChoicesResponse = NewResponse(WithResponseStatus(300))
)

// Client error responses
var (
	BadRequestsResponse      = NewResponse(WithResponseStatus(400))
	UnauthorizedResponse     = NewResponse(WithResponseStatus(401))
	PaymentRequiredResponse  = NewResponse(WithResponseStatus(402))
	ForbiddenResponse        = NewResponse(WithResponseStatus(403))
	NotFoundResponse         = NewResponse(WithResponseStatus(404))
	MethodNotAllowedResponse = NewResponse(WithResponseStatus(405))
	NotAcceptableResponse    = NewResponse(WithResponseStatus(406))
)

// Server error responses
var (
	InternalServerErrorResponse = NewResponse(WithResponseStatus(500))
)

func WithResponseBody(body []byte) func(*Response) {
	return func(r *Response) {
		r.body = body
	}
}

func WithResponseJSONBody(body []byte) func(*Response) {
	return func(r *Response) {
		r.headers["Content-Type"] = "application/json"

		r.body = body
	}
}

func WithResponseStatus(status int) func(*Response) {
	return func(r *Response) {
		r.status = status
	}
}

func WithResponseHeaders(headers map[string]string) func(*Response) {
	return func(r *Response) {
		for k, v := range headers {
			r.headers[k] = v
		}
	}
}

func WithResponseHeader(header string, value string) func(*Response) {
	return func(r *Response) {
		r.headers[header] = value
	}
}

func NewResponse(options ...func(*Response)) Response {
	r := Response{
		headers: make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}

	return r
}
