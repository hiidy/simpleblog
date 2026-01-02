package errorsx

import (
	"errors"
	"fmt"
	"net/http"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

type ErrorX struct {
	Code     int               `json:"code,omitempty"`
	Reason   string            `json:"reason,omitempty"`
	Message  string            `json:"message,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

func (err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", err.Code, err.Reason, err.Message, err.Metadata)
}

func (err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

func (err *ErrorX) WithMetadata(md map[string]string) *ErrorX {
	err.Metadata = md
	return err
}

func (err *ErrorX) KV(kvs ...string) *ErrorX {
	if err.Metadata == nil {
		err.Metadata = make(map[string]string)
	}

	for i := 0; i < len(kvs); i += 2 {
		if i+1 < len(kvs) {
			err.Metadata[kvs[i]] = kvs[i+1]
		}
	}

	return err
}

func (err *ErrorX) GRPCStatus() *status.Status {
	details := errdetails.ErrorInfo{Reason: err.Reason, Metadata: err.Metadata}

	s, _ := status.New(httpstatus.ToGRPCCode(err.Code), err.Message).WithDetails(&details)
	return s
}

func (err *ErrorX) WithRequestID(requestID string) *ErrorX {
	return err.KV("X-Request-ID", requestID)
}

func (err *ErrorX) Is(target error) bool {
	if errx := new(ErrorX); errors.As(target, &errx) {
		return errx.Code == err.Code && errx.Reason == err.Reason
	}
	return false
}

func Code(err error) int {
	if err == nil {
		return http.StatusOK //nolint:mnd
	}
	return FromError(err).Code
}

func Reason(err error) string {
	if err == nil {
		return ErrInternal.Reason
	}
	return FromError(err).Reason
}

func FromError(err error) *ErrorX {
	if err == nil {
		return nil
	}

	if errx := new(ErrorX); errors.As(err, &errx) {
		return errx
	}

	gs, ok := status.FromError(err)
	if !ok {
		return New(ErrInternal.Code, ErrInternal.Reason, err.Error())
	}

	ret := New(httpstatus.FromGRPCCode(gs.Code()), ErrInternal.Reason, gs.Message())

	for _, detail := range gs.Details() {
		if typed, ok := detail.(*errdetails.ErrorInfo); ok {
			ret.Reason = typed.Reason
			return ret.WithMetadata(typed.Metadata)
		}
	}

	return ret
}
