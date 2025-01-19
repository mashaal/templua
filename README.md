# TEMPLUA

Write web components in Lua, render them with Go.

## USAGE

`go run cmd/templua`

## COMPONENTS

Define a component:
```lua
-- components/card.lua
return {
    name = "my-card",
    template = {
        styles = ":host { display:block }",
        content = { H2 { props.title } }
    }
}
```

Use in Lua:
```lua
-- home.lua
Card { title = "Hello" }
```

Use in Go:
```go
lt := templates.NewLuaTemplate()
defer lt.Close()

// Pass props to template
props := map[string]interface{}{
    "title": "Hello from Go",
}

html := lt.RenderWithProps("templates/home.lua", props)
