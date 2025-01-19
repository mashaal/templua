local function render(params)
    params = params or {}
    local heading = params.heading or "Welcome"
    
    -- Load the watcher script
    local createWatcher = require("templates.watcher")
    local watcherScript = createWatcher()
    
    return Html {
        Head {
            Meta { charset="utf-8" },
            Meta { name="viewport", content="width-device-width, initial-scale=1" },
            Script { watcherScript }
        },
        Body {
            H1 { heading },
            P { "This is a test page." }
        }
    }
end

return render