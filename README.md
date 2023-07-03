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

## Configuring Snowflake-WAS

🚧 WIP 🚧

## License

This project is licensed under the MIT License - see [LICENSE](./LICENSE) file for details.
