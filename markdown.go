package main

import (
	"fmt"
	"net/http"
	// "syscall/js"

	"github.com/Joshcarp/sysl_testing/pkg/command"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

var mychan = make(chan string, 10000)
var mGlobal *Markdown
var info *http.Response

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
Server2:
        Login(data <: MobileApp.LoginData):
                return MobileApp.LoginResponse`,
	})
	// go keepAlive()
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

	var logger = logrus.New()
	command.Main2([]string{"sysl", "sd", "-o", "project.svg", "-s", "MobileApp <- Login", "tmp.sysl"}, fs, logger, command.Main3)

	// svg, err := fs.Open("project.svg")
	// check(err)
	// fmt.Println(svg)
	// this := make([]byte,0, 10000)
	this, err := afero.ReadFile(fs, "project.svg")
	check(err)

	foo := fmt.Sprintf("<img src=\"%s\">", string(this))
	fmt.Println(foo)
	return elem.Div(
		vecty.Markup(
			vecty.UnsafeHTML(
				foo),
		),
	)
}

// func keepAlive() {
// 	example := func(this js.Value, i []js.Value) interface{} {
// 		go func() {
// 			info, _ = http.Get("https://httpbin.org/get")
// 		}()
// 		return nil
// 	}
// 	js.Global().Set("example", js.FuncOf(example))
// 	select {}
// }

// 	js.Global().Set("example", js.FuncOf(example))
// 	select {}
// }

// func runSysl(m *Markdown) (res vecty.ComponentOrHTML) {

// }
func this() {
	fmt.Println("yes")
}

func check(err error) {
	if err != nil {
		// panic(err)
	}
}

var signal = make(chan int)

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
