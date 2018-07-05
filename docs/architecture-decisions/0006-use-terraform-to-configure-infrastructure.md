# 6. Use Terraform to Configure Infrastructure

Date: 2018-07-03

## Status

Accepted

## Context

Having a repeatable, predictable way to create and change infrastructure is essential to stability and reliability of applications. One leading candidate to allow writing infrastructure as code is [Terraform](https://www.terraform.io).

Other tools to consider might be Ansible, Puppet or Chef, but they are less suited to modern cloud infrastructure than Terraform as they were developed when running code on VMs was the norm. They are good options for Configuration Management, but less appropriate for managing the infrastructure itself.

Amazon CloudFormation is a closer a possibility, but it is only usable on the Amazon stack whereas Terraform can be used to manage any Cloud. This flexibility will allow us to manage infrastructure which span clouds which this project may require (such as backend processing on AWS and frontend APIs on Heroku).

Both CloudFormation and Terraform are good choices for Infrastructure Orchestration, but the AWS only restriction of CloudFormation makes it much less compelling to adopt.

See [3. Follow Twelve Factor methodology](0003-follow-twelve-factor-methodology.md)

## Decision

We will use Terraform to configure our Infrastructure.

## Consequences

We will have a repeatable, predictable way to create and change infrastructure.

Staff will need to learn Terraform.

Developers and Operations will be able to propose and review changes to Infrastructure prior to changes being made.

The same code review processes we use to ensure better software can be used to allow us to better understand our infrastructure as well as see exactly what changes are being proposed and why.
