# terraform-provider-flowdock
terraform provider to manage flowdock users in some organizations

# WIKI Page
https://stuffnz.atlassian.net/wiki/spaces/KIWIOPS/pages/775192616/terraform+flowdock+provider+implementation

## Initialization
1. place the plugin executables in the user plugins directory: ~/.terraform.d/plugins
2. run terraform init

## build version
go version: 1.13.5 darwin/amd64 <br/>
terraform version: Terraform v0.11.14

# highly recommended tool
VSCode

## how to setup env to develop providers

- install go tools from https://golang.org/doc/install <br/>
- create or move the go project out of GOPATH (run go env to check you GOPATH) and make it a module, not a project<br/>
- run 'go get -d ./...' to install dependencies if needed ( or run go install ) <br/>
- This project used Go Modules, so you will need to enable them using `export GO113MODULE=on`, otherwise your go commands (run, build and test) will fail.

## how to build it(or cross compile)
go build -o terraform-provider-flowdock for MacOS <br/>
or GOOS=linux GOARCH=amd64 go build -v . for Linux

## how to debug it
export TF_LOG=DEBUG and then add log.Printf() in your code

## how to debug it in IDE console
copy the configure files from example, build it in the project's root folder

## how to do acceptance tests
1. export TF_ACC=1
2. export FLOWDOCK_TOKEN=you own token
3. https://github.com/stretchr/testify#installation
4. go test -v 

## how to do unit tests
we are now using testfy to run the unit test, reference is https://github.com/stretchr/testify#installation the process is:
1. export TF_ACC=1
2. export FLOWDOCK_TOKEN=you own token
3. go get github.com/stretchr/testify
4. go test -v 

## how to destroy a resource
Explicitly specifying the name of the resource you want to destroy is a good habit
terraform destroy -target=flowdock_invitation.i1

## how to import existing user info into invitation resource
run terraform import flowdock_invitation.instanceName userId_flowName_orgName(instanceId)
for example:
$ terraform import flowdock_invitation.richard_xue_1_fairfaxmedia 123456_kiwiops-projects_stuff-kiwiops-projects


## To delete a specific resource, run the following command:
terraform destroy -target=resource_type.resource_name

## resources
https://github.com/eddiezane/terraform-provider-todoist