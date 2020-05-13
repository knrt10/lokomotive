# Controller Terraform module

This Terraform module aims to be a reusable module for generating controller nodes Ignition
configuration.

It build on top of [node](../node) module and adds some controller-specific settings on top of it,
like:
- extra `kubelet.service` dependencies etc.
- bootkube script and systemd unit
- etcd scripts and units
- controller labels
- controller taints

Additionally, it exposes various input variables, which allow to add platform-specific changes to the
ignition configuration.
