# gas-blog

## What?
A plugin for my webapp server in progress, [gas](https://github.com/moshee/gas). This should be a sort of template for all plugins to come.

Everything you see here is very much a work in progress and probably does not work at the time of your reading this.

## Working ideas
All plugins should begin with its name, much like any Go package, and import the `gas` package.

```go
package blog

import "gas"
```

### Dependencies
This package makes use of [gopgsqldriver](https://github.com/jbarham/gopgsqldriver) in conjunction with Go's `database/sql` package. It'll require a fixed version of jbarham's driver to work with the latest weekly build of Go.

Database support will probably come with gas instead of per plugin in the future.

To be continued...

## Comments, questions, concerns, hatemail
Contact moshee on Freenode or Rizon.
