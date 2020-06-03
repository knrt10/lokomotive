# How to setup third party OAuth for Grafana?

## Contents

* [Introduction](#introduction)
* [Prerequisites](#prerequisites)
* [Steps](#steps)
* [What's next?](#whats-next)

## Introduction

This document will help you to enable any supported auth provider on Grafana deployed as a part of Prometheus Operator.

## Prerequisites

- On Packet: You have a DNS entry in any DNS provider for `grafana.mydomain.net` against the Packet EIP.
- On AWS: You don't have to make any special DNS entries. Just make sure that the `grafana.ingress.host` value is `grafana.<CLUSTER NAME>.<AWS DNS ZONE>`.

## Steps

**NOTE**: This guide assumes that the underlying cloud platform is Packet and the OAuth provider is Github. For other OAuth providers the steps are the same just the secret parameters will change as mentioned in [Step 3](#step-3).

#### Step 1

- Create a Github OAuth application as documented in [Grafana docs](https://grafana.com/docs/grafana/latest/auth/github/).
- Set the **Homepage URL** to https://grafana.mydomain.net. This should be same as the `grafana.ingress.host` or `grafana.<CLUSTER NAME>.<AWS DNS ZONE>` as shown in [Step 2](#step-2).
- Set the **Authorization callback URL** to https://grafana.mydomain.net/login/github.
- Make a note of `Client ID` and `Client Secret` it will be needed in [Step 3](#step-3) and set as environment variable `GF_AUTH_GITHUB_CLIENT_ID` and `GF_AUTH_GITHUB_CLIENT_SECRET` respectively.

#### Step 2

Create `prometheus-operator.lokocfg` file with following contents:

```tf
component "prometheus-operator" {
  namespace = "monitoring"

  grafana {
    secret_env = var.grafana_secret_env
    ingress {
      host = "grafana.mydomain.net"
    }
  }
}
```

Observe the value of variable `secret_env` it should match the name of variable to be created in [Step 3](#step-3).

#### Step 3

Create `lokofg.vars` file or add following to an existing file, populate the values of this secret as needed:

```tf
grafana_secret_env = {
  "GF_AUTH_GITHUB_ENABLED" = "'true'"
  "GF_AUTH_GITHUB_ALLOW_SIGN_UP" = "'true'"
  "GF_AUTH_GITHUB_CLIENT_ID" = "YOUR_GITHUB_APP_CLIENT_ID"
  "GF_AUTH_GITHUB_CLIENT_SECRET" = "YOUR_GITHUB_APP_CLIENT_SECRET"
  "GF_AUTH_GITHUB_SCOPES" = "user:email,read:org"
  "GF_AUTH_GITHUB_AUTH_URL" = "https://github.com/login/oauth/authorize"
  "GF_AUTH_GITHUB_TOKEN_URL" = "https://github.com/login/oauth/access_token"
  "GF_AUTH_GITHUB_API_URL" = "https://api.github.com/user"
  "GF_AUTH_GITHUB_ALLOWED_ORGANIZATIONS" = "YOUR_GITHUB_ALLOWED_ORGANIZATIONS"
}
```

**NOTE**: In above configs the boolean value is set to `"'true'"` instead of plain `"true"` because Kubernetes expects the key value pair to be of type string and not boolean.

Modify the values of Github Auth configuration from

```ini
[auth.github]
enabled = true
client_id = YOUR_GITHUB_APP_CLIENT_ID
...
```

to look like following:

```yaml
"GF_AUTH_GITHUB_ENABLED" = "'true'"
"GF_AUTH_GITHUB_CLIENT_ID" = "YOUR_GITHUB_APP_CLIENT_ID"
```

The section name `[auth.github]` should be prepended with `GF_` and the name should be capitalised and `.` be replaced with `_`.

Deploy the prometheus operator using following command:

```bash
lokoctl component apply prometheus-operator
```

#### Step 4

Goto https://grafana.mydomain.net and now you will a special button **Sign in with GitHub**, use that to sign in with Github.

## What's next?

- Other auth providers for Grafana: https://grafana.com/docs/grafana/latest/auth/overview/#user-authentication-overview
