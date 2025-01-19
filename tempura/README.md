# Tempura

A minimal HTML templating engine using Lua. Write declarative HTML with the simplicity of Lua tables.

## Philosophy

Tempura believes in simplicity. No complex templating syntax, just Lua tables. No magic, just functions. Write HTML the way you think about it.

## Usage

```lua
-- Elements are just functions
Html({}, {
  Head({}, {
    Title({}, {"My Page"})
  }),
  Body({}, {
    H1({class = "title"}, {"Hello"}),
    P({"Write content directly"}),
    Div({id = "main"}, {
      "Mix ", Em({"and match"}), " content"
    })
  })
})
```

## Custom Elements

Create web components with shadow DOM, simply:

```lua
local function Card(props)
  return {
    name = "card-element",
    template = {
      styles = [[
        :host { display: block }
        h2 { color: #333 }
      ]],
      content = {
        H2({props.title}),
        P({props.content})
      }
    }
  }
end
```

## Install

```bash
go get github.com/mashaal/tempura
```

## Go

```go
lt := tempura.NewLuaTemplate()
defer lt.Close()

// Register a custom component
lt.RegisterCustomComponent("Card", "components/card.lua")

// Render template with variables
html, err := lt.RenderHTMLWithVars(template, map[string]interface{}{
  "title": "Welcome"
})
```

