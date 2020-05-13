# Worker Terraform module

This Terraform module aims to be a reusable module for generating worker nodes Ignition
configuration.

It build on top of [node](../node) module and adds some worker-specific settings on top of it,
like:
- kubeconfig file for kubelet
- iscsid service and bind-mounts for kubelet container
- default worker node labels

Additionally, it exposes various input variables, which allow to add platform-specific changes to the
ignition configuration.
