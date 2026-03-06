// Philograph - Main Application Module
import { initWebSocket } from './ws.js';
import { initControls } from './controls.js';
import { initGraph } from './graph.js';

const statusEl = document.getElementById('status');
const uploadForm = document.getElementById('upload-form');
const fileInput = document.getElementById('file-input');
const progressContainer = document.getElementById('progress-bar-container');
const progressBar = document.getElementById('progress-bar');
const progressText = document.getElementById('progress-text');
const controlsSection = document.getElementById('controls-section');

function setStatus(text) {
    statusEl.textContent = text;
}

function showProgress(percentage, message) {
    progressContainer.hidden = false;
    progressBar.style.width = percentage + '%';
    progressText.textContent = message;
    if (percentage >= 100) {
        setTimeout(() => { progressContainer.hidden = true; }, 1000);
    }
}

async function handleUpload(e) {
    e.preventDefault();
    const file = fileInput.files[0];
    if (!file) return;

    setStatus('Analyzing...');
    const formData = new FormData();
    formData.append('file', file);

    try {
        const res = await fetch('/api/v1/analyze', { method: 'POST', body: formData });
        const data = await res.json();
        if (!res.ok) {
            setStatus('Error: ' + (data.error || 'Unknown error'));
            return;
        }
        setStatus(`Done: ${data.nodes} nodes, ${data.edges} edges`);
        controlsSection.hidden = false;
        await loadResult();
    } catch (err) {
        setStatus('Error: ' + err.message);
    }
}

async function loadResult() {
    try {
        const res = await fetch('/api/v1/result');
        if (!res.ok) return;
        const graph = await res.json();
        initGraph(graph);
    } catch (err) {
        console.error('Failed to load result:', err);
    }
}

// Initialize
uploadForm.addEventListener('submit', handleUpload);
initWebSocket(showProgress);
initControls(loadResult, setStatus);
