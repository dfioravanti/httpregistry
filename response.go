package httpregistry

import "encoding/json"

// Response represents a response that we want to return if the registry finds a request that matches the incoming request.
// If the match happens then we will return a http response that matches the attributes defined in this struct.
type Response struct {
	Status  int               `json:"status"`
	Body    []byte            `json:"body,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

func (r Response) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic("cannot marshal request")
	}
	return string(bytes)
}

// Responses represents a slices of responses
type Responses = []Response

// ResponseOption represents a option that can be passed to NewResponse when creating a new response.
// NewResponse uses the Option patters to make it easy to configure the response behavior.
// For example
//
//	NewResponse(100, nil, WithResponseHeader("Content-Type", "application/json"))
//
// will return a response with the desired content type
type ResponseOption func(*Response)

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

// WithResponseHeader allows to add any header to a response.
// If multiple headers with the same name are defined only the last one is applied.
// To define multiple headers it is recommended to use WithResponseHeaders but it is possible to chain multiple calls of WithResponseHeader.
func WithResponseHeader(header string, value string) func(*Response) {
	return func(r *Response) {
		r.Headers[header] = value
	}
}

// WithResponseHeaders allows to add any number of headers to a response.
// If multiple headers with the same name are defined only the last one is applied.
func WithResponseHeaders(headers map[string]string) ResponseOption {
	return func(r *Response) {
		for k, v := range headers {
			r.Headers[k] = v
		}
	}
}

// NewJSONResponse creates a new Response with the desired statusCode, body and the various options applied,
// whose "Content-Type" header is set to "application/json".
// The JSON content type will overwrite any "Content-Type" header that is set via options.
func NewJSONResponse(statusCode int, body []byte, options ...ResponseOption) Response {
	r := Response{
		Status:  statusCode,
		Body:    body,
		Headers: make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}
	r.Headers["Content-Type"] = "application/json"
	return r
}

// NewResponse creates a new Response with the desired statusCode, body and the various options applied.
func NewResponse(statusCode int, body []byte, options ...func(*Response)) Response {
	r := Response{
		Status:  statusCode,
		Body:    body,
		Headers: make(map[string]string),
	}
	for _, o := range options {
		o(&r)
	}

	return r
}
