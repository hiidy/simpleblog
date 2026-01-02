package errno

import (
	"net/http"

	"github.com/hiidy/simpleblog/pkg/errorsx"
)

var ErrPostNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFOund.PostNotFound", Message: "Post not found"}
