local function createWatcher()
    local file = io.open("framework/watcher.js", "r")
    if not file then
        error("Could not open watcher.js")
    end
    local content = file:read("*all")
    file:close()
    return content
end

return createWatcher
