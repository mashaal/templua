# Templua

A minimalist templating engine that uses Lua for templating HTML, inspired by the simplicity of [100 Rabbits](https://100r.co/site/home.html) projects.

## Features

- **Lua-based Templates**: Write HTML templates using Lua's clean and expressive syntax
- **Live Reload**: Automatic page refresh during development when templates are modified
- **Minimal Dependencies**: Built with simplicity in mind, using only essential dependencies
- **HTML Element Functions**: Create HTML elements using simple Lua functions
- **Variable Support**: Pass variables from Go to your Lua templates

## Installation

```bash
go get github.com/yourusername/templua
```

## Dependencies

- [gopher-lua](https://github.com/yuin/gopher-lua): Lua VM implementation in Go
- [echo](https://github.com/labstack/echo): High performance web framework
- [gorilla/websocket](https://github.com/gorilla/websocket): WebSocket implementation for Go
- [fsnotify](https://github.com/fsnotify/fsnotify): Cross-platform file system notifications

## Usage

### Basic Template

```lua
-- templates/home.lua
local function render(params)
    params = params or {}
    local heading = params.heading or "Welcome"
    
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

### Go Server

```go
package main

import (
    "github.com/labstack/echo/v4"
    "templua/templates"
)

func main() {
    e := echo.New()
    lt := templates.NewLuaTemplate()
    defer lt.Close()

    e.GET("/", func(c echo.Context) error {
        vars := map[string]interface{}{
            "heading": "Welcome to Templua",
        }
        
        html, err := lt.RenderHTMLWithVars("templates/home.lua", vars)
        if err != nil {
            return err
        }
        
        return c.HTML(http.StatusOK, html)
    })

    e.Logger.Fatal(e.Start(":1323"))
}
```

## Development Mode

During development, Templua provides live reload functionality that automatically refreshes your browser when template files are modified. This feature is enabled by default when running the server locally.

### How Live Reload Works

1. A WebSocket connection is established between the browser and server
2. The server watches the templates directory for changes
3. When a template file is modified, all connected clients are notified
4. The browser automatically refreshes to show the updated content

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.