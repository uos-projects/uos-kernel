const API_BASE = "/api";

// 实体类型中文映射
const entityTypeMap = {
  "substation": "变电站",
  "voltage_level": "电压等级",
  "power_transformer": "电力变压器",
  "breaker": "断路器",
  "asset": "资产",
  "work": "工作单",
  "meter": "电表",
  "usage_point": "用电点",
};

// 关系类型中文映射
const relationTypeMap = {
  "contains": "包含",
  "partOf": "属于",
  "locatedIn": "位于",
  "realises": "实现",
  "targets": "目标",
  "measuredBy": "测量",
  "connectedVia": "通过连接",
  "connects": "连接",
  "boundTo": "绑定",
  "aggregates": "聚合",
  "hasEnd": "有端",
  "operatesOn": "操作",
  "includes": "包含",
  "feeds": "供电",
  "captures": "捕获",
  "associatedWith": "关联",
};

// 获取实体类型中文名称
function getEntityTypeName(type) {
  return entityTypeMap[type] || type;
}

// 获取关系类型中文名称
function getRelationTypeName(type) {
  return relationTypeMap[type] || type;
}

const modeSelect = document.getElementById("mode");
const valueInput = document.getElementById("value");
const queryBtn = document.getElementById("query-btn");
const snapshotsBtn = document.getElementById("snapshots-btn");
const statusEl = document.getElementById("status");
const timelineInput = document.getElementById("timeline");
const timelineValue = document.getElementById("timeline-value");
const timelineMarkers = document.getElementById("timeline-markers");

const entitiesBody = document.querySelector("#entities-table tbody");
const relationsBody = document.querySelector("#relations-table tbody");
const snapshotsBody = document.querySelector("#snapshots-table tbody");

const svg = d3.select("#graph");
// 获取 SVG 容器的实际宽度
function getSvgWidth() {
  const container = svg.node().parentElement;
  return container ? container.clientWidth - 32 : 800; // 减去 padding
}
const height = +svg.attr("height") || 600;
// 初始化 SVG 宽度
const initialWidth = getSvgWidth();
svg.attr("viewBox", `0 0 ${initialWidth} ${height}`);
svg.append("defs")
  .append("marker")
  .attr("id", "arrowhead")
  .attr("viewBox", "-0 -5 10 10")
  .attr("refX", 23)
  .attr("refY", 0)
  .attr("orient", "auto")
  .attr("markerWidth", 6)
  .attr("markerHeight", 6)
  .attr("xoverflow", "visible")
  .append("svg:path")
  .attr("d", "M 0,-5 L 10 ,0 L 0,5")
  .attr("fill", "#94a3b8");

let simulation;

queryBtn.addEventListener("click", () => {
  const mode = modeSelect.value;
  const value = valueInput.value.trim();
  if (!value) {
    statusEl.textContent = "请输入时间/ID";
    return;
  }
  fetchView(mode, value);
});

snapshotsBtn.addEventListener("click", () => loadSnapshots());

// 时间轴处理
function updateTimelineValue(dayOffset) {
  const baseDate = new Date("2025-01-01T00:00:00Z");
  const targetDate = new Date(baseDate);
  targetDate.setDate(baseDate.getDate() + dayOffset);
  const dateStr = targetDate.toISOString().split("T")[0];
  timelineValue.textContent = dateStr;
  return targetDate.toISOString();
}

// 存储 snapshots 数据，用于重新渲染标记点
let currentSnapshots = [];

// 更新时间轴标记点
function updateTimelineMarkers(snapshots) {
  currentSnapshots = snapshots || [];
  
  if (!currentSnapshots || currentSnapshots.length === 0) {
    timelineMarkers.innerHTML = "";
    return;
  }

  const baseDate = new Date("2025-01-01T00:00:00Z");
  const maxDays = 150; // 对应时间轴的 max 值
  
  // 提取唯一的时间点（按天）
  const uniqueDates = new Set();
  currentSnapshots.forEach(snap => {
    try {
      const snapDate = new Date(snap.committed_at);
      const dayOffset = Math.floor((snapDate - baseDate) / (1000 * 60 * 60 * 24));
      if (dayOffset >= 0 && dayOffset <= maxDays) {
        uniqueDates.add(dayOffset);
      }
    } catch (e) {
      // 忽略日期解析错误
    }
  });

  // 清空现有标记
  timelineMarkers.innerHTML = "";

  // 创建标记点
  Array.from(uniqueDates).sort((a, b) => a - b).forEach(dayOffset => {
    const marker = document.createElement("div");
    marker.className = "timeline-marker";
    const dateStr = updateTimelineValue(dayOffset).split("T")[0];
    marker.title = `Snapshot at ${dateStr}`;
    
    // 计算位置（百分比）
    const position = (dayOffset / maxDays) * 100;
    marker.style.left = `${position}%`;
    
    // 点击标记点跳转到对应时间
    marker.addEventListener("click", (e) => {
      e.stopPropagation();
      timelineInput.value = dayOffset;
      const isoString = updateTimelineValue(dayOffset);
      if (modeSelect.value === "biz") {
        valueInput.value = isoString;
        fetchView("biz", isoString);
      }
    });
    
    timelineMarkers.appendChild(marker);
  });
}

timelineInput.addEventListener("input", (e) => {
  const dayOffset = parseInt(e.target.value);
  const isoString = updateTimelineValue(dayOffset);
  // 如果当前是业务时间模式，自动查询
  if (modeSelect.value === "biz") {
    valueInput.value = isoString;
    fetchView("biz", isoString);
  }
});

// 初始化时间轴值
updateTimelineValue(parseInt(timelineInput.value));

async function fetchView(mode, value) {
  statusEl.textContent = "查询中...";
  try {
    const params = new URLSearchParams({ mode, value });
    const res = await fetch(`${API_BASE}/view?${params.toString()}`);
    if (!res.ok) {
      const err = await res.json();
      throw new Error(err.detail || res.statusText);
    }
    const data = await res.json();
    renderEntities(data.entities);
    renderRelations(data.relations);
    renderGraph(data.entities, data.relations);
    statusEl.textContent = `视图已加载: mode=${data.mode}, value=${data.value}`;
    
    // 如果使用业务时间模式，同步时间轴
    if (mode === "biz" && value.match(/^\d{4}-\d{2}-\d{2}/)) {
      try {
        const date = new Date(value);
        const baseDate = new Date("2025-01-01T00:00:00Z");
        const dayOffset = Math.floor((date - baseDate) / (1000 * 60 * 60 * 24));
        if (dayOffset >= 0 && dayOffset <= 150) {
          timelineInput.value = dayOffset;
          updateTimelineValue(dayOffset);
        }
      } catch (e) {
        // 忽略日期解析错误
      }
    }
  } catch (err) {
    console.error(err);
    statusEl.textContent = `查询失败: ${err.message}`;
  }
}

async function loadSnapshots() {
  statusEl.textContent = "读取 snapshot...";
  try {
    console.log("请求:", `${API_BASE}/snapshots`);
    const res = await fetch(`${API_BASE}/snapshots`);
    console.log("响应状态:", res.status, res.ok, "Content-Type:", res.headers.get("content-type"));
    if (!res.ok) {
      const text = await res.text();
      console.error("错误响应:", text);
      throw new Error(`获取 snapshot 失败: ${res.status} - ${text}`);
    }
    const text = await res.text();
    console.log("原始响应文本:", text);
    let data;
    try {
      data = JSON.parse(text);
    } catch (e) {
      console.error("JSON 解析失败:", e, "文本:", text);
      throw new Error(`JSON 解析失败: ${e.message}`);
    }
    console.log("解析后的数据:", data);
    snapshotsBody.innerHTML = "";
    if (!data.snapshots || data.snapshots.length === 0) {
      statusEl.textContent = "没有 snapshot 数据";
      return;
    }
    data.snapshots.forEach((snap) => {
      const tr = document.createElement("tr");
      // 提取表名（例如 ontology.grid.substation -> substation）
      const tableName = snap.table ? snap.table.split(".").pop() : "-";
      tr.innerHTML = `
        <td>${tableName}</td>
        <td>${snap.snapshot_id}</td>
        <td>${snap.committed_at}</td>
        <td>${snap.operation}</td>
      `;
      tr.addEventListener("click", () => {
        modeSelect.value = "snapshot";
        valueInput.value = snap.snapshot_id;
      });
      snapshotsBody.appendChild(tr);
    });
    statusEl.textContent = `共 ${data.snapshots.length} 个 snapshot`;
    
    // 更新时间轴标记点
    updateTimelineMarkers(data.snapshots);
  } catch (err) {
    statusEl.textContent = err.message;
  }
}

function renderEntities(entities) {
  entitiesBody.innerHTML = "";
  entities.forEach((e) => {
    const tr = document.createElement("tr");
    const entityTypeName = e.entity_type ? getEntityTypeName(e.entity_type) : "-";
    tr.innerHTML = `
      <td>${e.entity_id ?? e.entityId ?? "-"}</td>
      <td>${entityTypeName}</td>
      <td>${e.name ?? e.entity_id ?? "-"}</td>
      <td>${e.region ?? e.substation_id ?? "-"}</td>
      <td>${e.valid_from ?? "-"}</td>
      <td>${e.valid_to ?? "-"}</td>
    `;
    entitiesBody.appendChild(tr);
  });
}

function renderRelations(relations) {
  relationsBody.innerHTML = "";
  relations.forEach((r) => {
    const tr = document.createElement("tr");
    const relationTypeName = r.relation_type ? getRelationTypeName(r.relation_type) : "-";
    tr.innerHTML = `
      <td>${r.rel_id}</td>
      <td>${r.source_id}</td>
      <td>${r.target_id}</td>
      <td>${relationTypeName}</td>
      <td>${r.valid_from}</td>
      <td>${r.valid_to}</td>
    `;
    relationsBody.appendChild(tr);
  });
}

function renderGraph(entities, relations) {
  // 动态获取 SVG 宽度
  const currentWidth = getSvgWidth();
  // 更新 viewBox 以适应容器宽度
  svg.attr("viewBox", `0 0 ${currentWidth} ${height}`);
  
  const nodesMap = new Map();
  entities.forEach((e) => {
    nodesMap.set(e.entity_id, {
      id: e.entity_id,
      label: e.name || e.entity_id,
      type: e.entity_type,
    });
  });
  relations.forEach((r) => {
    if (!nodesMap.has(r.source_id)) {
      nodesMap.set(r.source_id, { id: r.source_id, label: r.source_id, type: "relation" });
    }
    if (!nodesMap.has(r.target_id)) {
      nodesMap.set(r.target_id, { id: r.target_id, label: r.target_id, type: "relation" });
    }
  });

  const nodes = Array.from(nodesMap.values());
  const links = relations.map((r) => ({
    source: r.source_id,
    target: r.target_id,
    relation_type: r.relation_type,
  }));

  svg.selectAll("*:not(defs)").remove();

  const linkGroup = svg.append("g").attr("class", "links");
  
  const link = linkGroup
    .selectAll("line")
    .data(links)
    .enter()
    .append("line")
    .attr("class", "link")
    .attr("stroke-opacity", 0.7);

  // 添加连线标签（关系类型）
  const linkLabels = linkGroup
    .selectAll("text")
    .data(links)
    .enter()
    .append("text")
    .attr("class", "link-label")
    .attr("font-size", "10px")
    .attr("fill", "#64748b")
    .attr("text-anchor", "middle")
    .text((d) => getRelationTypeName(d.relation_type));

  const node = svg
    .append("g")
    .selectAll("g")
    .data(nodes)
    .enter()
    .append("g")
    .attr("class", "node");

  node
    .append("circle")
    .attr("r", (d) => (d.type === "substation" ? 14 : 10))
    .attr("fill", (d) => (d.type === "substation" ? "#2563eb" : "#38bdf8"))
    .attr("stroke", "#1d4ed8")
    .attr("stroke-width", 1.5);

  const textElements = node
    .append("text")
    .text((d) => d.label)
    .attr("x", 15)
    .attr("y", 4)
    .attr("fill", "#000000")
    .attr("stroke", "none")
    .attr("font-size", "11px")
    .style("fill", "#000000")
    .style("stroke", "none");
  
  // 确保文字颜色不被覆盖，使用 each 直接设置属性
  textElements.each(function() {
    this.setAttribute("fill", "#000000");
    this.setAttribute("stroke", "none");
    this.style.setProperty("fill", "#000000", "important");
    this.style.setProperty("stroke", "none", "important");
  });

  // 添加拖拽功能
  const drag = d3
    .drag()
    .on("start", function(event, d) {
      if (!event.active) simulation.alphaTarget(0.3).restart();
      d.fx = d.x;
      d.fy = d.y;
    })
    .on("drag", function(event, d) {
      const [x, y] = d3.pointer(event, svg.node());
      d.fx = x;
      d.fy = y;
    })
    .on("end", function(event, d) {
      if (!event.active) simulation.alphaTarget(0);
      d.fx = null;
      d.fy = null;
    });

  // 在节点上应用拖拽行为
  node.call(drag);

  // 设置鼠标样式
  node.style("cursor", "move");

  simulation = d3
    .forceSimulation(nodes)
    .force(
      "link",
      d3
        .forceLink(links)
        .id((d) => d.id)
        .distance(80) // 减小距离，使图形更紧凑
    )
    .force("charge", d3.forceManyBody().strength(-150)) // 减小排斥力，使节点更靠近
    .force("center", d3.forceCenter(currentWidth / 2, height / 2))
    .force("collision", d3.forceCollide().radius(20)) // 添加碰撞检测，防止节点重叠
    .on("tick", () => {
      link
        .attr("x1", (d) => d.source.x)
        .attr("y1", (d) => d.source.y)
        .attr("x2", (d) => d.target.x)
        .attr("y2", (d) => d.target.y);

      // 更新连线标签位置（显示在连线中间）
      linkLabels
        .attr("x", (d) => (d.source.x + d.target.x) / 2)
        .attr("y", (d) => (d.source.y + d.target.y) / 2);

      node.attr("transform", (d) => `translate(${d.x},${d.y})`);
    });
}

// 窗口大小改变时重新计算 SVG 宽度
window.addEventListener("resize", () => {
  const currentWidth = getSvgWidth();
  svg.attr("viewBox", `0 0 ${currentWidth} ${height}`);
  if (simulation) {
    simulation.force("center", d3.forceCenter(currentWidth / 2, height / 2));
    simulation.alpha(0.3).restart();
  }
});

// 默认填充一个未来时间，避免空输入
valueInput.value = "2100-01-01T00:00:00Z";
fetchView("biz", valueInput.value);
// 自动加载 snapshots 以显示时间轴标记点
loadSnapshots();

