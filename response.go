package httpregistry

import (
	"net/http"
)

// The list of all status codes is available at
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status

// Information responses
var (
	ContinueResponse           = newResponseWithName("httpregistry.ContinueResponse").WithStatus(100)
	SwitchingProtocolsResponse = newResponseWithName("httpregistry.SwitchingProtocolsResponse").WithStatus(101)
	ProcessingResponse         = newResponseWithName("httpregistry.ProcessingResponse").WithStatus(102)
	EarlyHintsResponse         = newResponseWithName("httpregistry.EarlyHintsResponse").WithStatus(103)
)

// Successful responses
var (
	OkResponse                          = newResponseWithName("httpregistry.OkResponse").WithStatus(200)
	CreatedResponse                     = newResponseWithName("httpregistry.CreatedResponse").WithStatus(201)
	AcceptedResponse                    = newResponseWithName("httpregistry.AcceptedResponse").WithStatus(202)
	NonAuthoritativeInformationResponse = newResponseWithName("httpregistry.NonAuthoritativeInformationResponse").WithStatus(203)
	NoContentResponse                   = newResponseWithName("httpregistry.NoContentResponse").WithStatus(204)
	ResetContentResponse                = newResponseWithName("httpregistry.ResetContentResponse").WithStatus(205)
	PartialContentResponse              = newResponseWithName("httpregistry.PartialContentResponse").WithStatus(206)
	MultiStatusResponse                 = newResponseWithName("httpregistry.MultiStatusResponse").WithStatus(207)
	AlreadyReportedResponse             = newResponseWithName("httpregistry.AlreadyReportedResponse").WithStatus(208)
)

// Redirection messages
var (
	MultipleChoicesResponse   = newResponseWithName("httpregistry.MultipleChoicesResponse").WithStatus(300)
	MovedPermanentlyResponse  = newResponseWithName("httpregistry.MovedPermanentlyResponse").WithStatus(301)
	FoundResponse             = newResponseWithName("httpregistry.FoundResponse").WithStatus(302)
	SeeOtherResponse          = newResponseWithName("httpregistry.SeeOtherResponse").WithStatus(303)
	NotModifiedResponse       = newResponseWithName("httpregistry.NotModifiedResponse").WithStatus(304)
	TemporaryRedirectResponse = newResponseWithName("httpregistry.TemporaryRedirectResponse").WithStatus(307)
	PermanentRedirectResponse = newResponseWithName("httpregistry.PermanentRedirectResponse").WithStatus(308)
)

// Client error responses
var (
	BadRequestsResponse                 = newResponseWithName("httpregistry.BadRequestsResponse").WithStatus(400)
	UnauthorizedResponse                = newResponseWithName("httpregistry.UnauthorizedResponse").WithStatus(401)
	PaymentRequiredResponse             = newResponseWithName("httpregistry.PaymentRequiredResponse").WithStatus(402)
	ForbiddenResponse                   = newResponseWithName("httpregistry.ForbiddenResponse").WithStatus(403)
	NotFoundResponse                    = newResponseWithName("httpregistry.NotFoundResponse").WithStatus(404)
	MethodNotAllowedResponse            = newResponseWithName("httpregistry.MethodNotAllowedResponse").WithStatus(405)
	NotAcceptableResponse               = newResponseWithName("httpregistry.NotAcceptableResponse").WithStatus(406)
	ProxyAuthenticationRequiredResponse = newResponseWithName("httpregistry.ProxyAuthenticationRequiredResponse").WithStatus(407)
	RequestTimeoutResponse              = newResponseWithName("httpregistry.RequestTimeoutResponse").WithStatus(408)
	ConflictResponse                    = newResponseWithName("httpregistry.ConflictResponse").WithStatus(409)
	GoneResponse                        = newResponseWithName("httpregistry.GoneResponse").WithStatus(410)
	LengthRequiredResponse              = newResponseWithName("httpregistry.LengthRequiredResponse").WithStatus(411)
	PreconditionFailedResponse          = newResponseWithName("httpregistry.PreconditionFailedResponse").WithStatus(412)
	PayloadTooLargeResponse             = newResponseWithName("httpregistry.PayloadTooLargeResponse").WithStatus(413)
	URITooLongResponse                  = newResponseWithName("httpregistry.URITooLongResponse").WithStatus(414)
	UnsupportedMediaTypeResponse        = newResponseWithName("httpregistry.UnsupportedMediaTypeResponse").WithStatus(415)
	RangeNotSatisfiableResponse         = newResponseWithName("httpregistry.RangeNotSatisfiableResponse").WithStatus(416)
	ExpectationFailedResponse           = newResponseWithName("httpregistry.ExpectationFailedResponse").WithStatus(417)
	IAmATeapotResponse                  = newResponseWithName("httpregistry.IAmATeapotResponse").WithStatus(418)
	MisdirectedRequestResponse          = newResponseWithName("httpregistry.MisdirectedRequestResponse").WithStatus(421)
	UpgradeRequiredResponse             = newResponseWithName("httpregistry.UpgradeRequiredResponse").WithStatus(426)
	ReconditionRequiredResponse         = newResponseWithName("httpregistry.ReconditionRequiredResponse").WithStatus(428)
	RequestHeaderFieldsTooLargeResponse = newResponseWithName("httpregistry.RequestHeaderFieldsTooLargeResponse").WithStatus(431)
	UnavailableForLegalReasonsResponse  = newResponseWithName("httpregistry.UnavailableForLegalReasonsResponse").WithStatus(451)
)

// Server error responses
var (
	InternalServerErrorResponse     = newResponseWithName("httpregistry.InternalServerErrorResponse").WithStatus(500)
	NotImplementedResponse          = newResponseWithName("httpregistry.NotImplementedResponse").WithStatus(501)
	BadGatewayResponse              = newResponseWithName("httpregistry.BadGatewayResponse").WithStatus(502)
	ServiceUnavailableResponse      = newResponseWithName("httpregistry.ServiceUnavailableResponse").WithStatus(503)
	GatewayTimeoutResponse          = newResponseWithName("httpregistry.GatewayTimeoutResponse").WithStatus(504)
	HTTPVersionNotSupportedResponse = newResponseWithName("httpregistry.HTTPVersionNotSupportedResponse").WithStatus(505)
	VariantAlsoNegotiatesResponse   = newResponseWithName("httpregistry.VariantAlsoNegotiatesResponse").WithStatus(506)
	NotExtendedResponse             = newResponseWithName("httpregistry.NotExtendedResponse").WithStatus(510)
	NetworkAuthenticationResponse   = newResponseWithName("httpregistry.NetworkAuthenticationResponse").WithStatus(511)
)

// Response represents a response that we want to return if the registry finds a request that matches the incoming request.
// If the match happens then we will return a http response that matches the attributes defined in this struct.
type Response struct {
	name       string
	statusCode int
	body       []byte
	headers    map[string]string
}

// serveResponse emits the response encoded in Response to w
func (res Response) serveResponse(w http.ResponseWriter, _ *http.Request) {
	for k, v := range res.headers {
		w.Header().Add(k, v)
	}
	w.WriteHeader(res.statusCode)
	_, err := w.Write(res.body)
	if err != nil {
		panic("cannot write body of request")
	}
}

// WithName allows to add a name to a Response so that it can be better identified when debugging.
// By the default Response gets a sequential name that can be hard to identify if there are many of them.
// So if clarity is needed we recommend to change the default name.
func (res Response) WithName(name string) Response {
	res.name = name
	return res
}

// String marshal Response to string
func (res Response) String() string {
	return res.name
}

// WithStatus returns a new response with the StatusCode attribute set to statusCode
func (res Response) WithStatus(statusCode int) Response {
	res.statusCode = statusCode
	return res
}

// WithHeader returns a new response with the header header set to value
func (res Response) WithHeader(header string, value string) Response {
	res.headers[header] = value
	return res
}

// WithJSONHeader returns a new Response with the header `Content-Type` set to `application/json`
func (res Response) WithJSONHeader() Response {
	res.headers["Content-Type"] = "application/json"
	return res
}

// WithHeaders returns a new response with all the headers in headers applied.
// If multiple headers with the same name are defined only the last one is applied.
func (res Response) WithHeaders(headers map[string]string) Response {
	for k, v := range headers {
		res.headers[k] = v
	}
	return res
}

// WithBody returns a new request with the method body set to body
func (res Response) WithBody(body []byte) Response {
	res.body = body
	return res
}

// WithJSONBody returns a new response that will return body as body and will have
// the header `Content-Type` set to `application/json`.
// This method panics if body cannot be converted to JSON
func (res Response) WithJSONBody(body any) Response {
	res = res.WithJSONHeader()
	res.body = mustMarshalJSON(body)
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
	return newResponseWithName("")
}

// newResponseWithName creates a new Response with the name already set,
// this has the advantage of not increasing the counter in the default naming schema.
//
// This function is designed to be used in conjunction with other other receivers.
// For example
//
//	newResponseWithName("httpregistry.John response").
//		WithStatus(http.StatusOK).
//		WithJSONHeader().
//		WithBody([]byte("{\"user\": \"John Schmidt\"}"))
//
// The default response is a 200 without any body nor header
func newResponseWithName(name string) Response {
	r := Response{
		name:       name,
		statusCode: 200,
		body:       make([]byte, 0),
		headers:    make(map[string]string),
	}

	return r
}
