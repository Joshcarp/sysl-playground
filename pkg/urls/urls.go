package urls

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"syscall/js"
)

func EncodeUrl(input, cmd string) string {
	input = Encode(input)
	cmd = Encode(cmd)
	url := js.Global().Get("location").Get("hostname").String()
	port := js.Global().Get("location").Get("port").String()
	pathname := js.Global().Get("location").Get("pathname").String()

	if port != "" {
		port = ":" + port
	}
	pathname = strings.Replace(pathname, "/", "", 1)
	if pathname != "" {
		pathname = "/" + pathname
	}
	return fmt.Sprintf("http://%s%s/%s?input=%s&cmd=%s", url, port, pathname, input, cmd)

}

func LoadQueryParams() (url.Values, bool) {
	href := js.Global().Get("location").Get("href")
	str := fmt.Sprintf("%s", href)
	u, err := url.Parse(str)
	check(err)
	if len(u.Query()) == 0 {
		return u.Query(), false
	}
	return u.Query(), true
}

func Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func Decode(str string) string {
	this, _ := base64.StdEncoding.DecodeString(str)
	return string(this)
}
func DecodeQueryParams(in url.Values) (string, string) {
	foo := Decode(in.Get("input"))
	bar := Decode(in.Get("cmd"))
	return foo, bar
}
func DecodeURLString(in string) (string, string) {
	u, err := url.Parse(in)
	check(err)
	return DecodeQueryParams(u.Query())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
