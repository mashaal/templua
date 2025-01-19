local function createWatcher()
    return [[
        console.log('Initializing WebSocket connection...');
        
        function connectWebSocket() {
            if (window.location.hostname !== 'localhost') {
                console.log('Not localhost, skipping WebSocket connection');
                return;
            }

            const ws = new WebSocket('ws://' + window.location.host + '/ws');
            
            ws.onopen = function() {
                console.log('WebSocket connection established');
            };
            
            ws.onmessage = function(evt) {
                console.log('Received message:', evt.data);
                if (evt.data === 'reload') {
                    console.log('Reloading page...');
                    window.location.reload();
                }
            };
            
            ws.onclose = function(e) {
                console.log('WebSocket connection closed:', e.reason);
                // Try to reconnect after a delay
                setTimeout(connectWebSocket, 1000);
            };
            
            ws.onerror = function(error) {
                console.error('WebSocket error:', error);
            };

            // Handle page visibility changes
            document.addEventListener('visibilitychange', function() {
                if (document.hidden) {
                    ws.close();
                } else {
                    connectWebSocket();
                }
            });
        }

        // Start the initial connection
        connectWebSocket();
    ]]
end

return createWatcher
