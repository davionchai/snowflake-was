# Snowflake Warehouse Auto Scaling (Snowflake-WAS)

![Go Report Card](https://goreportcard.com/badge/github.com/davionchai/snowflake-was)
![License](https://img.shields.io/github/license/davionchai/snowflake-was)

`snowflake-was` is a Go written utility that helps to auto scale Snowflake warehouse size. It does so by monitoring warehouse activity and responding to different escalation points to upsize or downsize the warehouse size.

## Features

- Monitors query activity using Snowflake's [SHOW WAREHOUSES](https://docs.snowflake.com/en/sql-reference/sql/show-warehouses).
- Uses an escalation point mechanism to decide when to scale up or down.
- Written in Go for performance and concurrency.

## Getting Started

### Prerequisites

You will need the following tools:

- Go (version 1.19 or higher)
- A Snowflake account with `MODIFY` permission on specified warehouse

### Installation

Build the project:

```bash
go build
```

## Usage

After building the project, you can run the utility using the following command:

```bash
./snowflake-was
```

## Workflow

Refer to [workflow.md](./workflow.md)

## Configuring Snowflake-WAS

### Using yaml file
1. fill in the [config.yaml.template](./config.yaml.template) file and rename it to `config.yaml`

### Using environment variable
```bash
export WAS_SF_USERNAME=jdoe@email.com
export WAS_SF_ACCOUNT=hello.us-east-1
export WAS_SF_ROLE=wh_admin
export WAS_SF_WAREHOUSE_RUN=compute_wh_admin
export WAS_SF_AUTHENTICATOR=externalBrowser
export WAS_SF_WAREHOUSE_AUTOSCALE=compute_wh_analysts
export WAS_MIN_SIZE=xsmall
export WAS_MAX_SIZE=xxlarge
export WAS_QUEUED_THRESHOLD=5
export WAS_QUEUED_BASE_NUMBER=2
export WAS_DEFAULT_QUEUE_CHECKPOINT=5
export WAS_CYCLE_SECONDS=60
```

### Using cli
```bash
./snowflake-was was -h
```

Example

```bash
./snowflake-was was 
    --sf-username=jdoe@email.com \
    --sf-account=hello.us-east-1 \
    --sf-role=wh_admin \
    --sf-warehouse-run=compute_wh_admin \
    --sf-authenticator=externalBrowser \
    --sf-warehouse-autoscale=compute_wh_analysts \
    --min-size=xsmall \
    --max-size=xxlarge \
    --queued-threshold=5 \
    --queued-base-number=2 \
    --default-queue-checkpoint=5  \
    --cycle-seconds=60 \
```

## Known Issue
SSO reauthentication seems to be broken for snowflake go driver, meaning the application might throw error for every 4 hours. \
It's currently addressed in https://github.com/snowflakedb/gosnowflake/pull/836. \
Can implement [this approach](https://stackoverflow.com/questions/67069723/keep-retrying-a-function-in-golang) too.

## License

This project is licensed under the MIT License - see [LICENSE](./LICENSE) file for details.
