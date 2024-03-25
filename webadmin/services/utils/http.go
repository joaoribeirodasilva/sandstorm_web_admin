package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
)

var (
	STATUS = map[int]string{
		http.StatusContinue:                      "continue",
		http.StatusSwitchingProtocols:            "switching protocols",
		http.StatusProcessing:                    "processing",
		http.StatusEarlyHints:                    "early hints",
		http.StatusOK:                            "ok",
		http.StatusCreated:                       "created",
		http.StatusAccepted:                      "accepted",
		http.StatusNonAuthoritativeInfo:          "non authorative info",
		http.StatusNoContent:                     "no content",
		http.StatusResetContent:                  "reset content",
		http.StatusPartialContent:                "partial content",
		http.StatusMultiStatus:                   "multi status",
		http.StatusAlreadyReported:               "already reported",
		http.StatusIMUsed:                        "im used",
		http.StatusMultipleChoices:               "multiple choices",
		http.StatusMovedPermanently:              "moved permanently",
		http.StatusFound:                         "found",
		http.StatusSeeOther:                      "see other",
		http.StatusNotModified:                   "not modified",
		http.StatusUseProxy:                      "use proxy",
		http.StatusTemporaryRedirect:             "temporary redirect",
		http.StatusPermanentRedirect:             "permanent redirect",
		http.StatusBadRequest:                    "bad request",
		http.StatusUnauthorized:                  "unauthorized",
		http.StatusPaymentRequired:               "payment required",
		http.StatusForbidden:                     "forbidden",
		http.StatusNotFound:                      "not found",
		http.StatusMethodNotAllowed:              "method not allowed",
		http.StatusNotAcceptable:                 "not acceptable",
		http.StatusProxyAuthRequired:             "proxy auth required",
		http.StatusRequestTimeout:                "request timeout",
		http.StatusConflict:                      "conflict",
		http.StatusGone:                          "gone",
		http.StatusLengthRequired:                "length required",
		http.StatusPreconditionFailed:            "precondition failed",
		http.StatusRequestEntityTooLarge:         "request entity too large",
		http.StatusRequestURITooLong:             "request uri too large",
		http.StatusUnsupportedMediaType:          "unsuported media type",
		http.StatusRequestedRangeNotSatisfiable:  "request rang not satisfiable",
		http.StatusExpectationFailed:             "expectation failed",
		http.StatusTeapot:                        "I'm a tea pot",
		http.StatusMisdirectedRequest:            "misdirected request",
		http.StatusUnprocessableEntity:           "unprocessable entity",
		http.StatusLocked:                        "locked",
		http.StatusFailedDependency:              "failed dependency",
		http.StatusTooEarly:                      "too early",
		http.StatusUpgradeRequired:               "upgrade required",
		http.StatusPreconditionRequired:          "precondition required",
		http.StatusTooManyRequests:               "too many requests",
		http.StatusRequestHeaderFieldsTooLarge:   "request header fields too large",
		http.StatusUnavailableForLegalReasons:    "unavailable for legal reasons",
		http.StatusInternalServerError:           "internal server error",
		http.StatusNotImplemented:                "not implemented",
		http.StatusBadGateway:                    "bad gateway",
		http.StatusServiceUnavailable:            "service unavailable",
		http.StatusGatewayTimeout:                "gateway timeout",
		http.StatusHTTPVersionNotSupported:       "http version not supported",
		http.StatusVariantAlsoNegotiates:         "varian also negotiates",
		http.StatusInsufficientStorage:           "insufficient storage",
		http.StatusLoopDetected:                  "loop detected",
		http.StatusNotExtended:                   "not extended",
		http.StatusNetworkAuthenticationRequired: "network authentication required",
	}
)

func Download(url string, dest string, log *admin_log.Log) error {

	if log != nil {
		log.Write(fmt.Sprintf("downloading file '%s'", url), MODULE, admin_log.LOG_DEBUG)
	}

	file, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0660)
	if err != nil {
		return log.Write(fmt.Sprintf("failed to create dedstination file '%s'. ERR: %s", dest, err.Error()), MODULE, admin_log.LOG_ERROR)
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return log.Write(fmt.Sprintf("failed to download file '%s'. ERR: %s", url, err.Error()), MODULE, admin_log.LOG_ERROR)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return log.Write(fmt.Sprintf("failed to download file '%s'. ERR: %s", url, fmt.Sprintf("server returned status [%d] %s", resp.StatusCode, STATUS[resp.StatusCode])), MODULE, admin_log.LOG_ERROR)
	}

	if _, err = io.Copy(file, resp.Body); err != nil {
		return log.Write(fmt.Sprintf("failed to save downloaded file '%s'. ERR: %s", url, err.Error()), MODULE, admin_log.LOG_ERROR)
	}

	return nil
}
