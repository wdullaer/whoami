# Whoami
A simple reimplementation of https://github.com/containous/whoami as a grpc service.

The idea is to make this available as a docker image to use as a dummy service to test your infra configuration with.
It serves no real other purposes.

## TODO
* Implement all of the methods
* Implement tracing (could be useful to see if your tracing provider works)
* Implement structured logging (could be useful to see if the logging infrastructure works)