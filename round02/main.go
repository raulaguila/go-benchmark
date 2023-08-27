package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	"github.com/labstack/echo/v4"
	"github.com/savsgio/atreugo/v11"
)

const (
	paramOK    string = "paramOK"
	paramName  string = "channel"
	endpoint   string = "/bench"
	endpoint01 string = endpoint + "/{" + paramName + "}"
	endpoint02 string = endpoint + "/:" + paramName

	ginPort     string = ":8080"
	chiPort     string = ":8081"
	fiberPort   string = ":8082"
	gorillaPort string = ":8083"
	echoPort    string = ":8084"
	httpPort    string = ":8085"
	atreugoPort string = ":8086"
)

type ObjectExample struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func startFrameworks() {
	go initGinGonic()
	go initGoChi()
	go initGoFiber()
	go initGorillaMux()
	go initGoEcho()
	go initHttpServerMux()
	go initAtreugo()
}

func main() {
	log.Println("Starting frameworks")
	forever := make(chan bool)

	startFrameworks()

	<-forever
}

func initGinGonic() {
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()
	app.POST(endpoint02, func(c *gin.Context) {
		if value, ok := c.Params.Get(paramName); !ok || value != paramOK {
			c.Status(http.StatusBadRequest)
			return
		}

		objectExample := new(ObjectExample)
		if err := c.ShouldBindJSON(objectExample); err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, objectExample)
	})

	app.Run(ginPort)
}

func initGoChi() {
	app := chi.NewRouter()
	app.Post(endpoint01, func(w http.ResponseWriter, r *http.Request) {
		if channel := chi.URLParam(r, paramName); channel != paramOK {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		objectExample := new(ObjectExample)
		if err := json.NewDecoder(r.Body).Decode(objectExample); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		jsonValue, _ := json.Marshal(objectExample)
		w.Write(jsonValue)
	})

	http.ListenAndServe(chiPort, app)
}

func initGoFiber() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Post(endpoint02, func(c *fiber.Ctx) error {
		if value := c.Params(paramName); value != paramOK {
			return c.Status(http.StatusBadRequest).Send(nil)
		}

		objectExample := new(ObjectExample)
		if err := c.BodyParser(objectExample); err != nil {
			return c.Status(http.StatusBadRequest).Send(nil)
		}

		err := c.Status(http.StatusOK).JSON(objectExample)
		c.Set("content-type", "application/json; charset=utf-8")
		return err
	})

	app.Listen(fiberPort)
}

func initGorillaMux() {
	app := mux.NewRouter()
	app.HandleFunc(endpoint01, func(w http.ResponseWriter, r *http.Request) {
		if value, ok := mux.Vars(r)[paramName]; !ok || value != paramOK {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		objectExample := new(ObjectExample)
		if err := json.Unmarshal(b, objectExample); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		jsonValue, _ := json.Marshal(objectExample)
		w.Write(jsonValue)
	}).Methods("POST")

	http.ListenAndServe(gorillaPort, app)
}

func initGoEcho() {
	app := echo.New()
	app.HideBanner = true
	app.HidePort = true
	app.POST(endpoint02, func(c echo.Context) error {
		if value := c.Param(paramName); value != paramOK {
			return c.JSON(http.StatusBadRequest, nil)
		}

		objectExample := new(ObjectExample)
		if err := json.NewDecoder(c.Request().Body).Decode(objectExample); err != nil {
			return c.JSON(http.StatusBadRequest, nil)
		}

		return c.JSON(http.StatusOK, objectExample)
	})

	app.Start(echoPort)
}

func initHttpServerMux() {
	httpmux := http.NewServeMux()
	httpmux.HandleFunc(endpoint+"/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if value := strings.TrimPrefix(r.URL.Path, endpoint+"/"); value != paramOK {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(nil)
				return
			}

			objectExample := new(ObjectExample)
			if err := json.NewDecoder(r.Body).Decode(objectExample); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(nil)
				return
			}

			w.WriteHeader(http.StatusOK)
			resp, _ := json.Marshal(objectExample)
			w.Write(resp)
		}
	})
	http.ListenAndServe(httpPort, httpmux)
}

func initAtreugo() {
	server := atreugo.New(atreugo.Config{
		Addr:  "0.0.0.0" + atreugoPort,
		Debug: false,
	})

	server.POST(endpoint01, func(c *atreugo.RequestCtx) error {
		if value := c.UserValue(paramName).(string); value != paramOK {
			c.SetStatusCode(http.StatusBadRequest)
			return nil
		}

		objectExample := new(ObjectExample)
		if err := json.Unmarshal(c.PostBody(), objectExample); err != nil {
			c.SetStatusCode(http.StatusBadRequest)
			return nil
		}

		c.Response.Header.SetContentType("application/json")
		return c.JSONResponse(objectExample, http.StatusOK)
	})

	server.ListenAndServe()
}
