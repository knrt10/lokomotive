# Storage with Rook Ceph on Packet cloud

## Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Deploy storage worker pool](#deploy-storage-worker-pool)
  - [Config](#config)
  - [Deploy the worker pool](#deploy-the-worker-pool)
- [Deploy Rook](#deploy-rook)
  - [Config](#config-1)
  - [Deploy the component](#deploy-the-component)
- [Deploy Rook Ceph](#deploy-rook-ceph)
  - [Config](#config-2)
  - [Deploy the component](#deploy-the-component-1)
- [Access the Ceph dashboard](#access-ceph-dashboard)
- [Enable and access toolbox](#enable-and-access-toolbox)
- [Enable monitoring](#enable-monitoring)
- [Make default storage class](#make-default-storage-class)
- [Additional resources](#additional-resources)

## Introduction

This guide provides the steps for deploying a storage stack using the `rook` and `rook-ceph` Lokomotive component and explains how to access Ceph dashboard, Ceph toolbox and how to enable monitoring.

## Prerequisites

* A Lokomotive cluster deployed on a supported provider and accessible via `kubectl`.

## Deploy storage worker pool

#### Config

Deploy a cluster with at least one worker pool dedicated to `rook-ceph`. A dedicated worker pool configuration should look like the following:

```tf
cluster "packet" {
  ...

  worker_pool "storage" {
    count = 3
    node_type = "c2.medium.x86"

    labels = "storage.lokomotive.io=ceph"
    taints = "storage.lokomotive.io=ceph:NoSchedule"
  }
}
```

- The number of machines provided using `count` should be an odd number greater than three.
- Type of node, provided using `node_type`, should be one that has multiple disks like `c2.medium.x86` or `s1.large.x86`. Find out more servers [here](https://www.packet.com/cloud/servers/).
- To steer `rook-ceph` workload on these storage nodes provide `labels`.
- Provide `taints` so that other workload can be **steered away** by default. This setting is not mandatory, but isolating storage workloads from others is recommended.

#### Deploy the worker pool

Execute the following command to deploy the `storage` worker pool:

```bash
lokoctl cluster apply -v --skip-components
```

## Deploy `rook`

#### Config

Create a file named `storage.lokocfg` with the following contents:

```tf
component "rook" {
  node_selector = {
    "storage.lokomotive.io" = "ceph"
  }

  toleration {
    key      = "storage.lokomotive.io"
    operator = "Equal"
    value    = "ceph"
    effect   = "NoSchedule"
  }

  agent_toleration_key    = "storage.lokomotive.io"
  agent_toleration_effect = "NoSchedule"

  discover_toleration_key    = "storage.lokomotive.io"
  discover_toleration_effect = "NoSchedule"
}
```

- `node_selector` should match the `labels` variable provided in the `worker_pool`.
- `toleration` should match the `taints` variable mentioned in the `worker_pool`.
- `agent_toleration_key` and `discover_toleration_key` should match the `key` of the `taints` variable provided in the `worker_pool`.
- `agent_toleration_effect` and `discover_toleration_effect` should match the `effect` of the `taints` variable provided in the `worker_pool`.

For information about all the available configuration options for the `rook` component, visit the component's [configuration reference](../configuration-reference/components/rook.md).

#### Deploy the component

Execute the following command to deploy the `rook` component:

```bash
lokoctl component apply rook
```

Verify the operator pod in the `rook` namespace is in the `Running` state (this may take a few minutes):

```console
$ kubectl -n rook get pods -l app=rook-ceph-operator
NAME                                  READY   STATUS    RESTARTS   AGE
rook-ceph-operator-76d8687f95-6knf8   1/1     Running   0          2m
```

## Deploy `rook-ceph`

#### Config

Add following contents to the previously created file `storage.lokocfg`:

```tf
component "rook-ceph" {
  monitor_count = 3

  node_affinity {
    key      = "storage.lokomotive.io"
    operator = "Exists"
  }

  toleration {
    key      = "storage.lokomotive.io"
    operator = "Equal"
    value    = "ceph"
    effect   = "NoSchedule"
  }

  storage_class {
    enable = true
  }
}
```

- `monitor_count` should be an odd number greater than three and not higher than the `count` variable of workers in the `worker_pool`.
- `node_affinity` should match the `labels` variable provided in the `worker_pool`.
- `toleration` should match the `taints` variable provided in the `worker_pool`.

For information about all the available configuration options for the `rook-ceph` component, visit the component's [configuration reference](../configuration-reference/components/rook-ceph.md).

#### Deploy the component

Execute the following command to deploy the `rook-ceph` component:

```bash
lokoctl component apply rook-ceph
```

Verify the [OSD](https://docs.ceph.com/docs/master/glossary/#term-ceph-osd-daemon) pods in the `rook` namespace are in the `Running` state (this may take a few minutes):

```console
$ kubectl -n rook get pods -l app=rook-ceph-osd
NAME                               READY   STATUS    RESTARTS   AGE
rook-ceph-osd-0-6d4f69dbf9-26kzl   1/1     Running   0          15m
rook-ceph-osd-1-86c9597b84-lmh94   1/1     Running   0          15m
rook-ceph-osd-2-6d97697897-7bprl   1/1     Running   0          15m
rook-ceph-osd-3-5bfb9d86b-rk6v4    1/1     Running   0          15m
rook-ceph-osd-4-5b76cb9675-cxkdw   1/1     Running   0          15m
rook-ceph-osd-5-8c86f5c6c-6qxtz    1/1     Running   0          15m
rook-ceph-osd-6-5b9cc479b7-vjc9v   1/1     Running   0          15m
rook-ceph-osd-7-7b84d6cc48-b46z9   1/1     Running   0          15m
rook-ceph-osd-8-5868969f97-2bn9r   1/1     Running   0          15m
```

## Access the Ceph dashboard

Ceph dashboard provides valuable visual information. It is an essential tool to monitor the Ceph cluster. Here are steps on how to access it.

Obtain the password for the `admin` Ceph user by running the following command:

```bash
kubectl -n rook get secret rook-ceph-dashboard-password -o jsonpath="{['data']['password']}" | base64 --decode && echo
```

Execute the following command to forward port `8443` locally to the Ceph manager pod:

```bash
kubectl -n rook port-forward svc/rook-ceph-mgr-dashboard 8443
```

Now open the following URL: [https://localhost:8443](https://localhost:8443) and enter the username `admin` and the password obtained from the first step.

## Enable and access toolbox

Ceph is a complex software system, and not everything that happens in the Ceph cluster is visible at the `rook` layer of abstraction. So command-line interface to interact with Ceph cluster is useful to extract such hidden events and information. Ceph toolbox helps you access the ceph cluster using `ceph` CLI utility. Using the utility you can configure the Ceph cluster setting and debug the cluster.

To deploy the toolbox, the `rook-ceph` component config should set the variable `enable_toolbox` to `true`.

```tf
component "rook-ceph" {
  enable_toolbox = true
  ...
}
```

Execute the following command to apply the changes:

```bash
lokoctl component apply rook-ceph
```

Verify the toolbox pod in the `rook` namespace is in the `Running` state (this may take a few minutes):

```console
$ kubectl -n rook get deploy rook-ceph-tools
NAME              READY   UP-TO-DATE   AVAILABLE   AGE
rook-ceph-tools   1/1     1            1           39s
```

Execute the following command to access the toolbox pod:

```bash
kubectl -n rook exec -it $(kubectl -n rook get pods -l app=rook-ceph-tools -o name) bash
```

Once inside the pod you can run usual `ceph` commands:

```bash
ceph status
ceph osd status
ceph df
rados df
```

## Enable monitoring

Monitor `rook` and `rook-ceph` components using the `prometheus-operator` component. To enable your `rook` component config should have the variable `enable_monitoring` set to `true`.

> **NOTE:** Deploy the `prometheus-operator` component before.

```tf
component "rook" {
  enable_monitoring = true
  ...
}
```

Execute the following command to apply the changes:

```bash
lokoctl component apply rook
```

<!---
TODO: Once we have a tutorial on how to deploy and configure `prometheus-operator`, then point the following to that tutorial.
-->

## Make default storage class

It is recommended to make the storage class as default if `rook-ceph` is the only storage provider in the cluster. This setting helps to provision volumes for the [PVCs](https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims) created by workloads. The `rook-ceph` component config should look like the following:

```tf
component "rook-ceph" {
  ...

  storage_class {
    enable  = true
    default = true
  }
}
```

Execute the following command to apply the changes:

```bash
lokoctl component apply rook
```

Verify the StorageClass is default:

```console
$ kubectl get sc rook-ceph-block
NAME                        PROVISIONER             RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
rook-ceph-block (default)   rook.rbd.csi.ceph.com   Delete          Immediate           true                   8m17s
```

## Additional resources

- `rook` component [configuration reference](../configuration-reference/components/rook.md) guide.
- `rook-ceph` component [configuration reference](../configuration-reference/components/rook-ceph.md) guide.
- Rook docs:

  - [Ceph toolbox](https://rook.io/docs/rook/master/ceph-toolbox.html).
  - [Ceph dashboard](https://rook.io/docs/rook/master/ceph-dashboard.html).
  - [Ceph direct tools](https://rook.io/docs/rook/master/direct-tools.html).
  - [Ceph advanced configuration](https://rook.io/docs/rook/master/ceph-advanced-configuration.html).
  - [Disaster recovery](https://rook.io/docs/rook/master/ceph-disaster-recovery.html).
