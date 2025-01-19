local function render(params)
    params = params or {}
    local heading = params.heading or "Welcome"
    
    -- Load the watcher script
    local createWatcher = require("framework.watcher")
    local watcherScript = createWatcher()
    
    return Html({}, {
        Head({}, {
            Meta({charset="utf-8"}),
            Meta({name="viewport", content="width=device-width, initial-scale=1"}),
            Script({}, watcherScript)
        }),
        Body({}, {
            H1({heading}),
            P({"This is a test page with custom elements and shadow DOM."}),
            Card({title = "Card 1", content = "This is the content of Card 1."}),
            Card({title = "Card 2", content = "This is the content of Card 2."})
        })
    })
end

return render(vars)