# 🔥 Templua

Tired of clunky template engines? Meet Templua - where Go meets Lua for HTML templating that doesn't suck. Built on Echo because life's too short for slow servers.

## 🚀 Why Templua?

- **Lua-Powered Templates**: Write HTML like you're coding, not like you're stuck in 1995
- **Go-Lua Data Bridge**: Sling data between Go and Lua like a boss
- **Blazing Fast**: Echo framework under the hood = speed demon
- **Dead Simple API**: Because nobody has time for complexity
- **Type-Safe**: Keep your data clean and your errors at compile time
- **Zero BS**: Just the features you need, none you don't

## 🎮 Quick Start

Fire it up:

```bash
go run ./cmd/templua
```

Point your browser to `http://localhost:1323` and watch the magic happen.

## 💻 Show Me The Code

Here's what a Lua template looks like (`templates/home.lua`):

```lua
local function render(params)
    params = params or {}  -- Fail-safe mode: engaged
    local heading = params.heading or "Welcome to Templuta!"
    
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

Hook it up in Go (`cmd/templua/templua.go`):

```go
vars := map[string]interface{}{
    "heading": "Welcome to Dynamic Templua!",
}

html, err := lt.RenderHTMLWithVars(template, vars)
if err != nil {
    log.Printf("Failed to render template: %v", err)
    return fmt.Errorf("failed to render template: %v", err)
}
```

## 🎯 Supported Types

Templua speaks your language:
- `string`: For your text needs
- `int`: Count it up
- `float64`: Keep it precise
- `bool`: True that!

## 📁 Project Structure

```
templua/
├── cmd/
│   └── templua/
│       └── templua.go    # Where the action starts
├── templates/
│   ├── lua.go           # The engine room
│   └── home.lua         # Your first template
├── go.mod               # Keep it tidy
└── README.md           # You are here
```

## 🤘 Contributing

Got ideas? PRs welcome. Keep it clean, keep it mean, keep it working.
