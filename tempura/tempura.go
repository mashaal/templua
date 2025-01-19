package tempura

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

func (lt *LuaTemplate) RegisterCustomComponent(name string, loader string) {
	lt.registerCustomElement(name, loader)
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
}

func (lt *LuaTemplate) registerElement(name string, selfClosing bool) {
	lt.L.SetGlobal(name, lt.L.NewFunction(func(L *lua.LState) int {
		log.Printf("=== Creating %s element ===", name)

		// Handle the case where no arguments are provided
		if L.GetTop() == 0 {
			log.Printf("[%s] No arguments provided", name)
			if selfClosing {
				result := fmt.Sprintf("<%s/>", strings.ToLower(name))
				log.Printf("[%s] Empty self-closing element: %s", name, result)
				L.Push(lua.LString(result))
			} else {
				result := fmt.Sprintf("<%s></%s>", strings.ToLower(name), strings.ToLower(name))
				log.Printf("[%s] Empty element: %s", name, result)
				L.Push(lua.LString(result))
			}
			return 1
		}

		// Get arguments
		var attrs *lua.LTable
		var content []string

		// Check first argument
		firstArg := L.Get(1)
		if firstArg.Type() == lua.LTTable {
			// If there's a second argument, first is attrs, second is content
			if L.GetTop() > 1 && L.Get(2).Type() == lua.LTTable {
				attrs = L.ToTable(1)
				contentTable := L.ToTable(2)
				contentTable.ForEach(func(_, v lua.LValue) {
					if v != lua.LNil {
						content = append(content, v.String())
					}
				})
			} else {
				// If only one table argument, check if it has content-like values or attribute-like values
				tbl := L.ToTable(1)
				hasAttrs := false
				tbl.ForEach(func(k, _ lua.LValue) {
					if k.Type() == lua.LTString {
						hasAttrs = true
					}
				})
				if hasAttrs {
					attrs = tbl
				} else {
					// Treat as content
					tbl.ForEach(func(_, v lua.LValue) {
						if v != lua.LNil {
							content = append(content, v.String())
						}
					})
				}
			}
		}

		// Build the element
		var result strings.Builder
		result.WriteString("<")
		result.WriteString(strings.ToLower(name))

		// Add attributes if present
		if attrs != nil {
			attrs.ForEach(func(k, v lua.LValue) {
				if k.Type() == lua.LTString {
					result.WriteString(fmt.Sprintf(" %s=\"%s\"", k.String(), v.String()))
				}
			})
		}

		if selfClosing {
			result.WriteString("/>")
		} else {
			result.WriteString(">")
			result.WriteString(strings.Join(content, ""))
			result.WriteString("</")
			result.WriteString(strings.ToLower(name))
			result.WriteString(">")
		}

		L.Push(lua.LString(result.String()))
		return 1
	}))
}

func (lt *LuaTemplate) tableToString(tbl *lua.LTable) string {
	var parts []string
	tbl.ForEach(func(_, v lua.LValue) {
		if v != lua.LNil {
			parts = append(parts, v.String())
		}
	})
	return strings.Join(parts, "")
}

func (lt *LuaTemplate) pushString(L *lua.LState, str string) int {
	L.Push(lua.LString(str))
	return 1
}

// RenderHTML renders a Lua template string and returns the resulting HTML
func (lt *LuaTemplate) RenderHTML(script string) (string, error) {
	if err := lt.L.DoString(script); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	// Get the result from the Lua stack
	result := lt.L.Get(-1)
	lt.L.Pop(1)

	if result == lua.LNil {
		return "", fmt.Errorf("template returned nil")
	}

	// If result is a function, it's an error because we expect the template to return HTML directly
	if result.Type() == lua.LTFunction {
		return "", fmt.Errorf("template returned a function instead of HTML string")
	}

	return result.String(), nil
}

// RenderHTMLWithVars renders a Lua template string with variables and returns the resulting HTML
func (lt *LuaTemplate) RenderHTMLWithVars(script string, vars map[string]interface{}) (string, error) {
	// Create a new table for variables
	varsTable := lt.L.NewTable()
	for k, v := range vars {
		switch val := v.(type) {
		case string:
			varsTable.RawSetString(k, lua.LString(val))
		case int:
			varsTable.RawSetString(k, lua.LNumber(val))
		case float64:
			varsTable.RawSetString(k, lua.LNumber(val))
		case bool:
			varsTable.RawSetString(k, lua.LBool(val))
		case []interface{}:
			tbl := lt.L.NewTable()
			for i, item := range val {
				switch v := item.(type) {
				case string:
					tbl.RawSetInt(i+1, lua.LString(v))
				case int:
					tbl.RawSetInt(i+1, lua.LNumber(v))
				case float64:
					tbl.RawSetInt(i+1, lua.LNumber(v))
				case bool:
					tbl.RawSetInt(i+1, lua.LBool(v))
				}
			}
			varsTable.RawSetString(k, tbl)
		}
	}

	// Set the variables table as a global in Lua
	lt.L.SetGlobal("vars", varsTable)

	return lt.RenderHTML(script)
}
