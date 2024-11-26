package httpregistry

import (
	"encoding/json"
)

// The list of all status codes is available at
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status

// Information responses
var (
	ContinueResponse           = NewResponse().WithStatus(100)
	SwitchingProtocolsResponse = NewResponse().WithStatus(101)
	ProcessingResponse         = NewResponse().WithStatus(102)
	EarlyHintsResponse         = NewResponse().WithStatus(103)
)

// Successful responses
var (
	OkResponse                          = NewResponse().WithStatus(200)
	CreatedResponse                     = NewResponse().WithStatus(201)
	AcceptedResponse                    = NewResponse().WithStatus(202)
	NonAuthoritativeInformationResponse = NewResponse().WithStatus(203)
	NoContentResponse                   = NewResponse().WithStatus(204)
	ResetContentResponse                = NewResponse().WithStatus(205)
	PartialContentResponse              = NewResponse().WithStatus(206)
	MultiStatusResponse                 = NewResponse().WithStatus(207)
	AlreadyReportedResponse             = NewResponse().WithStatus(208)
)

// Redirection messages
var (
	MultipleChoicesResponse   = NewResponse().WithStatus(300)
	MovedPermanentlyResponse  = NewResponse().WithStatus(301)
	FoundResponse             = NewResponse().WithStatus(302)
	SeeOtherResponse          = NewResponse().WithStatus(303)
	NotModifiedResponse       = NewResponse().WithStatus(304)
	TemporaryRedirectResponse = NewResponse().WithStatus(307)
	PermanentRedirectResponse = NewResponse().WithStatus(308)
)

// Client error responses
var (
	BadRequestsResponse                 = NewResponse().WithStatus(400)
	UnauthorizedResponse                = NewResponse().WithStatus(401)
	PaymentRequiredResponse             = NewResponse().WithStatus(402)
	ForbiddenResponse                   = NewResponse().WithStatus(403)
	NotFoundResponse                    = NewResponse().WithStatus(404)
	MethodNotAllowedResponse            = NewResponse().WithStatus(405)
	NotAcceptableResponse               = NewResponse().WithStatus(406)
	ProxyAuthenticationRequiredResponse = NewResponse().WithStatus(407)
	RequestTimeoutResponse              = NewResponse().WithStatus(408)
	ConflictResponse                    = NewResponse().WithStatus(409)
	GoneResponse                        = NewResponse().WithStatus(410)
	LengthRequiredResponse              = NewResponse().WithStatus(411)
	PreconditionFailedResponse          = NewResponse().WithStatus(412)
	PayloadTooLargeResponse             = NewResponse().WithStatus(413)
	URITooLongResponse                  = NewResponse().WithStatus(414)
	UnsupportedMediaTypeResponse        = NewResponse().WithStatus(415)
	RangeNotSatisfiableResponse         = NewResponse().WithStatus(416)
	ExpectationFailedResponse           = NewResponse().WithStatus(417)
	IAmATeapotResponse                  = NewResponse().WithStatus(418)
	MisdirectedRequestResponse          = NewResponse().WithStatus(421)
	UpgradeRequiredResponse             = NewResponse().WithStatus(426)
	ReconditionRequiredResponse         = NewResponse().WithStatus(428)
	RequestHeaderFieldsTooLargeResponse = NewResponse().WithStatus(431)
	UnavailableForLegalReasonsResponse  = NewResponse().WithStatus(451)
)

// Server error responses
var (
	InternalServerErrorResponse     = NewResponse().WithStatus(500)
	NotImplementedResponse          = NewResponse().WithStatus(501)
	BadGatewayResponse              = NewResponse().WithStatus(502)
	ServiceUnavailableResponse      = NewResponse().WithStatus(503)
	GatewayTimeoutResponse          = NewResponse().WithStatus(504)
	HTTPVersionNotSupportedResponse = NewResponse().WithStatus(505)
	VariantAlsoNegotiatesResponse   = NewResponse().WithStatus(506)
	NotExtendedResponse             = NewResponse().WithStatus(510)
	NetworkAuthenticationResponse   = NewResponse().WithStatus(511)
)

// Response represents a response that we want to return if the registry finds a request that matches the incoming request.
// If the match happens then we will return a http response that matches the attributes defined in this struct.
type Response struct {
	StatusCode int               `json:"status_code"`
	Body       []byte            `json:"body,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

func (r Response) String() string {
	bytes, err := json.Marshal(r)
	if err != nil {
		panic("cannot marshal request")
	}
	return string(bytes)
}

// WithStatus returns a new response with the StatusCode attribute set to statusCode
func (r Response) WithStatus(statusCode int) Response {
	r.StatusCode = statusCode
	return r
}

// WithHeader returns a new response with the header header set to value
func (r Response) WithHeader(header string, value string) Response {
	r.Headers[header] = value
	return r
}

// WithJSONHeader returns a new Response with the header `Content-Type` set to `application/json`
func (r Response) WithJSONHeader() Response {
	r.Headers["Content-Type"] = "application/json"
	return r
}

// WithHeaders returns a new response with all the headers in headers applied.
// If multiple headers with the same name are defined only the last one is applied.
func (r Response) WithHeaders(headers map[string]string) Response {
	for k, v := range headers {
		r.Headers[k] = v
	}
	return r
}

// WithBody returns a new request with the method body set to body
func (r Response) WithBody(body []byte) Response {
	r.Body = body
	return r
}

// Responses represents a slices of responses
type Responses = []Response

// NewResponse creates a new Response.
// This function is designed to be used in conjunction with other other receivers.
// For example
//
//	NewResponse().
//		WithStatus(http.StatusOK).
//		WithJSONHeader().
//		WithBody([]byte("{\"user\": \"John Schmidt\"}"))
func NewResponse() Response {
	r := Response{
		StatusCode: 0,
		Body:       make([]byte, 0),
		Headers:    make(map[string]string),
	}

	return r
}
