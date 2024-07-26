# goji

A pluggable application framework for Go.

<div align="center">
  <img src="./logo.png" width="256px" />
</div>

## Overview

Goji is a lightweight, modular and composable application framework that provides the basic building blocks for building applications.

Each of the components is considered a DI module that can be composed into the application using the `fx` framework introduced earlier.

## Components

Goji provides the following components:

- [Database](./core/database): Database driver, powered by [ent](https://entgo.io/). It provides a custom database driver which supports multiple read replicas for use with the `ent` ORM.

- [Redis](./core/redis): Redis client, powered by [go-redis](https://github.com/redis/go-redis). It provides a redis connection manager which manages multiple connections to different redis instances.

- [Storage](./core/storage): Object storage client, supporting any S3-compatible storage provider using the [minio sdk](https://github.com/minio/minio-go), as well as a Google Cloud Storage client using the [google cloud storage sdk](https://pkg.go.dev/cloud.google.com/go/storage).

- [Queue](./core/queue): Queue client, supporting any Google Cloud PubSub queue provider using the [google cloud pubsub sdk](https://pkg.go.dev/cloud.google.com/go/pubsub).

- [Email](./core/email): Email client, supporting any SMTP email provider, as well as [Mailgun](https://www.mailgun.com).

- [Notification](./core/notification): Notification client, currently supporting [Slack](https://slack.com) notifications only.

- [Vault](./core/vault): Secret storage client, supporting [HashiCorp Vault](https://www.vaultproject.io), [Google Secret Manager](https://cloud.google.com/secret-manager), [Infisical](https://infisical.com) and a simple redis-based secret storage.

- [Crypt](./core/crypt): Symmetric encryption service, supporting AES encryption. It supports the Vault secret storage for dynamic key management and rotation.
