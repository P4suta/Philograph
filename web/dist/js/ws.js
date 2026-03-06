// WebSocket client for progress notifications
let ws = null;
let onProgress = null;

export function initWebSocket(progressCallback) {
    onProgress = progressCallback;
    connect();
}

function connect() {
    const protocol = location.protocol === 'https:' ? 'wss:' : 'ws:';
    ws = new WebSocket(`${protocol}//${location.host}/ws`);

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            if (onProgress) {
                onProgress(data.percentage, data.message);
            }
        } catch (e) {
            console.error('WS parse error:', e);
        }
    };

    ws.onclose = () => {
        setTimeout(connect, 3000);
    };

    ws.onerror = () => {
        ws.close();
    };
}
