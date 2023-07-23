# Webhook Bridge

This small web server serves as a bridge between different webhook dispatchers and receivers.

## Motivation

The motivation for this project came from a need of wanting to use flagsmith as a way to control feature flags in my gitops repo. Part of this requires that when a flag changes it needs to be able to trigger a workflow run in github to rebuild the manifests. Since flagsmith does not allow adding auth headers to their webhook configs I needed a way to bridge the gap to allow this functionality.

## Quick Start

```bash
docker run -d ghcr.io/graytonio/webhook-bridge:latest \
    -v ./config.yaml:/config.yaml \
    -p 8080:8080
```

## Config

Configuration of the webhooks are read in from a `config.yaml` file. An example configuration file is listed below.

```yaml
log_level: INFO
listen_address: 0.0.0.0:8080
webhooks:
  - path: my-service
    method: GET
    url: https://my-web-server/webhook
    headers:
        Authorization: Bearer my-token
```

### Top Level Config Description

| Config Key     | Description                                   | Default      |
|----------------|-----------------------------------------------|--------------|
| log_level      | Set the log level of the application          | INFO         |
| listen_address | What address the http server should listen on | 0.0.0.0:8080 |
| webhooks       | List of webhook handlers to setup             |              |

### Webhook Config Description

| Config Key | Description                                                                       | Default |
|------------|-----------------------------------------------------------------------------------|---------|
| path       | Path on the server that should trigger this webhook                               |         |
| method     | HTTP method that should be used to call the url                                   |         |
| url        | Webhook to call when this path is triggered                                       |         |
| headers    | A key: value map of the headers to attach to the headers of the triggered request |         |
