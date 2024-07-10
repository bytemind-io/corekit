package zrohandler

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/bytemind-io/corekit/zrohandler/response"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/color"
	"github.com/zeromicro/go-zero/core/iox"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/core/timex"
	"github.com/zeromicro/go-zero/core/utils"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	limitBodyBytes       = 1024
	defaultSlowThreshold = time.Millisecond * 500
)

var slowThreshold = syncx.ForAtomicDuration(defaultSlowThreshold)

// LogHandler returns a middleware that logs http request and response.
func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := utils.NewElapsedTimer()
		logs := new(LogCollector)
		lrw := response.NewWithCodeResponseWriter(w)

		var dup io.ReadCloser
		r.Body, dup = iox.LimitDupReadCloser(r.Body, limitBodyBytes)
		next.ServeHTTP(lrw, r.WithContext(WithLogCollector(r.Context(), logs)))
		r.Body = dup
		logBrief(r, lrw.Code, timer, logs)
	})
}

type detailLoggedResponseWriter struct {
	writer *response.WithCodeResponseWriter
	buf    *bytes.Buffer
}

func newDetailLoggedResponseWriter(writer *response.WithCodeResponseWriter,
	buf *bytes.Buffer) *detailLoggedResponseWriter {
	return &detailLoggedResponseWriter{
		writer: writer,
		buf:    buf,
	}
}

func (w *detailLoggedResponseWriter) Flush() {
	w.writer.Flush()
}

func (w *detailLoggedResponseWriter) Header() http.Header {
	return w.writer.Header()
}

// Hijack implements the http.Hijacker interface.
// This expands the Response to fulfill http.Hijacker if the underlying http.ResponseWriter supports it.
func (w *detailLoggedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacked, ok := w.writer.Writer.(http.Hijacker); ok {
		return hijacked.Hijack()
	}

	return nil, nil, errors.New("server doesn't support hijacking")
}

func (w *detailLoggedResponseWriter) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *detailLoggedResponseWriter) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}

// DetailedLogHandler returns a middleware that logs http request and response in details.
func DetailedLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := utils.NewElapsedTimer()
		var buf bytes.Buffer
		rw := response.NewWithCodeResponseWriter(w)
		lrw := newDetailLoggedResponseWriter(rw, &buf)

		var dup io.ReadCloser
		r.Body, dup = iox.DupReadCloser(r.Body)
		logs := new(LogCollector)
		next.ServeHTTP(lrw, r.WithContext(WithLogCollector(r.Context(), logs)))
		r.Body = dup
		logDetails(r, lrw, timer, logs)
	})
}

// SetSlowThreshold sets the slow threshold.
func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func dumpRequest(r *http.Request) string {
	reqContent, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err.Error()
	}

	return string(reqContent)
}

func isOkResponse(code int) bool {
	// not server error
	return code < http.StatusInternalServerError
}

func logBrief(r *http.Request, code int, timer *utils.ElapsedTimer, logs *LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	logger := logx.WithContext(r.Context()).WithDuration(duration)

	logger = logger.WithFields(
		logx.Field("status_code", wrapStatusCode(code)),
		logx.Field("method", wrapMethod(r.Method)),
		logx.Field("request_path", r.RequestURI),
		logx.Field("remote_ip", httpx.GetRemoteAddr(r)),
		logx.Field("user_agent", r.UserAgent()),
	)
	if duration > slowThreshold.Load() {
		logger = logger.WithFields(
			logx.Field("status_code", wrapStatusCode(code)),
			logx.Field("method", wrapMethod(r.Method)),
			logx.Field("request_path", r.RequestURI),
			logx.Field("remote_ip", httpx.GetRemoteAddr(r)),
			logx.Field("user_agent", r.UserAgent()),
			logx.Field("took", timex.ReprOfDuration(duration)),
			logx.Field("slowcall", "slowcall"),
		)
	}

	ok := isOkResponse(code)
	if !ok {
		buf.WriteString(fmt.Sprintf("\n%s", dumpRequest(r)))
	}

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("\n%s", body))
	}

	if ok {
		logger.Info(buf.String())
	} else {
		logger.Error(buf.String())
	}
}

func logDetails(r *http.Request, response *detailLoggedResponseWriter, timer *utils.ElapsedTimer,
	logs *LogCollector) {
	var buf bytes.Buffer
	duration := timer.Duration()
	code := response.writer.Code
	logger := logx.WithContext(r.Context())

	logger = logger.WithFields(
		logx.Field("status_code", wrapStatusCode(code)),
		logx.Field("method", wrapMethod(r.Method)),
		logx.Field("request_path", r.RequestURI),
		logx.Field("remote_ip", httpx.GetRemoteAddr(r)),
		logx.Field("user_agent", r.UserAgent()),
		logx.Field("duration", timex.ReprOfDuration(duration)),
		logx.Field("dump_request", dumpRequest(r)),
	)

	body := logs.Flush()
	if len(body) > 0 {
		buf.WriteString(fmt.Sprintf("%s\n", body))
	}

	if isOkResponse(code) {
		logger.Info(buf.String())
	} else {
		logger.Error(buf.String())
	}
}

func wrapMethod(method string) string {
	var colour color.Color
	switch method {
	case http.MethodGet:
		colour = color.BgBlue
	case http.MethodPost:
		colour = color.BgCyan
	case http.MethodPut:
		colour = color.BgYellow
	case http.MethodDelete:
		colour = color.BgRed
	case http.MethodPatch:
		colour = color.BgGreen
	case http.MethodHead:
		colour = color.BgMagenta
	case http.MethodOptions:
		colour = color.BgWhite
	}

	if colour == color.NoColor {
		return method
	}

	return logx.WithColorPadding(method, colour)
}

func wrapStatusCode(code int) string {
	var colour color.Color
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		colour = color.BgGreen
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		colour = color.BgBlue
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		colour = color.BgMagenta
	default:
		colour = color.BgYellow
	}

	return logx.WithColorPadding(strconv.Itoa(code), colour)
}
