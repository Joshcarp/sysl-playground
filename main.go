package main

import (
	"fmt"

	"github.com/Joshcarp/sysl-playground/pkg/syslUtil"
	"github.com/Joshcarp/sysl-playground/pkg/urls"
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
)

// PageView is our main page component.
type PageView struct {
	vecty.Core
	Input     string
	Command   string
	Link      string
	InputLink string
}

func main() {
	input, cmd := setup()

	vecty.SetTitle("SYSL Playground")
	vecty.RenderBody(&PageView{
		Input:   input,
		Command: cmd,
	})
}

func setup() (string, string) {
	playgroundUrl, _ := urls.LoadQueryParams()
	input, cmd := urls.DecodeQueryParams(playgroundUrl)

	if input == "" {
		input = `MobileApp:
        Login:
                        Server <- Login
        !type LoginData:
                        username <: string
                        password <: string
        !type LoginResponse:
                        message <: string
Server:
        Login(data <: MobileApp.LoginData):
                        return MobileApp.LoginResponse`
	}
	if cmd == "" {
		cmd = "sysl sd -o \"project.svg\" -s \"MobileApp <- Login\" tmp.sysl"
	}
	return input, cmd
}

// Render implements the vecty.Component interface.
func (p *PageView) Render() vecty.ComponentOrHTML {
	return elem.Body(
		vecty.Markup(
			vecty.Class("body"),
		),
		elem.Header(
			vecty.Text("Sysl Playground"),
			vecty.Markup(
				vecty.Style("font-family", "monospace"),
				vecty.Style("font-size", "25px"),
			),
		),
		elem.Article(
			vecty.Text("Welcome to the Sysl Playground"),
			vecty.Markup(
				vecty.Style("font-family", "monospace"),
			),
		),

		// Display a textarea on the right-hand side of the page.
		elem.Table(
			elem.TableRow(
				elem.TableData(
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
				elem.TableData(
					&Markdown{Input: p.Input, Command: p.Command},
				),
			),
			elem.TableRow(
				elem.TableData(
					elem.TextArea(
						vecty.Markup(
							vecty.Style("font-family", "monospace"),
							vecty.Property("rows", 1),
							vecty.Property("cols", 70),

							// When input is typed into the textarea, update the local
							// component state and rerender.
							event.Input(func(e *vecty.Event) {
								p.Command = e.Target.Get("value").String()
								vecty.Rerender(p)
							}),
						),
						vecty.Text(p.Command), // initial textarea text.
					),
				),
			)),
		elem.TableRow(
			elem.Button(
				vecty.Markup(
					vecty.UnsafeHTML("Share"),
					event.Click(func(e *vecty.Event) {
						p.Link = urls.EncodeUrl(p.Input, p.Command)
						vecty.Rerender(p)
					}),
				),
			),
		),
		elem.TableRow(
			elem.TableData(
				elem.TextArea(
					vecty.Markup(
						vecty.Style("font-family", "monospace"),
						vecty.Property("rows", 7),
						vecty.Property("cols", 70),
						vecty.Property("wrap", "hard"),
						event.Input(func(e *vecty.Event) {
							p.InputLink = e.Target.Get("value").String()
						}),
					),
					vecty.Text(p.Link),
				),
			),
		),
	)
}

// Markdown is a simple component which renders the Input markdown as sanitized
// HTML into a div.
type Markdown struct {
	vecty.Core
	Input   string `vecty:"prop"`
	Command string `vecty:"prop"`
}

// Render implements the vecty.Component interface.
func (m *Markdown) Render() (res vecty.ComponentOrHTML) {
	defer func() {
		// If panic, then print the error
		if r := recover(); r != nil {
			res = elem.Div(
				vecty.Markup(
					vecty.UnsafeHTML(fmt.Sprintf("%s", r)),
				),
			)
		}
	}()

	output, err := syslUtil.Parse(m.Input, m.Command)
	check(err)
	image := fmt.Sprintf("<img src=\"%s\">", string(output))

	return elem.Div(
		vecty.Markup(
			vecty.UnsafeHTML(
				image),
		),
	)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
