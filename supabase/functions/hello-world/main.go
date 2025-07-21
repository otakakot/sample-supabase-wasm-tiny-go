//go:generate sh -c "GOOS=js GOARCH=wasm tinygo build -o main.wasm ./main.go && cat main.wasm | deno run https://denopkg.com/syumai/binpack/mod.ts > mainwasm.ts"
package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"syscall/js"
)

var done = make(chan struct{})

func init() {
	js.Global().Set("handle", js.FuncOf(handle))
}

func handle(_ js.Value, args []js.Value) any {
	defer func() {
		done <- struct{}{}
	}()

	str := js.Global().Get("JSON").Call("stringify", args[0]).String()

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(strings.NewReader(str)).Decode(&body); err != nil {
		slog.Error(err.Error())

		return nil
	}

	return ToJSResponse(
		http.StatusOK,
		http.Header{"Content-Type": []string{"application/json"}},
		[]byte(`{"message":"Hello `+body.Name+`"}`),
	)
}

func main() {
	<-done
}

func ToJSResponse(
	statusCode int,
	headers http.Header,
	data []byte,
) js.Value {
	respInit := js.Global().Get("Object")

	respInit.Set("status", statusCode)

	respInit.Set("statusText", http.StatusText(statusCode))

	respInit.Set("headers", ToJSHeader(headers))

	dataJs := js.Global().Get("Uint8Array").New(len(data))

	js.CopyBytesToJS(dataJs, data)

	return js.Global().Get("Response").New(dataJs, respInit)
}

// ToJSHeader converts http.Header to JavaScript sides Headers.
//   - Headers: https://developer.mozilla.org/docs/Web/API/Headers
func ToJSHeader(header http.Header) js.Value {
	h := js.Global().Get("Headers").New()
	for key, values := range header {
		for _, value := range values {
			h.Call("append", key, value)
		}
	}
	return h
}
