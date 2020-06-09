# How to upgrade bootstrap Kubelet

## Contents

- [Introduction](#introduction)
- [Steps](#steps)
  - [Drain the node](#drain-the-node)
  - [Find out the IP and SSH](#find-out-the-ip-and-ssh)
  - [On the node](#on-the-node)
- [Caveats](#caveats)

## Introduction

Bootstrap Kubelet is the Kubelet running on the node as a systemd service outside Kubernetes cluster's control. Unlike most cluster components, `lokoctl` cannot update it(which may change in the future). This document enlists the manual steps involved in updating the bootstrap Kubelet.

## Steps

Apply the following steps to all the cluster nodes one at a time.

### Drain the node

> **Caution:** If you are using a local directory as a storage for a workload, it will be disturbed by this operation. So move the workload to another node and let the application replicate the data. If the application does not support data replication across instances, then expect downtime.

```bash
kubectl drain --ignore-daemonsets <node name>
```

### Find out the IP and SSH

Find the IP of the node by visiting the Cloud provider dashboard.

```
ssh core@<IP Address>
```

### On the node

Run the following commands.

> **NOTE**: Export the correct and latest Kubernetes version, before proceeding to other commands.

```
export latest_kube=<latest kubernetes version e.g. v1.18.0>
sudo sed -i "s|$(grep -i kubelet_image_tag /etc/kubernetes/kubelet.env)|KUBELET_IMAGE_TAG=${latest_kube}|g" /etc/kubernetes/kubelet.env
sudo systemctl restart kubelet
sudo journalctl -fu kubelet
```

Look at the logs carefully, if Kubelet fails to restart and instructs to do something like deleting a file, then delete it and restart the node.

```
sudo reboot
```

## Caveats

- If a node which is running storage workload like Rook Ceph, then verify that the Ceph cluster is in **`HEALTH_OK`** state. If it is in **any other state, do not proceed with upgrades**. When the cluster is in `HEALTH_OK` state a few minutes of downtime of OSDs will not be a problem for Ceph cluster, that happens during node reboot.
