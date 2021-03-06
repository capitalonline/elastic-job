package render

import (
	"github.com/gin-gonic/gin"
	"github.com/mongodb-job/mcode"
	"net/http"
	"time"
)

type Gin struct {
	C *gin.Context
	T time.Time
}

// New create a new gin Context
func New(c *gin.Context) *Gin {
	return &Gin{
		C: c,
		T: time.Now(),
	}
}

// JSON common json struct.
type JSON struct {
	// 业务错误码
	Code int `json:"code"`
	// 错误描述
	Msg string `json:"msg"`
	// 响应时长(ms)
	TTL int64 `json:"ttl"`
	// 响应数据
	Data interface{} `json:"data"`

	TraceID string `json:"trace_id"`
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (g *Gin) JSON(data interface{}, err error) {
	code := mcode.Cause(err)
	var (
		ttl int64 = 0
	)
	if !g.T.IsZero() {
		ttl = time.Now().Sub(g.T).Nanoseconds() / 1e6
	}
	trace, ok := g.C.Value("traceId").(string)
	if !ok {
		trace = ""
	}
	g.C.JSON(http.StatusOK, JSON{
		Code:    code.Code(),
		Msg:     code.Message(),
		TTL:     ttl,
		Data:    data,
		TraceID: trace,
	})
}
