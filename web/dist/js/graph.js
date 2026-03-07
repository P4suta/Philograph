// Graph rendering module using Sigma.js v2 + Graphology
import { forceLayout } from './layout.js';

const COMMUNITY_COLORS = [
    '#8dd3c7', '#ffffb3', '#bebada', '#fb8072', '#80b1d3',
    '#fdb462', '#b3de69', '#fccde5', '#d9d9d9', '#bc80bd',
];

let sigmaInstance = null;
let graphInstance = null;
let hoveredNode = null;

export function initGraph(graphData) {
    const container = document.getElementById('graph-container');

    // Check if libraries are loaded
    if (typeof graphology === 'undefined' || typeof Sigma === 'undefined') {
        console.warn('Sigma.js or Graphology not loaded. Falling back to basic display.');
        renderFallback(container, graphData);
        return;
    }

    // Cleanup previous instance
    if (sigmaInstance) {
        sigmaInstance.kill();
        sigmaInstance = null;
    }

    hoveredNode = null;
    graphInstance = new graphology.Graph();

    // Add nodes
    if (graphData.Nodes) {
        graphData.Nodes.forEach(node => {
            const size = 5 + node.DegreeCentrality * 25;
            const color = COMMUNITY_COLORS[node.CommunityID % COMMUNITY_COLORS.length];
            graphInstance.addNode(node.ID, {
                label: node.Label,
                x: Math.random() * 100,
                y: Math.random() * 100,
                size: size,
                color: color,
                frequency: node.Frequency,
                degree_centrality: node.DegreeCentrality,
                betweenness_centrality: node.BetweennessCentrality,
                eigenvector_centrality: node.EigenvectorCentrality,
                community_id: node.CommunityID,
            });
        });
    }

    // Add edges
    if (graphData.Edges) {
        graphData.Edges.forEach((edge, i) => {
            const size = 0.5 + edge.Weight * 4.5;
            graphInstance.addEdge(edge.SourceID, edge.TargetID, {
                size: size,
                color: '#555',
                raw_count: edge.RawCount,
            });
        });
    }

    // Render with nodeReducer/edgeReducer for hover highlighting
    sigmaInstance = new Sigma(graphInstance, container, {
        renderEdgeLabels: false,
        defaultEdgeColor: '#555',
        labelColor: { color: '#e0e0e0' },
        nodeReducer(node, data) {
            const res = { ...data };
            if (hoveredNode) {
                if (node === hoveredNode) {
                    res.highlighted = true;
                } else if (graphInstance.hasEdge(hoveredNode, node) || graphInstance.hasEdge(node, hoveredNode)) {
                    // neighbor — keep as is
                } else {
                    res.color = '#2a2a4a';
                    res.label = '';
                }
            }
            return res;
        },
        edgeReducer(edge, data) {
            const res = { ...data };
            if (hoveredNode) {
                const src = graphInstance.source(edge);
                const tgt = graphInstance.target(edge);
                if (src !== hoveredNode && tgt !== hoveredNode) {
                    res.hidden = true;
                }
            }
            return res;
        },
    });

    // Hover handlers
    sigmaInstance.on('enterNode', ({ node }) => {
        hoveredNode = node;
        sigmaInstance.refresh();
    });
    sigmaInstance.on('leaveNode', () => {
        hoveredNode = null;
        sigmaInstance.refresh();
    });

    // Node click handler
    sigmaInstance.on('clickNode', ({ node }) => {
        const attrs = graphInstance.getNodeAttributes(node);
        showNodeDetail(node, attrs);
    });

    // Run force-directed layout
    forceLayout(graphInstance, 100);

    // Build top terms ranking
    buildTopTerms();

    // Build community legend
    buildCommunityLegend();

    // Setup search
    setupSearch();
}

function showNodeDetail(nodeId, attrs) {
    const panel = document.getElementById('detail-panel');
    const detail = document.getElementById('node-detail');
    panel.hidden = false;

    // Build adjacent nodes list
    let neighborsHTML = '';
    if (graphInstance && graphInstance.hasNode(nodeId)) {
        const neighbors = graphInstance.neighbors(nodeId);
        if (neighbors.length > 0) {
            neighborsHTML = '<dt>Adjacent Nodes</dt><dd><ul id="adjacent-list">';
            neighbors.forEach(n => {
                const nAttrs = graphInstance.getNodeAttributes(n);
                neighborsHTML += `<li class="adjacent-item" data-node="${n}">${nAttrs.label}</li>`;
            });
            neighborsHTML += '</ul></dd>';
        }
    }

    detail.innerHTML = `
        <dt>Label</dt><dd>${attrs.label}</dd>
        <dt>Frequency</dt><dd>${attrs.frequency}</dd>
        <dt>Degree Centrality</dt><dd>${attrs.degree_centrality?.toFixed(4)}</dd>
        <dt>Betweenness</dt><dd>${attrs.betweenness_centrality?.toFixed(4)}</dd>
        <dt>Eigenvector</dt><dd>${attrs.eigenvector_centrality?.toFixed(6)}</dd>
        <dt>Community</dt><dd>${attrs.community_id}</dd>
        ${neighborsHTML}
    `;

    // Adjacent node click handlers
    detail.querySelectorAll('.adjacent-item').forEach(item => {
        item.addEventListener('click', () => {
            const targetNode = item.dataset.node;
            const targetAttrs = graphInstance.getNodeAttributes(targetNode);
            focusNode(targetNode);
            showNodeDetail(targetNode, targetAttrs);
        });
    });
}

function focusNode(nodeId) {
    if (!sigmaInstance || !graphInstance || !graphInstance.hasNode(nodeId)) return;
    const attrs = graphInstance.getNodeAttributes(nodeId);
    sigmaInstance.getCamera().animate({ x: attrs.x, y: attrs.y, ratio: 0.3 }, { duration: 300 });
}

function setupSearch() {
    const input = document.getElementById('search-input');
    if (!input) return;
    let timer = null;
    input.addEventListener('input', () => {
        clearTimeout(timer);
        timer = setTimeout(() => {
            const query = input.value.trim().toLowerCase();
            if (!query || !graphInstance) return;
            let found = null;
            graphInstance.forEachNode((node, attrs) => {
                if (!found && attrs.label.toLowerCase().includes(query)) {
                    found = node;
                }
            });
            if (found) {
                focusNode(found);
            }
        }, 300);
    });
}

function buildTopTerms() {
    const section = document.getElementById('top-terms-section');
    const list = document.getElementById('top-terms-list');
    if (!section || !list || !graphInstance) return;

    const nodes = [];
    graphInstance.forEachNode((node, attrs) => {
        nodes.push({ node, ...attrs });
    });
    nodes.sort((a, b) => (b.betweenness_centrality || 0) - (a.betweenness_centrality || 0));
    const top10 = nodes.slice(0, 10);

    list.innerHTML = '';
    top10.forEach(n => {
        const li = document.createElement('li');
        li.textContent = `${n.label} (${(n.betweenness_centrality || 0).toFixed(4)})`;
        li.classList.add('top-term-item');
        li.addEventListener('click', () => {
            focusNode(n.node);
            showNodeDetail(n.node, graphInstance.getNodeAttributes(n.node));
        });
        list.appendChild(li);
    });

    section.hidden = false;
    // Show detail panel for top terms even before node click
    document.getElementById('detail-panel').hidden = false;
}

function buildCommunityLegend() {
    const section = document.getElementById('legend-section');
    const legend = document.getElementById('community-legend');
    if (!section || !legend || !graphInstance) return;

    const communities = new Map();
    graphInstance.forEachNode((node, attrs) => {
        const cid = attrs.community_id;
        communities.set(cid, (communities.get(cid) || 0) + 1);
    });

    legend.innerHTML = '';
    const sorted = [...communities.entries()].sort((a, b) => a[0] - b[0]);
    sorted.forEach(([cid, count]) => {
        const item = document.createElement('div');
        item.classList.add('legend-item');
        const swatch = document.createElement('span');
        swatch.classList.add('legend-swatch');
        swatch.style.backgroundColor = COMMUNITY_COLORS[cid % COMMUNITY_COLORS.length];
        item.appendChild(swatch);
        item.appendChild(document.createTextNode(`Community ${cid} (${count} nodes)`));
        legend.appendChild(item);
    });

    section.hidden = false;
}

function renderFallback(container, graphData) {
    container.innerHTML = `
        <div style="padding:2rem;color:#e0e0e0;">
            <h3>Graph Summary</h3>
            <p>Nodes: ${graphData.Nodes?.length || 0}</p>
            <p>Edges: ${graphData.Edges?.length || 0}</p>
            <p style="margin-top:1rem;color:#a0a0b0;">
                To enable interactive visualization, add Sigma.js and Graphology libraries to web/dist/vendor/.
            </p>
        </div>
    `;
}
