# Go Starter Template

OmniViews's production-ready microservice template

## Goal
The goal of this template is to serve as the foundation for new projects, it is not expected to be a one-size-fits-all type of solution, opinions are kept at a minimal level by design. Before adding new stuff to the project, please ask yourself this question "is the feature generic enough that will be used in the majority of use cases?"

## Features

- Chi Router as default HTTP router, configured with a set of middlewares and some helper utilities
- Koanf is used for loading configuration files
- Uber's zap is configured to be the logging utility
- Prometheus' metrics available at /metrics
- K8s health endpoints

## Getting started

- Clone this repository
- Create a new git repository
- Change the remote repository with the following command

```
git remote set-url origin git@git.omnicloud.mx:omnicloud/development/your_new_shiny_repo.git
```

## Folder structure

- /cmd - Home of the application entry points, ex. main.go
- /config - Application configuration loader and type safe config structs
- /internal - Application sources for internal use
- /internal/adapters - API controllers, Event consumers, etc.
- /internal/repositories - Domain models and repositories live here
- /internal/services - The glue between adapters and repositories
- /mocks - Useful mocks for unit testing
- /resources - External resources, such as, db migrations, yaml files, etc.

## Customizations

This project contains some placeholders that need to be updated prior to start adding new code. Key places that require updates are:

- go.mod (update module name)

## Linter

[Golangci-linter](https://golangci-lint.run/) is used as our primary linter tool, this project has already a set of linters enabled with some sane defaults, please go to the golangci-lint website to learn how to install and integrate with your IDE of choice.

On demand checks can be triggered by typing the following command at your project root.

```
golangci-lint run
```

## Examples

The project contains some examples to help you get going, you can run the project with the following command.

```
go run cmd/main.go
```

Next, open your browser and hit the following URLs:

```
http://localhost:8080/v1/ok
http://localhost:8080/v1/error
```
