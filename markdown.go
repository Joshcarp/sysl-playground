package main

import (
	"fmt"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/spf13/afero"
)

var mychan = make(chan string, 10000)
var mGlobal *Markdown

func main() {
	vecty.SetTitle("sysl Demo")

	vecty.RenderBody(&PageView{
		Input: `
MobileApp:
        Login:
                Server <- Login
        !type LoginData:
                username <: string
                password <: string
        !type LoginResponse:
                message <: string
Server:
        Login(data <: MobileApp.LoginData):
                return MobileApp.LoginResponse`,
	})
	keepAlive()
}

// PageView is our main page component.
type PageView struct {
	vecty.Core
	Input string
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		// Display a textarea on the right-hand side of the page.
		elem.Div(
			vecty.Markup(
				vecty.Style("float", "right"),
			),
			elem.TextArea(
				vecty.Markup(
					vecty.Style("font-family", "monospace"),
					vecty.Property("rows", 14),
					vecty.Property("cols", 70),

					// When input is typed into the textarea, update the local
					// component state and rerender.
					event.Input(func(e *vecty.Event) {
						p.Input = e.Target.Get("value").String()
						vecty.Rerender(p)
					}),
				),
				vecty.Text(p.Input), // initial textarea text.
			),
		),

		// Render the markdown.
		&Markdown{Input: p.Input},
	)
}

// Markdown is a simple component which renders the Input markdown as sanitized
// HTML into a div.
type Markdown struct {
	vecty.Core
	Input string `vecty:"prop"`
}

// Render implements the vecty.Component interface.
func (m *Markdown) Render() (res vecty.ComponentOrHTML) {
	defer func() {
		if r := recover(); r != nil {
			res = elem.Div(
				vecty.Markup(
					vecty.UnsafeHTML(fmt.Sprintf("%s", r)),
				),
			)
		}
	}()
	fs := afero.NewMemMapFs()
	f, err := fs.Create("/tmp.sysl")
	check(err)

	_, e := f.Write([]byte(m.Input))
	check(e)

	// function definition
	// exposing to JS
	// js.Global().Set("add", js.FuncOf(example))

	// var logger = logrus.New()
	// this := decimal.MustParse64(m.Input)
	// command.Main2([]string{"sysl", "sd", "-o", "project.svg", "-s", "MobileApp <- Login", "tmp.sysl"}, fs, logger, command.Main3)
	// http.Get("https://httpbin.org/get")


	// svg, err := fs.Open("project.svg")
	// check(err)
	// fmt.Println(svg)
	// this := make([]byte, 10000)
	// svg.Read(this)
	// keepAlive()
	return elem.Div(
		vecty.Markup(
			vecty.UnsafeHTML("string(this)"),
		),
	)
}

// func runSysl(m *Markdown) (res vecty.ComponentOrHTML) {

// }
func this() {
	fmt.Println("yes")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var signal = make(chan int)

func keepAlive() {
	for {
		<-signal
	}
}

// // Render implements the vecty.Component interface.
// func (m *Markdown) Render2() (res vecty.ComponentOrHTML) {
// 	defer func() {
// 		if r := recover(); r != nil {

// 			res = elem.Div(
// 				vecty.Markup(
// 					vecty.UnsafeHTML(fmt.Sprintf("%s", r)),
// 				),
// 			)
// 		}
// 	}()
// 	fs := afero.NewMemMapFs()
// 	f, err := fs.Create("/tmp.sysl")
// 	check(err)

// 	_, e := f.Write([]byte(m.Input))
// 	check(e)

// 	var logger = logrus.New()

// 	// if rc != 0 {
// 	// 	panic(rc)
// 	// }
// 	svg, err := fs.Open("/project.svg")
// 	check(err)
// 	fmt.Println(svg)
// 	this := make([]byte, 10000)
// 	svg.Read(this)

// 	return elem.Div(
// 		vecty.Markup(
// 			vecty.UnsafeHTML(string(this)),
// 		),
// 	)
// }

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
