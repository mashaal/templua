local function Card(props)
    props = props or {}
    local title = props.title or "Card Title"
    local content = props.content or "Card content goes here"

    -- Define card styles
    local styles = [[
        :host {
            display: block;
            border: 1px solid #ddd;
            border-radius: 8px;
            padding: 16px;
            margin: 16px;
        }
        h2 {
            margin: 0;
            color: #333;
            font-size: 1.5em;
        }
        p {
            margin-top: 8px;
            color: #666;
        }
    ]]

    -- Return component definition
    return {
        name = "card-component", -- Custom element name
        template = {
            styles = styles,
            content = {
                H2({title}),
                P({content})
            }
        }
    }
end

return Card
