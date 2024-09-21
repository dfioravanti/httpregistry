package httpregistry

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
	ContinueResponse           = NewResponse(100, nil)
	SwitchingProtocolsResponse = NewResponse(101, nil)
	ProcessingResponse         = NewResponse(102, nil)
	EarlyHintsResponse         = NewResponse(103, nil)
)

// Successful responses
var (
	OkResponse                          = NewResponse(200, nil)
	CreatedResponse                     = NewResponse(201, nil)
	AcceptedResponse                    = NewResponse(202, nil)
	NonAuthoritativeInformationResponse = NewResponse(203, nil)
	NoContentResponse                   = NewResponse(204, nil)
	ResetContentResponse                = NewResponse(205, nil)
	PartialContentResponse              = NewResponse(206, nil)
	MultiStatusResponse                 = NewResponse(207, nil)
	AlreadyReportedResponse             = NewResponse(208, nil)
)

// Redirection messages
var (
	MultipleChoicesResponse   = NewResponse(300, nil)
	MovedPermanentlyResponse  = NewResponse(301, nil)
	FoundResponse             = NewResponse(302, nil)
	SeeOtherResponse          = NewResponse(303, nil)
	NotModifiedResponse       = NewResponse(304, nil)
	TemporaryRedirectResponse = NewResponse(307, nil)
	PermanentRedirectResponse = NewResponse(308, nil)
)

// Client error responses
var (
	BadRequestsResponse                 = NewResponse(400, nil)
	UnauthorizedResponse                = NewResponse(401, nil)
	PaymentRequiredResponse             = NewResponse(402, nil)
	ForbiddenResponse                   = NewResponse(403, nil)
	NotFoundResponse                    = NewResponse(404, nil)
	MethodNotAllowedResponse            = NewResponse(405, nil)
	NotAcceptableResponse               = NewResponse(406, nil)
	ProxyAuthenticationRequiredResponse = NewResponse(407, nil)
	RequestTimeoutResponse              = NewResponse(408, nil)
	ConflictResponse                    = NewResponse(409, nil)
	GoneResponse                        = NewResponse(410, nil)
	LengthRequiredResponse              = NewResponse(411, nil)
	PreconditionFailedResponse          = NewResponse(412, nil)
	PayloadTooLargeResponse             = NewResponse(413, nil)
	URITooLongResponse                  = NewResponse(414, nil)
	UnsupportedMediaTypeResponse        = NewResponse(415, nil)
	RangeNotSatisfiableResponse         = NewResponse(416, nil)
	ExpectationFailedResponse           = NewResponse(417, nil)
	IAmATeapotResponse                  = NewResponse(418, nil)
	MisdirectedRequestResponse          = NewResponse(421, nil)
	UpgradeRequiredResponse             = NewResponse(426, nil)
	ReconditionRequiredResponse         = NewResponse(428, nil)
	RequestHeaderFieldsTooLargeResponse = NewResponse(431, nil)
	UnavailableForLegalReasonsResponse  = NewResponse(451, nil)
)

// Server error responses
var (
	InternalServerErrorResponse     = NewResponse(500, nil)
	NotImplementedResponse          = NewResponse(501, nil)
	BadGatewayResponse              = NewResponse(502, nil)
	ServiceUnavailableResponse      = NewResponse(503, nil)
	GatewayTimeoutResponse          = NewResponse(504, nil)
	HTTPVersionNotSupportedResponse = NewResponse(505, nil)
	VariantAlsoNegotiatesResponse   = NewResponse(506, nil)
	NotExtendedResponse             = NewResponse(510, nil)
	NetworkAuthenticationResponse   = NewResponse(511, nil)
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

func NewJSONResponse(statusCode int, body []byte, options ...func(*Response)) Response {
	r := Response{
		status:  statusCode,
		body:    body,
		headers: map[string]string{"Content-Type": "application/json"},
	}
	for _, o := range options {
		o(&r)
	}

	return r
}

func NewResponse(statusCode int, body []byte, options ...func(*Response)) Response {
	r := Response{
		status:  statusCode,
		body:    body,
		headers: make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}

	return r
}
