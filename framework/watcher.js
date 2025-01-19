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
    
    ws.onclose = function() {
        console.log('WebSocket connection closed. Reconnecting in 1s...');
        setTimeout(connectWebSocket, 1000);
    };
    
    ws.onerror = function(err) {
        console.error('WebSocket error:', err);
        ws.close();
    };
}

connectWebSocket();
