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

The fastest way to get started is to run `nks config init`:

```
bash-4.4$ nks config init
Creating config file...
Enter your NKS API Token: 
<NKS API TOKEN>
Setting default org..
Setting Provider Key...
Setting SSH Key...
bash-4.4$
```

## Usage 

Commands follow this structure: `nks <NOUN> <VERB> [FLAGS]`.
examples:
```
nks clusters list
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

### Resources
[NetApp/hci-nks-demo](https://github.com/NetApp/hci-nks-demo/tree/nks-shell/nks) - NKS tutorial with usage examples
