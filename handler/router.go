package handler

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mongodb-job/internal/render"
	"github.com/mongodb-job/internal/worker_manager"
	"github.com/mongodb-job/mcode"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var httpClient = &http.Client{}

func Start() {

	gin.SetMode(gin.ReleaseMode)

	router := loadRestRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 80),
		Handler:        router,
		ReadTimeout:    300 * time.Second,
		WriteTimeout:   300 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
func loadRestRouter() *gin.Engine {
	router := gin.Default()

	// health api
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := router.Group("/inner/v1")
	v1.POST("/job/save", HandleJobSave)
	v1.DELETE("/job/delete", HandlerJobDelete)
	v1.POST("/job/update", HandleJobUpdate)
	v1.DELETE("/job/kill", HandlerJobKill)

	forgien := router.Group("/api/v1")
	forgien.GET("/worker", func(c *gin.Context) {
		c.JSON(http.StatusOK, worker_manager.WorkerManagerInstance.All())
	})
	// 获取所有任务接口
	forgien.GET("/pipeline", ShowPipeline)
	// 获取节点上的任务接口
	forgien.GET("/plan", ShowWorkerPipeline)
	// 获取任务历史记录接口
	forgien.GET("/snapshot", ShowPipelineHistory)
	forgien.POST("/job/save", HandleJobSave)
	forgien.DELETE("/job/delete", HandlerJobDelete)
	forgien.POST("/job/update", HandleJobUpdate)

	forgien.POST("/file", func(ctx *gin.Context) {
		r := render.New(ctx)
		var form newForm
		if err := ctx.ShouldBind(&form); err != nil {
			r.JSON("", mcode.RequestErr)
			return
		}
		dst := form.Target + form.UploadKey.Filename
		if err := ctx.SaveUploadedFile(form.UploadKey, dst); err != nil {
			r.JSON("", err)
			return
		}
		waitSyncFile, err := os.Open(dst)
		if err != nil {
			r.JSON("", err)
			return
		}
		defer func() {
			waitSyncFile.Close()
		}()

		// sync to worker TODO 当前节点不多，后续需要改成并发
		for _, node := range worker_manager.WorkerManagerInstance.All() {
			//node.Host
			url := "http://" + node.Host + "/inner/api/v1/worker/file"

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("upload", filepath.Base(dst))
			if err != nil {
				logrus.Print(err)
				r.JSON("", err)
				return
			}
			_, err = io.Copy(part, waitSyncFile)
			_ = writer.WriteField("target", form.Target)
			writer.Close()

			request, err := http.NewRequest("POST", url, body)
			if err != nil {
				logrus.Print(err)
				r.JSON("", err)
			}
			request.Header.Add("Content-Type", writer.FormDataContentType())
			response, err := httpClient.Do(request)
			if err != nil {
				logrus.Print(err)
				r.JSON("", err)
			}
			responseByte, _ := ioutil.ReadAll(response.Body)
			logrus.Println(string(responseByte))
			request.Body.Close()
		}

		r.JSON("", mcode.OK)
		return
	})

	forgien.GET("/log", JobLogSearch)
	forgien.GET("/worker/status", WorkerStatus)
	forgien.GET("/job/migrate", JobMigrate)
	return router
}

func StartWorker() {

	gin.SetMode(gin.ReleaseMode)

	router := loadWorkerRestRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", 80),
		Handler:        router,
		ReadTimeout:    300 * time.Second,
		WriteTimeout:   300 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}
func loadWorkerRestRouter() *gin.Engine {
	router := gin.Default()

	// health api
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := router.Group("/inner/v1")
	v1.POST("/worker/job", WorkersJob)
	v1.POST("/worker/file", func(ctx *gin.Context) {
		r := render.New(ctx)
		var form newForm
		if err := ctx.ShouldBind(&form); err != nil {
			r.JSON("", mcode.RequestErr)
			return
		}
		dst := form.Target + form.UploadKey.Filename
		if err := ctx.SaveUploadedFile(form.UploadKey, dst); err != nil {
			r.JSON("", err)
			return
		}
		r.JSON("", mcode.OK)
		return
	})

	return router
}

type newForm struct {
	UploadKey *multipart.FileHeader `form:"upload"`
	Target    string                `form:"target"`
}
