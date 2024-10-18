# Simple Temporal Workflow

This project is a part of temporal exploration.
We aim to create a simple workflow to register new company to the HubSpot using temporal.

## How To Run
### 1. Install temporal server on local
```
brew install temporal
```
Once installed, we can start temporal server with command:
```
temporal server start-dev
```
More about how to setup a local temporal service : https://learn.temporal.io/getting_started/go/dev_environment/#set-up-a-local-temporal-service-for-development-with-temporal-cli
### 2. Adjust `.env` file
copy `.env-example` to `.env` file and adjust the value
### 3. Start worker
```
go run worker/main.go
```
### 4. Create new workflow
```
go run starter/main.go
```