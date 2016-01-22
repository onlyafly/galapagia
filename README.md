# galapagia

An artificial evolution simulator, for simulating concepts of evolution on artificial life.

## Setup

### Install all NPM packages

    $ sudo npm install

## Developing

### Starting the app

Run gulp, which will compile web assets, build and run the server, and automatically rebuild and restart the server on changes.

    $ gulp

### Adding a new dependency

**Note: The only time you need to use godep is when you add a new dependency to an application or update an existing dependency that is already vendored in your application.**

1. Use go-get to get the dependency:

       $ go get -u github.com/russross/blackfriday

2. Add a relevant import in a Go file:

       import "github.com/russross/blackfriday"

3. Vendor the additional import (records it, copies it to Godeps/, rewrites imports):

       $ godep save -r ./...

### Updating an existing dependency

    $ go get -u <dependency>

    $ godep update <dependency>
