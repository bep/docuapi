**DocuAPI** is a beautiful multilingual API documentation theme for [Hugo](http://gohugo.io/). This theme is built on top of the beautiful work of [Robert Lord](https://github.com/lord) and others on the [Slate](https://github.com/lord/slate) project ([Apache 2 License](https://github.com/lord/slate/blob/master/LICENSE)).

![Screenshot DocuAPI Excample site](https://raw.githubusercontent.com/bep/docuapi/master/images/screenshot.png)

## Use

See the [exampleSite](/exampleSite) and more specific its site [configuration](/exampleSite/config.toml) for the available options.

**Most notable:** This theme will use all the (non drafts) pages in the site and build a single-page API documentation. Using `weight` in the page front matter is the easiest way to control page order.

If you want a different page selection, please provide your own `layouts/index.html` template.

## Develop the Theme

**Note:** In most situations you will be well off just using the theme and maybe in some cases provide your own template(s). Please refer to the [Hugo Documentation](http://gohugo.io/overview/introduction/) for that.

But you may find styling issues, etc., that you want to fix. Pull requests of this type is welcomed.

**If you find issues that obviously belongs to  [Slate](https://github.com/lord/slate), then please report/fix them there, and we will pull in the latest changes here when that is done.**

This project provides a very custom asset bundler in [bundler.go](bundler.go) written in Go.

It depends on `libsass` to build, so you will need `gcc` (a C compiler) to build it for your platform. If that is present, you can try:

* `go get -u -v .`
* `go run bundler.go` (this will clone Slate to a temp folder)
* Alternative  to the above if you already have Slate cloned somewhere: `go run main.go -slate=/path/to/Slate`

With `make` available, you can get a fairly enjoyable live-reloading development experience for all artifacts by running:

* `hugo server` in your Hugo site project.
* `make serve` in the theme folder.




