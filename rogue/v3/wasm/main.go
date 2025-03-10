//go:build wasm
// +build wasm

package main

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"syscall/js"

	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

// This will be set by -ldflags
var ImageTag = "development"

var panicCallback js.Value

func registerPanicCallback(_ js.Value, args []js.Value) interface{} {
	panicCallback = args[0]
	return nil
}

func catchPanic() {
	if r := recover(); r != nil {
		fmt.Printf("[wasm] panic: %s\n%s", r, debug.Stack())
		if panicCallback.Type() == js.TypeFunction {
			panicCallback.Invoke(js.ValueOf(fmt.Sprintf("%s\n%s", r, debug.Stack())))
		}
	}
}

var done = make(chan struct{})

func jsValueToMap(obj js.Value) map[string]interface{} {
	result := make(map[string]interface{})
	keys := js.Global().Get("Object").Call("keys", obj)

	for i := 0; i < keys.Length(); i++ {
		key := keys.Index(i).String()
		value := obj.Get(key).String()
		result[key] = value
	}

	return result
}

func main() {
	js.Global().Set("RegisterPanicCallback", js.FuncOf(registerPanicCallback))

	js.Global().Set("NewRogue", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer catchPanic()
		author := args[0].String()
		instance := v3.NewRogueForQuill(author)

		return map[string]interface{}{
			"Version": js.FuncOf(func(this js.Value, args []js.Value) any {
				return "v3"
			}),
			"ImageTag": js.FuncOf(func(this js.Value, args []js.Value) any {
				return ImageTag
			}),
			"Panic": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				panic(args[0].String())
			}),
			"SetAuthor": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				author := args[0].String()
				instance.Author = author
				return nil
			}),
			"GetText": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				return instance.GetText()
			}),
			"Paste": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 4 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 4 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeNumber || args[1].Type() != js.TypeNumber || args[2].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected number, number and object, got %T, %T and %T", args[0], args[1], args[2]),
					}
				}

				visIx := args[0].Int()
				spanLen := args[1].Int()

				formatMap := jsValueToMap(args[2])
				formatMap["e"] = "true"
				formatMap["en"] = "true"
				format, err := v3.MapToFormatV3(formatMap)
				if err != nil {
					fmt.Printf("ERROR [wasm] MapToFormatV3(%v): %v\n", formatMap, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				spanFormat, ok := format.(v3.FormatV3Span)
				if !ok {
					return map[string]interface{}{
						"error": "expected span format",
					}
				}

				jsArray := args[3]
				if jsArray.Type() != js.TypeObject || !jsArray.InstanceOf(js.Global().Get("Array")) {
					return map[string]interface{}{
						"error": "Third argument must be an array",
					}
				}

				var pasteItems []v3.PasteItem
				arrayLen := jsArray.Length()
				for i := 0; i < arrayLen; i++ {
					jsObj := jsArray.Index(i)
					if jsObj.Type() != js.TypeObject {
						return map[string]interface{}{
							"error": fmt.Sprintf("Item at index %d is not an object", i),
						}
					}

					mime := jsObj.Get("mime")
					kind := jsObj.Get("kind")
					data := jsObj.Get("data")
					if mime.Type() != js.TypeString || kind.Type() != js.TypeString || data.Type() != js.TypeString {
						return map[string]interface{}{
							"error": fmt.Sprintf("Item at index %d has invalid properties", i),
						}
					}

					pasteItems = append(pasteItems, v3.PasteItem{
						Kind: kind.String(),
						Mime: mime.String(),
						Data: data.String(),
					})
				}

				ops, cursorID, err := instance.Paste(visIx, spanLen, spanFormat, pasteItems)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				opsArr := make([]interface{}, 0, len(ops))
				for _, op := range ops {
					opsArr = append(opsArr, op.AsArr())
				}

				return map[string]interface{}{
					"ops":      opsArr,
					"cursorID": cursorID.AsJS(),
				}
			}),
			"ToOps": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				ops, err := instance.ToOps()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				opsArr := make([]interface{}, 0, len(ops))
				for _, op := range ops {
					opsArr = append(opsArr, op.AsArr())
				}

				return opsArr
			}),
			"PasteHtml": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 2 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeNumber || args[1].Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected number and string, got %T and %T", args[0], args[1]),
					}
				}

				idx := args[0].Int()
				html := args[1].String()

				mop, err := instance.PasteHtml(idx, html)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"GetPlaintext": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 2 && len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2 or 3 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[0]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()

				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				endID, err := instance.AfterIDToEndID(afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				var plaintext string
				if len(args) == 3 {
					ca, err := ContentAddressArg(args, 2)
					if err != nil {
						return map[string]interface{}{
							"error": fmt.Sprintf("content address: %s", err.Error()),
						}
					}

					plaintext, err = instance.GetPlaintext(startID, endID, ca)
				} else {
					plaintext, err = instance.GetPlaintext(startID, endID, nil)
				}

				if err != nil {
					fmt.Printf("ERROR [wasm] GetPlaintext(%v, %v): %s\n", startID, endID, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return plaintext
			}),
			"GetHtml": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 2 && len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2-3 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[0]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()

				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				endID, err := instance.AfterIDToEndID(afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				showIds := true
				if len(args) == 3 {
					if args[2].Type() != js.TypeBoolean {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected boolean, got %T", args[2]),
						}
					}
					showIds = args[2].Bool()
				}

				html, err := instance.GetHtml(startID, endID, showIds, true)
				if err != nil {
					fmt.Printf("ERROR [wasm] GetHtml(%v, %v, %v): %s\n", startID, endID, showIds, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return html
			}),
			"GetHtmlXRay": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 2 && len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2-3 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[0]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()

				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				endID, err := instance.AfterIDToEndID(afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				showIds := true
				if len(args) == 3 {
					if args[2].Type() != js.TypeBoolean {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected boolean, got %T", args[2]),
						}
					}
					showIds = args[2].Bool()
				}

				html, err := instance.GetHtmlXRay(startID, endID, showIds, true)
				if err != nil {
					fmt.Printf("ERROR [wasm] GetHtml(%v, %v, %v): %s\n", startID, endID, showIds, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return html
			}),
			"GetHtmlAtAddress": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 3 && len(args) != 4 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 3 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[0]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()

				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				if args[2].Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %T", args[2]),
					}
				}

				address := args[2].String()

				ca := v3.ContentAddress{}
				err := json.Unmarshal([]byte(address), &ca)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				endID, err := instance.AfterIDToEndID(afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				showIds := true
				if len(args) == 4 {
					if args[3].Type() != js.TypeBoolean {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected boolean, got %T", args[3]),
						}
					}
					showIds = args[3].Bool()
				}

				html, err := instance.GetHtmlAt(startID, endID, &ca, showIds, true)
				if err != nil {
					fmt.Printf("ERROR [wasm] GetHtmlAt(%v, %v, %v): %s\n", startID, endID, ca, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return html
			}),
			"GetAddress": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				var err error
				startID, endID := v3.NoID, v3.NoID

				if len(args) == 0 {
					startID, err = instance.Rope.GetTotID(0)
					if err != nil {
						return map[string]interface{}{
							"error": err.Error(),
						}
					}

					endID, err = instance.Rope.GetTotID(instance.TotSize - 1)
					if err != nil {
						return map[string]interface{}{
							"error": err.Error(),
						}
					}
				} else if len(args) == 2 {
					if args[0].Type() != js.TypeObject {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected object, got %T", args[0]),
						}
					}

					sau := args[0].Index(0).String()
					slt := args[0].Index(1).Int()

					if args[1].Type() != js.TypeObject {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected object, got %T", args[1]),
						}
					}

					fau := args[1].Index(0).String()
					flt := args[1].Index(1).Int()

					startID = v3.ID{Author: sau, Seq: slt}
					endID = v3.ID{Author: fau, Seq: flt}
				} else {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 0 or 2 arguments, got %d", len(args)),
					}
				}

				ca, err := instance.GetAddress(startID, endID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				caBytes, err := json.Marshal(ca)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return string(caBytes)
			}),
			"GetHtmlDiffBetween": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 4 && len(args) != 5 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 4 or 5 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[0]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()

				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				if args[2].Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %T", args[2]),
					}
				}

				fromContentAddress := v3.ContentAddress{}
				err := json.Unmarshal([]byte(args[2].String()), &fromContentAddress)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				if args[3].Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %T", args[2]),
					}
				}

				toContentAddress := v3.ContentAddress{}
				err = json.Unmarshal([]byte(args[3].String()), &toContentAddress)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				endID, err := instance.AfterIDToEndID(afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				showIds := true
				if len(args) == 5 {
					if args[4].Type() != js.TypeBoolean {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected boolean, got %T", args[4]),
						}
					}
					showIds = args[4].Bool()
				}

				html, err := instance.GetHtmlDiffBetween(startID, endID, &fromContentAddress, &toContentAddress, showIds, true)
				if err != nil {
					fmt.Printf("ERROR [wasm] GetHtmlDiffBetween(%v, %v, %v, %v): %s\n", startID, endID, fromContentAddress, toContentAddress, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return html
			}),
			"GetHtmlDiff": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 3 && len(args) != 4 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 3 or 4 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[0]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()

				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				if args[2].Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %T", args[2]),
					}
				}

				address := args[2].String()

				ca := v3.ContentAddress{}
				err := json.Unmarshal([]byte(address), &ca)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				endID, err := instance.AfterIDToEndID(afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				showIds := true
				if len(args) == 4 {
					if args[3].Type() != js.TypeBoolean {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected boolean, got %T", args[3]),
						}
					}
					showIds = args[3].Bool()
				}

				html, err := instance.GetHtmlDiff(startID, endID, &ca, showIds, true)
				if err != nil {
					fmt.Printf("ERROR [wasm] GetHtmlDiff(%v, %v, %v): %s\n", startID, endID, ca, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return html
			}),
			"HighlightSpan": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()
				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				beforeID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				startIx, length, err := instance.HighlightSpan(beforeID, afterID)
				if err != nil {
					return map[string]interface{}{
						"error":  err.Error(),
						"index":  startIx,
						"length": length,
					}
				}

				return map[string]interface{}{
					"index":  startIx,
					"length": length,
				}
			}),
			"GetMarkdown": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()
				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				beforeID := v3.ID{Author: sau, Seq: slt}
				afterID := v3.ID{Author: fau, Seq: flt}

				md, err := instance.GetMarkdown(beforeID, afterID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
						"text":  "",
					}
				}

				return map[string]interface{}{
					"text": md,
				}
			}),
			"Insert": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				if len(args) != 2 {
					return map[string]interface{}{
						"error": "expected 2 args",
					}
				}

				if args[0].Type() != js.TypeNumber || args[1].Type() != js.TypeString {
					errStr := fmt.Sprintf("expected number and string got %s (%s) and %s (%s)",
						args[0].String(),
						args[0].Type().String(),
						args[1].String(),
						args[1].Type().String(),
					)
					return map[string]interface{}{
						"error": errStr,
					}
				}

				visIx := args[0].Int()
				text := args[1].String()

				if visIx < 0 {
					return map[string]interface{}{
						"error": "visIx < 0",
					}
				}

				result, err := instance.Insert(visIx, text)
				if err != nil {
					fmt.Printf("ERROR [wasm] Insert(%d, %q): %s\n", visIx, text, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}
				return result.AsArr()
			}),
			"RichInsert": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				if len(args) != 4 {
					return map[string]interface{}{
						"error": "expected 4 args",
					}
				}

				if args[0].Type() != js.TypeNumber || args[1].Type() != js.TypeNumber || args[3].Type() != js.TypeString {
					errStr := fmt.Sprintf("expected number, number and string got %s (%s), %s (%s) and %s (%s)",
						args[0].String(),
						args[0].Type().String(),
						args[1].String(),
						args[1].Type().String(),
						args[3].String(),
						args[3].Type().String(),
					)
					return map[string]interface{}{
						"error": errStr,
					}
				}

				formatMap := jsValueToMap(args[2])
				formatMap["e"] = "true"
				formatMap["en"] = "true"
				format, err := v3.MapToFormatV3(formatMap)
				if err != nil {
					fmt.Printf("ERROR [wasm] MapToFormatV3(%v): %v\n", formatMap, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				spanFormat, ok := format.(v3.FormatV3Span)
				if !ok {
					return map[string]interface{}{
						"error": "expected span format",
					}
				}

				visIx := args[0].Int()
				selLen := args[1].Int()
				text := args[3].String()

				if visIx < 0 {
					return map[string]interface{}{
						"error": "visIx < 0",
					}
				}

				if selLen < 0 {
					return map[string]interface{}{
						"error": "selLen < 0",
					}
				}

				ops, cursorID, err := instance.RichInsert(visIx, selLen, spanFormat, text)
				if err != nil {
					fmt.Printf("ERROR [wasm] RichInsert(%d, %d, %v, %q): %s\n", visIx, selLen, format, text, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				opsArr := make([]interface{}, 0, len(ops))
				for _, op := range ops {
					opsArr = append(opsArr, op.AsArr())
				}

				return map[string]interface{}{
					"ops":      opsArr,
					"cursorID": cursorID.AsJS(),
				}
			}),
			"RichDelete": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				if len(args) != 2 {
					return map[string]interface{}{
						"error": fmt.Sprintf("invalid args: %v", args),
					}
				}

				visIx := args[0].Int()
				length := args[1].Int()

				op, startIx, err := instance.RichDelete(visIx, length)
				if err != nil {
					fmt.Printf("ERROR [wasm] RichDelete err: %s\n", err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				if op == nil {
					return map[string]interface{}{}
				}

				return map[string]interface{}{
					"op":      op.AsArr(),
					"startIx": startIx,
				}
			}),
			"RichDeleteLine": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				if len(args) != 1 {
					return map[string]interface{}{
						"error": fmt.Sprintf("invalid args: %v", args),
					}
				}

				visIx := args[0].Int()

				op, startIx, err := instance.RichDeleteLine(visIx)
				if err != nil {
					fmt.Printf("ERROR [wasm] RichDeleteLine err: %s\n", err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				if op == nil {
					return map[string]interface{}{}
				}

				return map[string]interface{}{
					"op":      op.AsArr(),
					"startIx": startIx,
				}
			}),
			"RichDeleteWord": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				if len(args) != 2 {
					return map[string]interface{}{
						"error": fmt.Sprintf("invalid args: %v", args),
					}
				}

				visIx := args[0].Int()
				forward := args[1].Bool()

				op, startIx, err := instance.RichDeleteWord(visIx, forward)
				if err != nil {
					fmt.Printf("ERROR [wasm] RichDeleteWord err: %s\n", err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				if op == nil {
					return map[string]interface{}{}
				}

				return map[string]interface{}{
					"op":      op.AsArr(),
					"startIx": startIx,
				}
			}),
			"InsertMarkdown": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				visIx := args[0].Int()
				text := args[1].String()

				mop, actions, err := instance.InsertMarkdown(visIx, text)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Errorf("InsertMarkdown(%d, %q): %w", visIx, text, err),
					}
				}

				return map[string]interface{}{
					"mop":     mop.AsArr(),
					"actions": actions.AsJS(),
				}
			}),
			"Delete": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				if len(args) != 2 {
					return map[string]interface{}{
						"error": fmt.Sprintf("invalid args: %v", args),
					}
				}

				visIx := args[0].Int()
				length := args[1].Int()

				if visIx < 0 {
					return map[string]interface{}{
						"error": fmt.Sprintf("invalid visIx: %d", visIx),
					}
				}
				if length < 0 {
					return map[string]interface{}{
						"error": fmt.Sprintf("invalid length: %d", length),
					}
				}

				mop, err := instance.Delete(visIx, length)
				if err != nil {
					fmt.Printf("ERROR [wasm] Delete err: %s\n", err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"Format": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				visIx := args[0].Int()
				length := args[1].Int()
				format := jsValueToMap(args[2])
				f, err := v3.MapToFormatV3(format)
				if err != nil {
					fmt.Printf("ERROR [wasm] MapToFormatV3(%v): %v\n", format, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				mop, err := instance.Format(visIx, length, f)
				if err != nil {
					fmt.Printf("ERROR [wasm] Format(%d, %d, %v): %v\n", visIx, length, f, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"FormatLineByID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()
				id := v3.ID{Author: sau, Seq: slt}

				format := jsValueToMap(args[1])
				f, err := v3.MapToFormatV3(format)
				if err != nil {
					fmt.Printf("ERROR [wasm] MapToFormatV3(%v): %v\n", format, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				mop, err := instance.FormatLineByID(id, f)
				if err != nil {
					fmt.Printf("ERROR [wasm] FormatLineByID(%v, %v): %v\n", id, f, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"Undo": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()
				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				endID := v3.ID{Author: fau, Seq: flt}

				/*var err error
				if startID == endID {
					startID, err = instance.Rope.GetTotID(0)
					if err != nil {
						return map[string]interface{}{
							"error": err.Error(),
						}
					}

					endID, err = instance.Rope.GetTotID(instance.TotSize - 1)
					if err != nil {
						return map[string]interface{}{
							"error": err.Error(),
						}
					}
				}*/

				// TODO: when undoing we need to return new startID and endID spans
				// to update the selection in the editor to the correct location
				// potentially we could use this same patter for insert, delete and paste

				// just undo whole doc for now
				startID, err := instance.Rope.GetTotID(0)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				endID, err = instance.Rope.GetTotID(instance.TotSize - 1)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				mop, err := instance.Undo(startID, endID)
				if err != nil {
					fmt.Printf("ERROR [wasm] Undo(%v, %v): %s\n", startID, endID, err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"Redo": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				mop, err := instance.Redo()
				if err != nil {
					fmt.Printf("ERROR [wasm] Redo(): %s\n", err)
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"Rewind": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 3 arguments, got %d", len(args)),
					}
				}
				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %s", args[0].Type()),
					}
				}
				if args[1].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %s", args[1].Type()),
					}
				}
				addressArg := args[2]
				if addressArg.Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %s", addressArg.Type()),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()
				fau := args[1].Index(0).String()
				flt := args[1].Index(1).Int()
				addressStr := addressArg.String()

				var address v3.ContentAddress
				err := json.Unmarshal([]byte(addressStr), &address)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				startId := v3.ID{Author: sau, Seq: slt}
				endId := v3.ID{Author: fau, Seq: flt}

				mop, err := instance.Rewind(startId, endId, address)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return mop.AsArr()
			}),
			"MergeOp": js.FuncOf(func(this js.Value, args []js.Value) any {
				if len(args) != 1 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 1 argument, got %d", len(args)),
					}
				}
				msgArg := args[0]
				if msgArg.Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %s", msgArg.Type()),
					}
				}
				msgJson := msgArg.String()

				var msg v3.Message
				err := json.Unmarshal([]byte(msgJson), &msg)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				actions, err := instance.MergeOp(msg.Op)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return actions.AsJS()
			}),
			"ApplyDiff": js.FuncOf(func(this js.Value, args []js.Value) any {
				author := args[0].String()
				newText := args[1].String()
				sau := args[2].Index(0).String()
				slt := args[2].Index(1).Int()
				fau := args[3].Index(0).String()
				flt := args[3].Index(1).Int()

				startID := v3.ID{Author: sau, Seq: slt}
				finalID := v3.ID{Author: fau, Seq: flt}

				ops, actions, err := instance.ApplyMarkdownDiff(author, newText, startID, finalID)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("ApplyMarkdownDiff(%q, %q, %v, %v): %v", author, newText, startID, finalID, err),
					}
				}

				return map[string]interface{}{
					"mop":     ops.AsArr(),
					"actions": actions.AsJS(),
				}
			}),
			"Size": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				return instance.VisSize
			}),
			"GetTextAt": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				ix := args[0].Int()
				_, node, err := instance.Rope.GetNode(ix)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				txt := v3.Uint16ToStr(node.Val.Text)
				return map[string]interface{}{
					"text": txt,
				}
			}),
			"GetID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				ix := args[0].Int()
				id, err := instance.Rope.GetVisID(ix)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"GetVisID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				id, err := instance.Rope.GetVisID(args[0].Int())
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"GetFirstID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				id, err := instance.GetFirstID()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"GetLastID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				id, err := instance.GetLastID()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"GetFirstTotID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				id, err := instance.GetFirstTotID()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"GetLastTotID": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				id, err := instance.GetLastTotID()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"GetIndex": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				id := v3.ID{
					Author: args[0].Index(0).String(),
					Seq:    args[0].Index(1).Int(),
				}

				visible, total, err := instance.Rope.GetIndex(id)
				if err != nil {
					return map[string]interface{}{
						"error":   err.Error(),
						"visible": visible,
						"total":   total,
					}
				}

				return map[string]interface{}{
					"visible": visible,
					"total":   total,
				}
			}),
			"VisLeftOf": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				id := v3.ID{
					Author: args[0].Index(0).String(),
					Seq:    args[0].Index(1).Int(),
				}

				leftID, err := instance.VisLeftOf(id)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return leftID.AsJS()
			}),
			"VisRightOf": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				id := v3.ID{
					Author: args[0].Index(0).String(),
					Seq:    args[0].Index(1).Int(),
				}

				rightID, err := instance.VisRightOf(id)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return rightID.AsJS()
			}),
			"TotLeftOf": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				id := v3.ID{
					Author: args[0].Index(0).String(),
					Seq:    args[0].Index(1).Int(),
				}

				leftID, err := instance.TotLeftOf(id)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return leftID.AsJS()
			}),
			"GetCurSpanFormat": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				startID, err := IdArg(args, 0)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("start id: %s", err.Error()),
					}
				}

				endID, err := IdArg(args, 1)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("end id: %s", err.Error()),
					}
				}

				format, err := instance.GetCurSpanFormat(startID, endID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return format.AsMap()
			}),
			"GetCurLineFormat": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				startID, err := IdArg(args, 0)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("start id: %s", err.Error()),
					}
				}

				endID, err := IdArg(args, 1)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("end id: %s", err.Error()),
					}
				}

				format, err := instance.GetCurLineFormat(startID, endID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return format.AsMap()
			}),
			"GetSelection": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 2 && len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2 or 3 arguments, got %d", len(args)),
					}
				}

				startID, err := IdArg(args, 0)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("start id: %s", err.Error()),
					}
				}

				endID, err := IdArg(args, 1)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("end id: %s", err.Error()),
					}
				}

				isScrub := args[2].Bool()

				var address *v3.ContentAddress
				if isScrub && instance.ScrubState != nil {
					address = instance.ScrubState.CurAddress
				}

				selection, err := instance.GetSelection(startID, endID, address, true)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return selection.AsJS()
			}),
			"GetSelectionAt": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()

				if len(args) != 2 && len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2 or 3 arguments, got %d", len(args)),
					}
				}

				startID, err := IdArg(args, 0)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("start id: %s", err.Error()),
					}
				}

				endID, err := IdArg(args, 1)
				if err != nil {
					return map[string]interface{}{
						"error": fmt.Sprintf("end id: %s", err.Error()),
					}
				}

				var address *v3.ContentAddress = nil

				if len(args) == 3 {
					address, err = ContentAddressArg(args, 2)
					if err != nil {
						return map[string]interface{}{
							"error": fmt.Sprintf("content address: %s", err.Error()),
						}
					}
				}

				selection, err := instance.GetSelection(startID, endID, address, true)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return selection.AsJS()
			}),
			"IDFromIDAndOffset": js.FuncOf(func(this js.Value, args []js.Value) any {
				if len(args) != 2 && len(args) != 3 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2 or 3 arguments, got %d", len(args)),
					}
				}

				if args[0].Type() != js.TypeObject {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected object, got %T", args[1]),
					}
				}

				if args[1].Type() != js.TypeNumber {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected number, got %T", args[1]),
					}
				}

				sau := args[0].Index(0).String()
				slt := args[0].Index(1).Int()
				targetID := v3.ID{Author: sau, Seq: slt}

				offset := args[1].Int()

				var err error
				var ca *v3.ContentAddress = nil

				if instance.ScrubState != nil {
					ca = instance.ScrubState.CurAddress
				} else if len(args) == 3 {
					ca, err = ContentAddressArg(args, 2)
					if err != nil {
						return map[string]interface{}{
							"error": fmt.Sprintf("content address: %s", err.Error()),
						}
					}
				}

				id, err := instance.IDFromIDAndOffset(targetID, offset, ca)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return id.AsJS()
			}),
			"RenderOp": js.FuncOf(func(this js.Value, args []js.Value) any {
				if len(args) != 1 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 1 argument, got %d", len(args)),
					}
				}
				msgArg := args[0]
				if msgArg.Type() != js.TypeString {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected string, got %s", msgArg.Type()),
					}
				}
				msgJson := msgArg.String()

				var msg v3.Message
				err := json.Unmarshal([]byte(msgJson), &msg)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				span, err := instance.RenderOp(msg.Op)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return map[string]interface{}{
					"firstBlockID": span.FirstBlockID.AsJS(),
					"lastBlockID":  span.LastBlockID.AsJS(),
					"toStartID":    span.ToStartID.AsJS(),
					"toEndID":      span.ToEndID.AsJS(),
					"html":         span.Html,
				}
			}),
			"ScrubInit": js.FuncOf(func(this js.Value, args []js.Value) any {
				if len(args) > 2 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 2 or fewer arguments, got %d", len(args)),
					}
				}

				var startID, endID *v3.ID

				if len(args) >= 1 {
					if args[0].Type() == js.TypeObject {
						sau := args[0].Index(0).String()
						slt := args[0].Index(1).Int()
						startID = &v3.ID{Author: sau, Seq: slt}
					} else {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected object, got %T", args[0]),
						}
					}
				}

				if len(args) >= 2 {
					if args[1].Type() == js.TypeObject {
						sau := args[1].Index(0).String()
						slt := args[1].Index(1).Int()
						endID = &v3.ID{Author: sau, Seq: slt}
					} else {
						return map[string]interface{}{
							"error": fmt.Sprintf("expected object, got %T", args[1]),
						}
					}
				}

				n, err := instance.ScrubInit(startID, endID)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return n
			}),
			"ScrubTo": js.FuncOf(func(this js.Value, args []js.Value) any {
				// acepts a single int as the index to rewind to
				if len(args) != 1 {
					return map[string]interface{}{
						"error": fmt.Sprintf("expected 1 argument, got %d", len(args)),
					}
				}

				ix := args[0].Int()

				rewindStep, err := instance.ScrubTo(ix)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				if rewindStep == nil {
					return false
				}

				return map[string]interface{}{
					"firstBlockID":  rewindStep.Span.FirstBlockID.AsJS(),
					"lastBlockID":   rewindStep.Span.LastBlockID.AsJS(),
					"cursorStartID": rewindStep.CursorStartID.AsJS(),
					"cursorEndID":   rewindStep.CursorEndID.AsJS(),
					"html":          rewindStep.Html,
				}
			}),
			"ScrubRevert": js.FuncOf(func(this js.Value, args []js.Value) any {
				if instance.ScrubState == nil {
					return map[string]interface{}{
						"error": "no rewind state",
					}
				}

				mop, err := instance.Rewind(
					instance.ScrubState.StartID,
					instance.ScrubState.EndID,
					*instance.ScrubState.CurAddress,
				)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				cursorID, err := instance.NearestAt(instance.ScrubState.EndID, nil)
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				instance.ScrubState = nil

				return map[string]interface{}{
					"mop":      mop.AsArr(),
					"cursorID": cursorID.AsJS(),
				}
			}),
			"ScrubMax": js.FuncOf(func(this js.Value, args []js.Value) any {
				if instance.ScrubState == nil {
					return 0
				}

				if instance.ScrubState.IdTree == nil {
					return instance.OpIndex.Size() - 1
				}

				return instance.ScrubState.IdTree.Size - 1
			}),
			"ScrubExit": js.FuncOf(func(this js.Value, args []js.Value) any {
				instance.ScrubState = nil
				return nil
			}),
			"CanRedo": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				return instance.CanRedo()
			}),
			"CanUndo": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				return instance.CanUndo()
			}),
			"DocStats": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				stats, err := instance.DocStats()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return stats.AsJS()
			}),
			"OpStats": js.FuncOf(func(this js.Value, args []js.Value) any {
				defer catchPanic()
				stats, err := instance.OpStats()
				if err != nil {
					return map[string]interface{}{
						"error": err.Error(),
					}
				}

				return stats.AsJS()
			}),
			"Exit": js.FuncOf(func(this js.Value, args []js.Value) any {
				Exit()
				return nil
			}),
		}
	}))

	<-done // Wait indefinitely
}

func Exit() {
	done <- struct{}{}
}

func IdArg(args []js.Value, ix int) (v3.ID, error) {
	if len(args) <= ix {
		return v3.ID{}, fmt.Errorf("expected at least %d arguments, got %d", ix+1, len(args))
	}
	if args[ix].Type() != js.TypeObject {
		return v3.ID{}, fmt.Errorf("expected %dth argument to an id(object), got %s", ix, args[ix].Type())
	}

	idArgs := args[ix]
	if idArgs.Index(0).Type() != js.TypeString || idArgs.Index(1).Type() != js.TypeNumber {
		return v3.ID{}, fmt.Errorf("expected %dth argument to an id [string, number], got [%s, %s]", ix, idArgs.Index(0).Type(), idArgs.Index(1).Type())
	}

	au := args[ix].Index(0).String()
	sq := args[ix].Index(1).Int()

	return v3.ID{Author: au, Seq: sq}, nil
}

func ContentAddressArg(args []js.Value, ix int) (*v3.ContentAddress, error) {
	if len(args) <= ix {
		return nil, fmt.Errorf("expected at least %d arguments, got %d", ix+1, len(args))
	}

	if args[ix].IsNull() {
		return nil, nil
	}

	if args[ix].Type() != js.TypeString {
		return nil, fmt.Errorf("expected %dth argument to an contentAddress(string), got %s", ix, args[ix].Type())
	}

	address := args[ix].String()

	ca := v3.ContentAddress{}
	err := json.Unmarshal([]byte(address), &ca)
	if err != nil {
		return nil, err
	}

	return &ca, nil
}
