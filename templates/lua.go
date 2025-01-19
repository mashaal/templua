package templates

import (
	"fmt"
	"log"
	"strings"
	lua "github.com/yuin/gopher-lua"
)

type LuaTemplate struct {
	L *lua.LState
}

func NewLuaTemplate() *LuaTemplate {
	L := lua.NewState()
	lt := &LuaTemplate{L: L}
	lt.registerComponents()
	return lt
}

func (lt *LuaTemplate) Close() {
	lt.L.Close()
}

func (lt *LuaTemplate) registerComponents() {
	// Document structure
	rootElements := []string{"Html", "Head", "Body"}
	sectionElements := []string{"Header", "Nav", "Main", "Article", "Section", "Aside", "Footer", "H1", "H2", "H3", "H4", "H5", "H6"}
	textElements := []string{"Div", "P", "Pre", "Blockquote", "Ol", "Ul", "Li", "Dl", "Dt", "Dd", "Figure", "Figcaption"}
	inlineElements := []string{"A", "Em", "Strong", "Small", "S", "Cite", "Q", "Dfn", "Abbr", "Code", "Var", "Samp", "Kbd", "Sub", "Sup", "I", "B", "U", "Mark", "Ruby", "Rt", "Rp", "Bdi", "Bdo", "Span"}
	selfClosingElements := []string{"Br", "Hr", "Wbr", "Img", "Source", "Track", "Embed", "Col", "Input"}

	// Register all elements
	for _, name := range rootElements {
		lt.registerElement(name, false)
	}
	for _, name := range sectionElements {
		lt.registerElement(name, false)
	}
	for _, name := range textElements {
		lt.registerElement(name, false)
	}
	for _, name := range inlineElements {
		lt.registerElement(name, false)
	}
	for _, name := range selfClosingElements {
		lt.registerElement(name, true)
	}

	// Special components
	lt.L.SetGlobal("Meta", lt.L.NewFunction(func(L *lua.LState) int {
		log.Printf("=== Creating Meta element ===")
		attrs := L.CheckTable(1)
		var attrList []string
		attrs.ForEach(func(key, value lua.LValue) {
			if k, ok := key.(lua.LString); ok {
				log.Printf("[Meta] Found attribute %v = %v", k, value)
				attrList = append(attrList, fmt.Sprintf(`%s="%s"`, k, value.String()))
			}
		})
		attrStr := ""
		if len(attrList) > 0 {
			attrStr = " " + strings.Join(attrList, " ")
		}
		result := fmt.Sprintf("<meta%s>", attrStr)
		log.Printf("[Meta] Final rendered element: %s", result)
		L.Push(lua.LString(result))
		return 1
	}))
}

func (lt *LuaTemplate) registerElement(name string, selfClosing bool) {
	lt.L.SetGlobal(name, lt.L.NewFunction(func(L *lua.LState) int {
		log.Printf("=== Creating element: %s ===", name)

		// Handle the case where no arguments are provided
		if L.GetTop() == 0 {
			log.Printf("[%s] No arguments provided", name)
			if selfClosing {
				result := fmt.Sprintf("<%s>", strings.ToLower(name))
				log.Printf("[%s] Empty self-closing element: %s", name, result)
				L.Push(lua.LString(result))
			} else {
				result := fmt.Sprintf("<%s></%s>", strings.ToLower(name), strings.ToLower(name))
				log.Printf("[%s] Empty element: %s", name, result)
				L.Push(lua.LString(result))
			}
			return 1
		}

		// Get the first argument which should be a table
		tbl := L.CheckTable(1)
		log.Printf("[%s] Table has %d elements", name, tbl.Len())

		// First, collect all content in order
		var orderedContent []string
		for i := 1; i <= tbl.Len(); i++ {
			value := tbl.RawGetInt(i)
			log.Printf("[%s] Processing ordered value %d: %v (type: %T)", name, i, value, value)
			switch v := value.(type) {
			case lua.LString:
				orderedContent = append(orderedContent, string(v))
			case lua.LNumber:
				orderedContent = append(orderedContent, fmt.Sprintf("%v", float64(v)))
			case lua.LBool:
				orderedContent = append(orderedContent, fmt.Sprintf("%v", bool(v)))
			case *lua.LTable:
				// For nested elements, we expect them to return a string
				if str := lua.LVAsString(value); str != "" {
					orderedContent = append(orderedContent, str)
				}
			default:
				log.Printf("[%s] Warning: unexpected value type %T", name, v)
			}
		}

		// Process attributes (string keys)
		var attrList []string
		tbl.ForEach(func(key, value lua.LValue) {
			if k, ok := key.(lua.LString); ok {
				if s := k.String(); s != "1" && !strings.HasPrefix(s, "_") {
					log.Printf("[%s] Found attribute %v = %v", name, k, value)
					attrList = append(attrList, fmt.Sprintf(`%s="%s"`, k, value.String()))
				}
			}
		})

		// Build the attribute string
		attrStr := ""
		if len(attrList) > 0 {
			attrStr = " " + strings.Join(attrList, " ")
		}
		log.Printf("[%s] Final attributes: %s", name, attrStr)

		// Build the final HTML
		contentStr := strings.Join(orderedContent, "")
		log.Printf("[%s] Content to be inserted: %s", name, contentStr)

		var result string
		tagName := strings.ToLower(name)
		if selfClosing {
			result = fmt.Sprintf("<%s%s>", tagName, attrStr)
		} else {
			if name == "Html" {
				result = fmt.Sprintf("<!DOCTYPE html><%s%s>%s</%s>", tagName, attrStr, contentStr, tagName)
			} else {
				result = fmt.Sprintf("<%s%s>%s</%s>", tagName, attrStr, contentStr, tagName)
			}
		}

		log.Printf("[%s] Final rendered element: %s", name, result)
		L.Push(lua.LString(result))
		return 1
	}))
}

func (lt *LuaTemplate) RenderHTML(script string) (string, error) {
	log.Printf("Starting template rendering with script:\n%s", script)

	if err := lt.L.DoString(script); err != nil {
		log.Printf("Template execution error: %v", err)
		return "", fmt.Errorf("template execution error: %v", err)
	}

	// Get the result from the stack
	result := lt.L.Get(-1)
	lt.L.Pop(1)

	if result == lua.LNil {
		log.Printf("Template did not produce any output")
		return "", fmt.Errorf("template did not produce any output")
	}

	html := result.String()
	log.Printf("Generated HTML:\n%s", html)
	return html, nil
}

func (lt *LuaTemplate) RenderHTMLWithVars(script string, vars map[string]interface{}) (string, error) {
	if err := lt.L.DoString(script); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	// Get the render function result from the stack
	result := lt.L.Get(-1)
	lt.L.Pop(1)

	if result.Type() != lua.LTFunction {
		return "", fmt.Errorf("template must return a function")
	}

	// Create Lua table from vars
	paramsTable := lt.L.NewTable()
	for k, v := range vars {
		switch val := v.(type) {
		case string:
			paramsTable.RawSetString(k, lua.LString(val))
		case int:
			paramsTable.RawSetString(k, lua.LNumber(val))
		case float64:
			paramsTable.RawSetString(k, lua.LNumber(val))
		case bool:
			paramsTable.RawSetString(k, lua.LBool(val))
		default:
			log.Printf("Unsupported type for key %s: %T", k, v)
		}
	}

	// Call the render function with params
	err := lt.L.CallByParam(lua.P{
		Fn:      result,
		NRet:    1,
		Protect: true,
	}, paramsTable)
	if err != nil {
		return "", fmt.Errorf("failed to call template function: %v", err)
	}

	// Get the rendered HTML
	rendered := lt.L.Get(-1)
	lt.L.Pop(1)

	return rendered.String(), nil
}
