package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/exercise/beer/model"
	"github.com/exercise/beer/service"
	"github.com/gin-gonic/gin"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

var (
	ErrRespFailedQueryRead = model.ErrorResponse{Message: "Read query parameters failed"}
	ErrRespFailedBodyRead  = model.ErrorResponse{Message: "Read body failed"}
)

type BeerHandler struct {
	beerService service.BeerService
}

func NewBeerHandler(beerService service.BeerService) BeerHandler {
	return BeerHandler{beerService: beerService}

}

func (b BeerHandler) FindAll(c *gin.Context) {

	q := model.PagingQuery{
		Limit:  10,
		Offset: 0,
	}

	limit := c.Query("limit")
	offset := c.Query("offset")
	if (limit != "") || (offset != "") {
		Limit, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrRespFailedQueryRead)

			return
		}

		Offset, err := strconv.ParseInt(offset, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrRespFailedQueryRead)

			return
		}

		if Limit != 0 {
			q.Limit = int(Limit)
		}

		if Offset != 0 {
			q.Offset = int(Offset)
		}

	}

	search := c.Query("search")

	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusBadRequest, ErrRespFailedQueryRead)

		return
	}

	beers, count, err := b.beerService.FindAll(search, q.Limit, q.Offset)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Find beers failed",
		})

		return
	}

	c.JSON(http.StatusOK, model.GetsResponse{
		Message: "success",
		Result:  beers,
		PagingQuery: model.RespPagingQuery{
			Offset: q.Offset,
			Limit:  q.Limit,
			Total:  count,
		},
	})

}

func (b BeerHandler) Create(c *gin.Context) {

	var req model.Beer

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)

		var validationErrs validator.ValidationErrors
		if ok := errors.As(err, &validationErrs); ok {
			c.JSON(http.StatusBadRequest, ErrRespFailedBodyRead)
			return
		}

		c.JSON(http.StatusBadRequest, ErrRespFailedBodyRead)

		return
	}

	beers, err := b.beerService.Create(req)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Find beers failed",
		})

		return
	}

	_, err = b.Logs(c)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Create log failed",
		})

		return
	}

	c.JSON(http.StatusCreated, model.DefaultResponse{
		Message: "success",
		Result:  beers,
	})

}

func (b BeerHandler) Update(c *gin.Context) {

	beerId, err := strconv.ParseInt(c.Param("beerId"), 10, 64)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: "User ID must be an integer",
		})

		return
	}

	var reqBody model.Beer
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		_ = c.Error(err)

		var validationErrs validator.ValidationErrors
		if errors.As(err, &validationErrs) {
			c.JSON(http.StatusBadRequest, ErrRespFailedBodyRead)

			return
		}

		c.JSON(http.StatusBadRequest, ErrRespFailedBodyRead)

		return
	}

	findBeer, err := b.beerService.Find(uint(beerId))
	if errors.Is(err, service.ErrBeerNotFound) {
		_ = c.Error(err)

		c.JSON(http.StatusNotFound,
			model.ErrorResponse{Message: "User not found"})

		return
	}

	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError,
			model.ErrorResponse{Message: "Update user failed"})

		return
	}

	findBeer.Name = reqBody.Name
	findBeer.Description = reqBody.Description
	findBeer.Category = reqBody.Category
	findBeer.Image = reqBody.Image

	beer, err := b.beerService.Update(findBeer)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError,
			model.ErrorResponse{Message: "Update beer failed"})

		return
	}

	_, err = b.Logs(c)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Create log failed",
		})

		return
	}

	c.JSON(http.StatusCreated, model.DefaultResponse{
		Message: "success",
		Result:  beer,
	})
}

func (b BeerHandler) Delete(c *gin.Context) {
	beerId, err := strconv.ParseInt(c.Param("beerId"), 10, 64)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid beer ID"})

		return
	}

	beer, err := b.beerService.Find(uint(beerId))
	if errors.Is(err, service.ErrBeerNotFound) {
		_ = c.Error(err)

		c.JSON(http.StatusNotFound,
			service.ErrBeerNotFound)

		return
	}

	err = b.beerService.Delete(beer.ID)

	if errors.Is(err, service.ErrBeerNotFound) {
		_ = c.Error(err)

		c.JSON(http.StatusNotFound,
			model.ErrorResponse{Message: "Beer not found"})

		return
	}

	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError,
			model.ErrorResponse{Message: "Update beer failed"})

		return
	}

	_, err = b.Logs(c)
	if err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Message: "Create log failed",
		})

		return
	}

	c.Status(http.StatusNoContent)
}

func (b BeerHandler) Upload(c *gin.Context) {

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("./images/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := "http://localhost:8080/file/" + filename
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
	var req struct {
		Type  string                `form:"type"`
		Image *multipart.FileHeader `form:"image"`
	}

	if err := c.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		_ = c.Error(err)

		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Message: err.Error(),
		})

		return
	}

}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (b BeerHandler) Logs(c *gin.Context) (interface{}, error) {

	var log model.Logs
	log.TimeStamp = time.Now().Format("02-01-2006 15:04:05")
	log.ClientIP = c.ClientIP()
	log.Method = c.Request.Method
	log.Path = c.Request.RequestURI
	log.StatusCode = c.Writer.Status()
	Logs, err := b.beerService.CreateLogOnMongoDB(log)

	return Logs, err

}
