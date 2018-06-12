Aws Working Sugar
===

AWS Working Sugar provides a set of useful tools for
	your day to day AWS duties.


## Install

~~~bash
go get github.com/Dal-Papa/awsugar
~~~

## Usage

### Options

```
  -d, --dry-run         Toggle a list-only mode without executing any action.
  -h, --help            help for awsugar
  -r, --region string   Choose the region to execute the actions in (default "us-west-2")
```
## awsugar clean

Clean your AWS account in various places

### Synopsis

Clean your AWS account in various places including:
	
	- Soft kill an EC2 instance with a snapshot first
	- Remove deprecated ELB without target instances
	- Remove available volumes and snapshot them
	- Release unattached Elastic IPs and Network Interfaces
	- Remove unused Security Groups
	- Remove unused Launch Configurations

```
awsugar clean [type] [flags]
```

### Options

```
  -h, --help          help for clean
  -s, --sweet-clean   allow some preparation before cleaning (snapshot, etc.) (default true)
```

## awsugar search

Search through various AWS services

### Synopsis

Provides some helpers to search through services in AWS.
	
	Allows to search for an IP in Route53.

```
awsugar search [type] [flags]
```

### Options

```
  -h, --help         help for search
      --ip ipSlice   list of IPs to search (default [])
```

### Options inherited from parent commands

```
  -d, --dry-run         Toggle a list-only mode without executing any action.
  -r, --region string   Choose the region to execute the actions in (default "us-west-2")
```
