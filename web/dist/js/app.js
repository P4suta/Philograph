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

let analysisStartTime = null;

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

    analysisStartTime = performance.now();
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
        updateStats(graph);
    } catch (err) {
        console.error('Failed to load result:', err);
    }
}

function updateStats(graphData) {
    const section = document.getElementById('stats-section');
    const detail = document.getElementById('stats-detail');
    if (!section || !detail) return;

    const nodeCount = graphData.Nodes?.length || 0;
    const edgeCount = graphData.Edges?.length || 0;

    const communities = new Set();
    if (graphData.Nodes) {
        graphData.Nodes.forEach(n => communities.add(n.CommunityID));
    }

    let elapsed = '—';
    if (analysisStartTime) {
        const ms = performance.now() - analysisStartTime;
        elapsed = ms < 1000 ? `${Math.round(ms)}ms` : `${(ms / 1000).toFixed(1)}s`;
    }

    detail.innerHTML = `
        <dt>Nodes</dt><dd>${nodeCount}</dd>
        <dt>Edges</dt><dd>${edgeCount}</dd>
        <dt>Communities</dt><dd>${communities.size}</dd>
        <dt>Processing Time</dt><dd>${elapsed}</dd>
    `;
    section.hidden = false;
}

async function onAnalysisComplete() {
    controlsSection.hidden = false;
    await loadResult();
}

// Initialize
uploadForm.addEventListener('submit', handleUpload);
initWebSocket(showProgress, onAnalysisComplete);
initControls(loadResult, setStatus);

// Try to load existing result (for CLI-initiated analysis)
loadResult();
