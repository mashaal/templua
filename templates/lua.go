package templates

import (
	"fmt"
	"github.com/yuin/gopher-lua"
	"log"
	"strings"
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

func (lt *LuaTemplate) renderCustomElement(name string, template *lua.LTable) string {
	var styles string
	var contentParts []string

	// Get styles
	if stylesTbl := template.RawGetString("styles"); stylesTbl != lua.LNil {
		styles = stylesTbl.String()
	}

	// Get content
	if contentTbl := template.RawGetString("content"); contentTbl != lua.LNil {
		if content, ok := contentTbl.(*lua.LTable); ok {
			for i := 1; i <= content.Len(); i++ {
				if item := content.RawGetInt(i); item != lua.LNil {
					contentParts = append(contentParts, item.String())
				}
			}
		}
	}

	// Build the custom element HTML
	result := fmt.Sprintf(`<%s>
  <template shadowrootmode="open">
    <style>%s</style>
    %s
  </template>
</%s>`, name, styles, strings.Join(contentParts, "\n    "), name)

	return result
}

func (lt *LuaTemplate) registerCustomElement(name string, loader string) {
	lt.L.SetGlobal(name, lt.L.NewFunction(func(L *lua.LState) int {
		// Load the component module
		if err := L.DoFile(loader); err != nil {
			log.Printf("Error loading custom element %s: %v", name, err)
			L.Push(lua.LString(""))
			return 1
		}

		// Get the component function
		component := L.Get(-1)
		L.Pop(1)

		// Call the component function with props
		var props *lua.LTable
		if L.GetTop() > 0 && L.Get(1).Type() == lua.LTTable {
			props = L.ToTable(1)
		} else {
			props = L.NewTable()
		}

		err := L.CallByParam(lua.P{
			Fn:      component,
			NRet:    1,
			Protect: true,
		}, props)

		if err != nil {
			log.Printf("Error calling custom element %s: %v", name, err)
			L.Push(lua.LString(""))
			return 1
		}

		// Get the result
		result := L.Get(-1)
		L.Pop(1)

		// If result is a table with name and template
		if resultTbl, ok := result.(*lua.LTable); ok {
			if name := resultTbl.RawGetString("name"); name != lua.LNil {
				if template := resultTbl.RawGetString("template"); template != lua.LNil {
					if templateTbl, ok := template.(*lua.LTable); ok {
						html := lt.renderCustomElement(name.String(), templateTbl)
						L.Push(lua.LString(html))
						return 1
					}
				}
			}
		}

		log.Printf("Invalid result from custom element %s", name)
		L.Push(lua.LString(""))
		return 1
	}))
}

func (lt *LuaTemplate) registerComponents() {
	// Register HTML elements
	elements := []string{
		"Html", "Head", "Body", "Title",
		"Div", "P", "Span", "A", "Img",
		"H1", "H2", "H3", "H4", "H5", "H6",
		"Ul", "Ol", "Li", "Table", "Tr", "Td",
		"Form", "Input", "Button", "Label",
		"Script", "Style", "Link", "Template",
		"Meta", "Img", "Input", "Br", "Hr",
	}

	// Register all elements
	for _, name := range elements {
		isSelfClosing := false
		switch name {
		case "Meta", "Img", "Input", "Br", "Hr":
			isSelfClosing = true
		}
		lt.registerElement(name, isSelfClosing)
	}

	// Register Style helper
	lt.L.SetGlobal("Style", lt.L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() < 1 {
			return 0
		}

		// If argument is a string, use it directly
		if L.Get(1).Type() == lua.LTString {
			L.Push(L.Get(1))
			return 1
		}

		// Otherwise, process as a table
		var styles []string
		tbl := L.ToTable(1)
		tbl.ForEach(func(k, v lua.LValue) {
			if k.Type() == lua.LTString {
				selector := k.String()
				rules := v.String()
				styles = append(styles, fmt.Sprintf("%s { %s }", selector, rules))
			}
		})

		L.Push(lua.LString(strings.Join(styles, "\n")))
		return 1
	}))

	// Register custom elements
	lt.registerCustomElement("Card", "templates/components/card.lua")
}

func (lt *LuaTemplate) registerElement(name string, selfClosing bool) {
	lt.L.SetGlobal(name, lt.L.NewFunction(func(L *lua.LState) int {
		log.Printf("=== Creating %s element ===", name)

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

		// Special handling for Script with direct string argument
		if name == "Script" && L.GetTop() == 1 && L.Get(1).Type() == lua.LTString {
			script := L.ToString(1)
			result := fmt.Sprintf("<script>%s</script>", script)
			L.Push(lua.LString(result))
			return 1
		}

		// Get the first argument
		var attrs *lua.LTable
		var content []string

		// If first argument is a table, it could be attributes or content
		if L.Get(1).Type() == lua.LTTable {
			tbl := L.ToTable(1)

			// Check if it's an attribute table (has string keys)
			hasStringKeys := false
			tbl.ForEach(func(key, _ lua.LValue) {
				if key.Type() == lua.LTString {
					hasStringKeys = true
				}
			})

			if hasStringKeys {
				attrs = tbl
				// If there's a second argument, it's content
				if L.GetTop() > 1 {
					content = append(content, lua.LVAsString(L.Get(2)))
				}
			} else {
				// It's a content table
				for i := 1; i <= tbl.Len(); i++ {
					content = append(content, lua.LVAsString(tbl.RawGetInt(i)))
				}
			}
		} else {
			// Single argument is content
			content = append(content, lua.LVAsString(L.Get(1)))
		}

		// Build attributes string
		var attrList []string
		if attrs != nil {
			attrs.ForEach(func(key, value lua.LValue) {
				if k, ok := key.(lua.LString); ok {
					attrList = append(attrList, fmt.Sprintf(`%s="%s"`, k, value.String()))
				}
			})
		}

		// Join attributes
		attrStr := ""
		if len(attrList) > 0 {
			attrStr = " " + strings.Join(attrList, " ")
		}

		// Join content
		contentStr := strings.Join(content, "")

		// Build final HTML
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

func (lt *LuaTemplate) pushString(L *lua.LState, str string) int {
	L.Push(lua.LString(str))
	return 1
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
