# NKS CLI

A Command-Line Utility for NKS - https://nks.netapp.io

## Installation

The easiest way to install is by using `go get`:
```
go get -u gitlab.com/sgryczan/nks-cli/nks
```

## Configuration

nks-cli supports the following methods of configuration:
* Environment variables
* Config file

The fastest way to get started is to run `nks config init --set-defaults`:

```
bash-4.4$ nks config init --set-defaults
Creating config file...
Enter your NKS API Token: <TOKEN>
Setting default org..
Default SSH keys not set. Configuring...
bash-4.4$
```

The config file is located in `~/.nks.yaml`. This file contains a number of items pertaining to the context of the current request, including the default organization, cluster, and provider keysets.

To list the current configuration, run `nks config list`:

```
➜ ✗ nks config list
Current configuration:

api_url = https://api.nks.netapp.io
org_id = 21622
ssh_keyset = 58796
api_token = <TOKEN>
```

#### Example - Set a different default organization

```
➜ ✗ nks organization list
NAME             ID                  
Gentle Bird      21622 (current)     
HCI-DEV          24007               

➜ ✗ nks config set organization --id 24007

➜ ✗ nks organization list 
NAME             ID                  
Gentle Bird      21622               
HCI-DEV          24007 (current)     
```

## Usage Examples

Commands follow this structure: `nks <NOUN> <VERB> [FLAGS]`.
examples:

#### Example - List clusters
*note - the `--all` flag can be specified to list hidden clusters (i.e. service clusters*)
```
➜ ✗ nks clusters list --all
NAME                      ID        PROVIDER     NODES     K8s_VERSION     STATE                 
AWS-US-West-2a-1          26352     aws          3         v1.15.3         running (default)     
GCE-West1a-2              26350     gce          3         v1.15.3         running               
Crimson Snowflake         25743     hci          7         v1.14.3         running               
service-cluster-fluxt     25686     hci          4         v1.14.3         running
```

#### Example - Download KubeConfig for target cluster

`nks config set cluster` can be used to download the kubectl file for a target cluster. This can be handy for working with service clusters:

```
➜ ✗ nks clusters list --all          
NAME                      ID        PROVIDER     NODES     K8s_VERSION     STATE                 
AWS-US-West-2a-1          26352     aws          3         v1.15.3         running               
GCE-West1a-2              26350     gce          3         v1.15.3         running               
Crimson Snowflake         25743     hci          7         v1.14.3         running (default)     
service-cluster-fluxt     25686     hci          4         v1.14.3         running               

➜ ✗ nks config set cluster --id 25686

➜ ✗ kubectl get nodes                
NAME                      STATUS   ROLES    AGE   VERSION
net56aqmmv-master-1       Ready    master   18d   v1.14.3
net56aqmmv-pool-1-42g6r   Ready    <none>   18d   v1.14.3
net56aqmmv-pool-1-9gchl   Ready    <none>   18d   v1.14.3
net56aqmmv-pool-1-wtzhn   Ready    <none>   18d   v1.14.3
```

#### Example - Create a new cluster on HCI

```
➜  repos nks clusters create \
    --name mycluster \
    --provider hci \
    --num-workers 3 \
    --worker-size l \
    --master-size xl

NAME          ID        PROVIDER     NODES     K8s_VERSION     STATE     
mycluster     26529     hci          4         v1.14.3         draft
```


## Help

Most commands can be run with `-h` for a help menu:

```
bash-4.4$ nks -h
A command line utility for NKS

Usage:
  nks [command]

Available Commands:
  clusters     manage cluster resources
  config       nks cli configuration
  help         Help about any command
  keysets      add or edit keysets
  organization manage organizations
  repositories manage chart repositories
  solutions    mnanage solutions
  version      display CLI version

Flags:
      --config string   config file (default is $HOME/.nks.yaml)
  -h, --help            help for nks

Use "nks [command] --help" for more information about a command.
bash-4.4$
```

### Enable Shell Completion

#### Bash

Drop the following into your ~/.bashrc:

```
. /etc/profile.d/bash_completion.sh
. <(nks version --generatecompletion)
```

#### Zsh
Drop the following into your ~/.zshrc:

```
nks version -z > /home/nks/.oh-my-zsh/plugins/git/_nks
compinit
```

### Resources
[NetApp/hci-nks-demo](https://github.com/NetApp/hci-nks-demo/tree/nks-shell/nks) - NKS tutorial with usage examples
