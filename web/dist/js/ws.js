// WebSocket client for progress notifications
let ws = null;
let onProgress = null;
let onComplete = null;

export function initWebSocket(progressCallback, completeCallback) {
    onProgress = progressCallback;
    onComplete = completeCallback || null;
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
            if (data.stage === 'complete' && onComplete) {
                onComplete();
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
