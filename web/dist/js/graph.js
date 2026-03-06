// Graph rendering module using Sigma.js v2 + Graphology
const COMMUNITY_COLORS = [
    '#e94560', '#4ecdc4', '#ffe66d', '#95e1d3', '#f38181',
    '#aa96da', '#fcbad3', '#a8d8ea', '#ffd3b6', '#c7ecee',
];

let sigmaInstance = null;
let graphInstance = null;

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

    // Render
    sigmaInstance = new Sigma(graphInstance, container, {
        renderEdgeLabels: false,
        defaultEdgeColor: '#555',
        labelColor: { color: '#e0e0e0' },
    });

    // Node click handler
    sigmaInstance.on('clickNode', ({ node }) => {
        const attrs = graphInstance.getNodeAttributes(node);
        showNodeDetail(attrs);
    });

    // Run ForceAtlas2 layout if available
    if (typeof graphologyLayoutForceAtlas2 !== 'undefined') {
        graphologyLayoutForceAtlas2.assign(graphInstance, {
            iterations: 100,
            settings: { gravity: 1, scalingRatio: 10 },
        });
    }
}

function showNodeDetail(attrs) {
    const panel = document.getElementById('detail-panel');
    const detail = document.getElementById('node-detail');
    panel.hidden = false;
    detail.innerHTML = `
        <dt>Label</dt><dd>${attrs.label}</dd>
        <dt>Frequency</dt><dd>${attrs.frequency}</dd>
        <dt>Degree Centrality</dt><dd>${attrs.degree_centrality?.toFixed(4)}</dd>
        <dt>Betweenness</dt><dd>${attrs.betweenness_centrality?.toFixed(4)}</dd>
        <dt>Eigenvector</dt><dd>${attrs.eigenvector_centrality?.toFixed(6)}</dd>
        <dt>Community</dt><dd>${attrs.community_id}</dd>
    `;
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
