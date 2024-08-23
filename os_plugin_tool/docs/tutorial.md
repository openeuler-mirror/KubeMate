# 使用方法

## 安装 KubeMate 插件

安装或编译 KubeMate 插件，请参考安装说明。

## 功能介绍

主要有备份、升级、回滚三个功能。

### 备份

备份功能目前仅可备份到 NFS 服务器。

#### 准备配置文件

```yaml
# /opt/kubemate/config/backup.yaml
nfs_server: ""
nfs_path: ""
```

> nfs_server： NFS 服务器地址，IP 或 域名
>
> nfs_path：NFS 服务器存储路径（绝对路径）

#### 调用方法

```shell
# 假设当前操作系统为 openEuler
UniversalOS Backup
```

### 升级

#### 准备配置文件
```yaml
# /opt/kubemate/config/upgrade.yaml
repo: |
```
#### 调用方法

```shell
# 假设当前操作系统为 openEuler
UniversalOS Upgrade
```

### 回滚

回滚功能目前使用 NFS+iPXE 方式提供。

#### 准备配置文件

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

> nfs_server：NFS 服务器地址，IP 或 域名
>
> nfs_path：NFS 服务器存储路径（绝对路径）
>
> hostname：待回滚机器的 hostname
>
> user：iPXE 服务器用户名
>
> password：iPXE 服务器密码
>
> ipxe_server：iPXE 服务器地址
>
> ipxe_root_path：iPXE 服务器根路径，用于客户端通过 HTTP 服务获取配置文件
>
> ssh_port：iPXE 服务器 ssh 服务端口，默认可填写 22

#### 调用方法

```shell
# 假设当前操作系统为 openEuler
UniversalOS Rollback
```
