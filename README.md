# TWIPLA 3AS Go SDK

[![GoDoc](https://pkg.go.dev/badge/github.com/twipla/3as-go-sdk)](https://pkg.go.dev/github.com/twipla/3as-go-sdk)

A simple API wrapper for integrating the Analytics as a Service (3AS) APIs provided by TWIPLA.

## Getting started

1. [Create an RSA Key Pair (PEM format)](#creating-an-rsa-key-pair)
2. Send the resulting public key (`jwtRS256.key.pub`) to the TWIPLA Dev Team
3. [Install the library](#installation)
4. [Use the SDK instance](#how-to-use-the-library) to interact with the API

## Installation
```sh
go get github.com/twipla/3as-go-sdk
```

## How to use the library

Please refer to the example on [pkg.go.dev](https://pkg.go.dev/github.com/twipla/3as-go-sdk)

## Creating an RSA Key pair

1. Create the keypair: `ssh-keygen -t rsa -b 2048 -m PEM -f jwtRS256.key`
2. Convert the public key to PEM: `openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub`

## Concepts

### Terms

- **INTP (Integration Partner)**\
  The company that is integrating the analytics as a service solution (3AS)
- **STPs (Server Touchpoints)**\
  Credits used to measure data usage for a given website
- **Intpc (INTPC integration partner customer)**\
  One user of the INTP, can have many websites
- **Website**\
  The website where data will be tracked. It has a subscription with a package with a certain limit of STPs.
  This subscription can be upgraded or downgraded.
  When the website is created a tracking code snippet is returned that must be embedded within the websites HTML.
- **Package**\
  A package has a price and contains a certain number of STPs. They are used when upgrading/downgrading the subscription of a website.
