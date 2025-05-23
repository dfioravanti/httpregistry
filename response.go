package httpregistry

import (
	"encoding/json"
	"net/http"
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

// createResponse emits the response encoded in Response to w
func (res Response) createResponse(w http.ResponseWriter, _ *http.Request) {
	for k, v := range res.Headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(res.StatusCode)
	_, err := w.Write(res.Body)
	if err != nil {
		panic("cannot write body of request")
	}
}

// String marshal Response to string
func (res Response) String() string {
	bytes, err := json.Marshal(res)
	if err != nil {
		panic("cannot marshal request")
	}
	return string(bytes)
}

// WithStatus returns a new response with the StatusCode attribute set to statusCode
func (res Response) WithStatus(statusCode int) Response {
	res.StatusCode = statusCode
	return res
}

// WithHeader returns a new response with the header header set to value
func (res Response) WithHeader(header string, value string) Response {
	res.Headers[header] = value
	return res
}

// WithJSONHeader returns a new Response with the header `Content-Type` set to `application/json`
func (res Response) WithJSONHeader() Response {
	res.Headers["Content-Type"] = "application/json"
	return res
}

// WithHeaders returns a new response with all the headers in headers applied.
// If multiple headers with the same name are defined only the last one is applied.
func (res Response) WithHeaders(headers map[string]string) Response {
	for k, v := range headers {
		res.Headers[k] = v
	}
	return res
}

// WithBody returns a new request with the method body set to body
func (res Response) WithBody(body []byte) Response {
	res.Body = body
	return res
}

// WithJSONBody returns a new response that will return body as body and will have
// the header `Content-Type` set to `application/json`.
// This method panics if body cannot be converted to JSON
func (res Response) WithJSONBody(body any) Response {
	res = res.WithJSONHeader()
	res.Body = mustMarshalJSON(body)
	return res
}

// NewResponse creates a new Response.
// This function is designed to be used in conjunction with other other receivers.
// For example
//
//	NewResponse().
//		WithStatus(http.StatusOK).
//		WithJSONHeader().
//		WithBody([]byte("{\"user\": \"John Schmidt\"}"))
//
// The default response is a 200 without any body nor header
func NewResponse() Response {
	r := Response{
		StatusCode: 200,
		Body:       make([]byte, 0),
		Headers:    make(map[string]string),
	}

	return r
}
