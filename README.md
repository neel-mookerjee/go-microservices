# What is it?
This is a simple set of microservices written in GO and deployed on K8

It reads from DynamoDB and returns the objects as HTTP response

## Makefile
```
internal-services go/compile       compile go programs
internal-services docker/tag/list  list the existing tagged images
internal-services docker/build     build and tag the Docker image. vars:tag
internal-services docker/push      push the Docker image to ECR. vars:tag
internal-services helm/install     Deploy stack into kubernetes. vars: stack
internal-services helm/delete      delete stack from reference. vars: stack
internal-services deploy           Compiles, builds and deploys a stack for a tag. vars: tag, stack
internal-services redeploy         Compiles, builds and deploys a stack for a tag. vars: tag, stack
internal-services help             this helps
```
