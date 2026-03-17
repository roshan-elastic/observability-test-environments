# Contributing

## Build and Test

Oblt-cli is a command line tool develop in Go, the requirements to compile it are Go 1.15+ and make installed in the machine. All the generated files are copied to a `build` folder inside the project.

```
git clone git@github.com:elastic/observability-test-environments.git
cd observability-test-environments
git checkout ccs_deployment_2
cd tools/oblt-cli
make clean lint build test
```

To install the tool in your Go bin folder you can use `go install`

*NOTE* tests are a work in progress, there are 2 that failed, even though. the feature works (manual tested).

## Add a new command

If you'd like to extend this tool, please use below standard:

`oblt-cli [entity]+ [verb]? [flags]*`
