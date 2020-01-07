package main

import (
	"fmt"

	"github.com/Joshcarp/sysl-playground/attempt2/files"
	"github.com/Joshcarp/sysl_testing/pkg/command"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/sirupsen/logrus"
)

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
func (m *Markdown) Render() vecty.ComponentOrHTML {
	// Render the markdown input into HTML using Blackfriday.
	// unsafeHTML := blackfriday.Run([]byte(m.Input))

	var fs = files.Fs{}
	f, err := fs.Create("tmp.sysl")
	check(err)

	_, e := f.Write([]byte(m.Input))
	check(e)

	var logger = logrus.New()
	fmt.Println(logger)
	fmt.Println("this")
	// rc := 0
	rc := command.Main2([]string{"sysl", "pb", "-o", "project.textpb", "tmp.sysl"}, fs, logger, command.Main3)
	if rc != 0 {
		// panic(rc)
	}
	// g, err := fs.Create("project.svg")
	// check(err)
	// defer g.Close()
	svg, err := fs.Open("project.textpb")
	check(err)
	fmt.Println(svg)

	return elem.Div(
		vecty.Markup(
			vecty.UnsafeHTML("this isn't some text"),
		),
	)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
