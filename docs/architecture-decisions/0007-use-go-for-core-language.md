# 7. Use Go for Core Language

Date: 2018-07-05

## Status

Accepted

## Context

The choice of which programming language to use is governed by a host of different factors. Most languages can be made to work for whatever task is required, though some are better at certain tasks than others.

We expect the nature of the work in this project to benefit from concurrency, so choosing a language with good support for this is important. Since we will be deploying to AWS Lambda (See [5. Use AWS Lambda](0005-use-aws-lambda.md)), we are further limited to using one of the supported languages. Other considerations include ease of packaging and distribution, excellent data streaming abilities, and a healthy ecosystem of 3rd party libraries.

## Decision

Use Go for the core application language.

## Consequences

Go is a new language for us, so there is some inherent risk involved in this choice. If we decide we are not making progress quickly enough due to lack of language familiarity, we can fall back to Python.
