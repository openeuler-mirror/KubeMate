# Tutorial

## Install KubeMate plugins

For installing or building KubeMate plugins, see the installation instructions.

## Function Introduction

Three functions: backup, upgrade, and rollback.

### Backup

Backup can only backup to NFS server.

#### Prepare config file

```yaml
# /opt/kubemate/config/backup.yaml
nfs_server: ""
nfs_path: ""
```

> nfs_server： NFS server address，IP or domain
>
> nfs_path：NFS server storage path（absolute path）

#### How to use

```shell
# If you use openEuler
UniversalOS Backup
```

### Upgrade

#### Prepare config file
```yaml
# /opt/kubemate/config/upgrade.yaml
repo: |
```

#### How to use
```shell
# If you use openEuler
UniversalOS Upgrade
```

### Rollback

Rollback is provided using NFS+iPXE.

#### Prepare config file

```yaml
# /opt/kubemate/config/rollback.yaml
nfs_server: ""
nfs_path: ""
hostname: ""
user: ""
password: ""
ipxe_server: ""
ipxe_root_path: ""
ssh_port: ""
```

> nfs_server：NFS server address，IP or domain
>
> nfs_path：NFS server storage path（absolute path）
>
> hostname：hostname of the machine to be rollback
>
> user：iPXE server user
>
> password：iPXE server password
>
> ipxe_server：iPXE server address
>
> ipxe_root_path：IPXE server root path, used for clients to obtain configuration files through HTTP service
>
> ssh_port：IPXE server SSH service port, default to 22

#### How to use

```shell
# If you use openEuler
UniversalOS Rollback
```
