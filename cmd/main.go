package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pdkovacs/fork-fem-htmx-proj/internal/incrementor"
)

const (
	indexView                 = "index"
	blocksView                = "blocks"
	memoryConsumptionTestView = "memory-consumption-test"
)

type Templates struct {
	templates map[string]*template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	parts := strings.Split(name, "/")
	fmt.Printf(">>>>>>>>> name: %s, parts[0]: %s, parts[1]: %s\n", name, parts[0], parts[1])
	return t.templates[parts[0]].ExecuteTemplate(w, parts[1], data)
}

func NewTemplates() *Templates {
	tmpl := make(map[string]*template.Template)
	tmpl[indexView] = template.Must(template.ParseFiles("views/index.html", "views/scripts.html", "views/base.html"))
	tmpl[blocksView] = template.Must(template.ParseFiles("views/blocks/index.html", "views/scripts.html", "views/base.html"))
	tmpl[memoryConsumptionTestView] = template.Must(template.ParseGlob("views/memory-consumption-test/*"))
	tmpl[memoryConsumptionTestView].ParseFiles("views/scripts.html", "views/base.html")
	return &Templates{
		templates: tmpl,
	}
}

type Block struct {
	Id int
}

type Blocks struct {
	Start  int
	Next   int
	More   bool
	Blocks []Block
}

func main() {
	e := echo.New()
	e.Renderer = NewTemplates()

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index/base", Blocks{})
	})

	e.GET("/index", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index/base", Blocks{})
	})

	e.GET("/memory-consumption-test", func(c echo.Context) error {
		return c.Render(http.StatusOK, fmt.Sprintf("%s/base", memoryConsumptionTestView), struct{}{})
	})

	e.POST("/memory-consumption-test/start", func(c echo.Context) error {
		form, formErr := c.FormParams()
		if formErr != nil {
			fmt.Printf(">>>>> formErr: %#v\n", formErr)
			return formErr
		}
		fmt.Printf(">>>>>>> increment_size: %s increment_interval: %s\n", form["increment_size"], form["increment_interval"])
		if len(form["increment_size"]) == 0 {
			fmt.Printf(">>>> increment_size missing")
			return fmt.Errorf("increment_size missing")
		}
		if len(form["increment_interval"]) == 0 {
			fmt.Printf(">>>> increment_interval missing")
			return fmt.Errorf("increment_interval missing")
		}
		incrementSize, sizeParseErr := strconv.Atoi(form["increment_size"][0])
		if sizeParseErr != nil {
			fmt.Printf(">>>> The value of the increment_size field %s is not a valid integer: %v", form["increment_size"][0], sizeParseErr)
			return sizeParseErr
		}
		incrementInterval, intervalParseErr := strconv.Atoi(form["increment_interval"][0])
		if intervalParseErr != nil {
			fmt.Printf(">>>> The value of the increment_interval field %s is not a valid integer: %v", form["increment_interval"][0], intervalParseErr)
			return intervalParseErr
		}
		incrementor.StartIncrementing(incrementSize, incrementInterval)
		return c.Render(http.StatusOK, fmt.Sprintf("%s/in-progress", memoryConsumptionTestView), struct{ Consumed int }{Consumed: incrementor.GetConsumed()})
	})

	e.GET("/memory-consumption-test/poll", func(c echo.Context) error {
		return c.Render(http.StatusOK, fmt.Sprintf("%s/consumed", memoryConsumptionTestView), struct{ Consumed int }{Consumed: incrementor.GetConsumed()})
	})

	e.POST("/memory-consumption-test/stop", func(c echo.Context) error {
		incrementor.SuspendIncrementing()
		return c.Render(http.StatusOK, fmt.Sprintf("%s/start-button", memoryConsumptionTestView), struct{}{})
	})

	e.GET("/blocks", func(c echo.Context) error {
		startStr := c.QueryParam("start")
		start, err := strconv.Atoi(startStr)
		if err != nil {
			start = 0
		}

		blocks := []Block{}
		for i := start; i < start+10; i++ {
			blocks = append(blocks, Block{Id: i})
		}

		template := "blocks"
		if start == 0 {
			template = "blocks-index"
		}
		return c.Render(http.StatusOK, template, Blocks{
			Start:  start,
			Next:   start + 10,
			More:   start+10 < 100,
			Blocks: blocks,
		})
	})

	e.Logger.Fatal(e.Start(":42069"))
}
