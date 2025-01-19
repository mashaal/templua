local function render(params)
    params = params or {}  -- Initialize empty table if no params passed
    local heading = params.heading or "Welcome to Templua!"
    
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