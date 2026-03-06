// Parameter controls module
let reloadCallback = null;
let statusCallback = null;

export function initControls(onReload, onStatus) {
    reloadCallback = onReload;
    statusCallback = onStatus;

    // Slider bindings
    bindSlider('window-size', 'window-size-val');
    bindSlider('min-freq', 'min-freq-val');
    bindSlider('max-nodes', 'max-nodes-val');

    // Debounced config update on slider change
    const sliders = ['window-size', 'min-freq', 'max-nodes'];
    sliders.forEach(id => {
        document.getElementById(id).addEventListener('change', () => updateConfig());
    });

    document.getElementById('metric').addEventListener('change', () => updateConfig());

    // Export buttons
    document.getElementById('export-json').addEventListener('click', () => exportGraph('json'));
    document.getElementById('export-gexf').addEventListener('click', () => exportGraph('gexf'));
}

function bindSlider(sliderId, valueId) {
    const slider = document.getElementById(sliderId);
    const valueEl = document.getElementById(valueId);
    slider.addEventListener('input', () => { valueEl.textContent = slider.value; });
}

async function updateConfig() {
    const config = {
        window_size: parseInt(document.getElementById('window-size').value),
        min_frequency: parseInt(document.getElementById('min-freq').value),
        max_nodes: parseInt(document.getElementById('max-nodes').value),
        metric: document.getElementById('metric').value,
    };

    if (statusCallback) statusCallback('Re-analyzing...');

    try {
        const res = await fetch('/api/v1/config', {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(config),
        });
        const data = await res.json();
        if (!res.ok) {
            if (statusCallback) statusCallback('Error: ' + (data.error || 'Unknown'));
            return;
        }
        if (statusCallback) statusCallback(`Updated: ${data.nodes} nodes, ${data.edges} edges`);
        if (reloadCallback) await reloadCallback();
    } catch (err) {
        if (statusCallback) statusCallback('Error: ' + err.message);
    }
}

function exportGraph(format) {
    window.open(`/api/v1/export?format=${format}`, '_blank');
}
