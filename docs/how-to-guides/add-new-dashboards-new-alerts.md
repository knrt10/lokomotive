# Add custom monitoring configs

## Contents

* [Introduction](#introduction)
* [Prerequisites](#prerequisites)
* [Add custom Grafana dashboards](#add-custom-grafana-dashboards)
* [Add new ServiceMonitors](#add-new-service-monitors)
  * [Default Prometheus operator setting](#default-prometheus-operator-setting)
  * [Custom Prometheus operator setting:](#custom-prometheus-operator-setting)
* [Add custom alerts for Alertmanager](#add-custom-alerts-for-alertmanager)
  * [Default Prometheus operator setting](#default-prometheus-operator-setting-1)
  * [Custom Prometheus operator setting:](#custom-prometheus-operator-setting-1)
* [Additional resources](#additional-resources)

## Introduction

Lokomotive ships monitoring stack in the form of prometheus-operator component. It monitors control plane and essential components by default. However, users can also use it to monitor their applications running on lokomotive. To enable such application monitoring user can add custom grafana dashboards, service monitors and alerting rules.

This guide shows how to add custom grafana dashboards, service monitor to discover and scrape applications, alerting rules for alertmanager.

## Prerequisites

- Deploy prometheus-operator component on lokomotive.

## Add custom Grafana dashboards

Create a ConfigMap with keys as the dashboard file names and values as JSON dashboard. See the following example:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: dashboard
  labels:
    grafana_dashboard: "true"
data:
  metallb.json: |
    {
      "annotations": {
[REDACTED]
```

Add the label `grafana_dashboard: "true"` so that grafana automatically picks up the dashboards in the ConfigMaps across the cluster.

## Add new ServiceMonitors

#### Default Prometheus operator setting:

Create a ServiceMonitor with the required configuration and make sure to add the following label, so that the prometheus-operator will track it:

```yaml
metadata:
  labels:
    release: prometheus-operator
```

#### Custom Prometheus operator setting:

Deploy the prometheus-operator with the following setting, and it watches all ServiceMonitors across the cluster:

```tf
watch_labeled_service_monitors = "false"
```

Then there is no need to add any label to ServiceMonitor, at all. Create a ServiceMonitor, and prometheus-operator tracks it.

## Add custom alerts for Alertmanager

#### Default Prometheus operator setting:

Create a PrometheuRule object with the required configuration and make sure to add the following labels, so that prometheus-operator will track it:

```yaml
metadata:
  labels:
    release: prometheus-operator
    app: prometheus-operator
```

#### Custom Prometheus operator setting:

Deploy the prometheus-operator with the following setting, and it watches all PrometheusRules across the cluster:

```tf
watch_labeled_prometheus_rules = "false"
```

Then there is no need to add any label to PrometheusRule, at all. Create a PrometheusRule, and prometheus-operator tracks it.

## Additional resources

- ServiceMonitor API docs https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#servicemonitor
- PrometheusRule API docs https://github.com/coreos/prometheus-operator/blob/master/Documentation/api.md#prometheusrule
