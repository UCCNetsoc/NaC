local grafana = import 'grafonnet-7.0/grafana.libsonnet';
local dashboard = grafana.dashboard;
local variable = grafana.template.query;
local graph = grafana.panel.graph;
local stat = grafana.panel.stat;
local gauge = grafana.panel.gauge;
local table = grafana.panel.table;
local prometheus = grafana.target.prometheus;

local vm_summary = table.new(
  datasource="Prometheus",
  title="Resource Allocation Summary",
)
.setGridPos(
  h=10,
  w=14,
  x=0,
  y=0,
)
.addTarget(
  prometheus.new(
    expr='sum(pve_guest_info{node="${host}"}) by (id)',
    format="table",
    instant=true,
  )
);

dashboard.new(
  editable=false,
  refresh="30s",
  title="PVE",
)
.setTime(
  from="now-1h",
)
.addTemplate(
  variable.new(
    datasource="Prometheus",
    label="Host",
    name="host",
    query="label_values(pve_node_info, name)",
    refresh=1,
    sort=1,
  )
)
.addPanels(
  [
    vm_summary
  ]
)