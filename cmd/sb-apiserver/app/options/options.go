package options

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/hiidy/simpleblog/internal/apiserver"
	genericoptions "github.com/hiidy/simpleblog/pkg/options"
)

var availableServerModes = sets.New(
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
	apiserver.GinServerMode,
)

type ServerOptions struct {
	ServerMode  string                      `json:"server-mode" mapstructure:"server-mode"`
	JWTKey      string                      `json:"jwt-key" mapstructure:"jwt-key"`
	Expiration  time.Duration               `json:"expiration" mapstructure:"expiration"`
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`
}

func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:  "grpc-gateway",
		JWTKey:      "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:  2 * time.Hour,
		GRPCOptions: genericoptions.NewGRPCOptions(),
	}
	opts.GRPCOptions.Addr = ":6666"
	return opts
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long")
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "The expiration duration of JWT Tokens")
	o.GRPCOptions.AddFlags(fs)
}

func (o *ServerOptions) Validate() error {
	errs := []error{}

	// ServerMode가 유효한지 검증
	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: must be one of %v", availableServerModes.UnsortedList()))
	}

	// JWTKey 길이 검증
	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("JWTKey must be at least 6 characters long"))
	}

	if slices.Contains([]string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}, o.ServerMode) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}
	// 모든 오류를 병합하여 반환
	return utilerrors.NewAggregate(errs)
}

func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:  o.ServerMode,
		JWTKey:      o.JWTKey,
		Expiration:  o.Expiration,
		GRPCOptions: o.GRPCOptions,
	}, nil
}
