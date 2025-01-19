# Templua

Templua is a lightweight templating engine that bridges Go and Lua, designed for building elegant HTML templates with the full power of a programming language. Built on the Echo framework, it provides both performance and simplicity.

## Core Features

- **Programmatic Templates**: Write templates using Lua's clean and expressive syntax
- **Seamless Integration**: Direct data flow between Go and Lua
- **Performance**: Built on Echo's high-performance HTTP framework
- **Type Safety**: Robust type checking between Go and Lua
- **Simplicity**: Clear, predictable API design
- **Flexibility**: Full access to Lua's programming capabilities

## Getting Started

Initialize the server:

```bash
go run ./cmd/templua
```

The server will be available at `http://localhost:1323`.

## Template Design

A template in Templua (`templates/home.lua`):

```lua
local function render(params)
    params = params or {}
    local heading = params.heading or "Welcome to Templua"
    
    return Html {
        Head {
            Meta { charset="utf-8" },
            Meta { name="viewport", content="width-device-width, initial-scale=1" }
        },
        Body {
            H1 { heading },
            P { "This is a test page." }
        }
    }
end

return render
```

Integration in Go (`cmd/templua/templua.go`):

```go
vars := map[string]interface{}{
    "heading": "Welcome to Templua",
}

html, err := lt.RenderHTMLWithVars(template, vars)
if err != nil {
    log.Printf("Failed to render template: %v", err)
    return fmt.Errorf("failed to render template: %v", err)
}
```

## Type System

Templua supports a core set of data types for Go-Lua communication:
- `string`: Text values
- `int`: Integer numbers
- `float64`: Floating-point numbers
- `bool`: Boolean values

## Architecture

```
templua/
├── cmd/
│   └── templua/
│       └── templua.go    # Application entry point
├── templates/
│   ├── lua.go           # Template engine implementation
│   └── home.lua         # Template definition
├── go.mod               # Module definition
└── README.md           # Documentation
```

## Contributing

Contributions are welcome. Please ensure your changes maintain the project's focus on simplicity and reliability.
