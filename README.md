# TEMPLUA

## Purpose

Templua is a template engine that uses Lua as its expression language. It is designed to be simple, predictable, and maintainable.

## Method

Templates are Lua functions that return HTML structure. Data flows from Go to Lua through a minimal interface. The design emphasizes readability over convenience.

## Usage

Start:
```sh
go run ./cmd/templua
```

Write (`templates/home.lua`):
```lua
local function render(params)
    params = params or {}
    local heading = params.heading or "Welcome"
    
    return Html {
        Head {
            Meta { charset="utf-8" },
            Meta { name="viewport", content="width-device-width, initial-scale=1" }
        },
        Body {
            H1 { heading }
        }
    }
end

return render
```

Use (`cmd/templua/templua.go`):
```go
vars := map[string]interface{}{
    "heading": "Welcome"
}

html, err := lt.RenderHTMLWithVars(template, vars)
```

## Types

- string
- int
- float64
- bool

## Structure

```
templua/
├── cmd/
│   └── templua/
│       └── templua.go    # Application entry point
├── templates/
│   ├── lua.go           # Template engine implementation
│   └── home.lua         # Template definition
└── go.mod               # Module definition
```