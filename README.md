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

### Templating
HTML templates should be referred to like importing packages (like `g.Render("blog/index", post)`). Templates are kept in the `html` directory for each plugin for modular-ness, but the `html/` should be left out when referring to them in code. File extension should also be ignored.

If desired, a template file called `_defs_.html.got` may be placed in each plugins `html/` directory. If one exists, it will be parsed along with each template file. `_defs_.html.got` should only contain template definitions, declared in Go's templating language as `{{define "name"}}…{{end}}`. One defs file can contain any arbitrary number of definitions. The definitions may be used in any other template file of the same package with `{{template "name"}}`.

### Dependencies
This package makes use of [gopgsqldriver](https://github.com/jbarham/gopgsqldriver) in conjunction with Go's `database/sql` package. It'll require a fixed version of jbarham's driver to work with the latest weekly build of Go. Attempting to build his package will make it obvious enough what to fix. There's a pull request for this fix already waiting.

Database support will probably come with gas instead of per plugin in the future.

## TODO

- All of the javascripts
	- Perhaps an animated background with <canvas>
- The new post page
- Pagination
- More links and stuff
- Make things use css image sprites more
- Tags
- Comments

To be continued...

## Comments, questions, concerns, hatemail
Contact moshee on Freenode or Rizon.
