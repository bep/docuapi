[![Netlify Status](https://api.netlify.com/api/v1/badges/49111249-0a1a-4b5a-a3ab-45d00732fdb3/deploy-status)](https://app.netlify.com/sites/docuapi/deploys)

**DocuAPI** is a beautiful multilingual API documentation theme for [Hugo](http://gohugo.io/). This theme is built on top of the beautiful work of [Robert Lord](https://github.com/lord) and others on the [Slate](https://github.com/slatedocs/slate) project ([Apache 2 License](https://github.com/slatedocs/slate/blob/master/LICENSE)). The JS part has recently been rewritten from Jquery to [AlpineJS](https://alpinejs.dev/).

<br/>

> Visit the [demo site](https://docuapi.netlify.app/).

<br/>

![Screenshot DocuAPI Example site](https://raw.githubusercontent.com/bep/docuapi/master/images/screenshot.png)

## Use

Import the theme in your Hugo config:

```toml
[[module.imports]]
path = "github.com/bep/docuapi/v3"
```

v3 requires Hugo v0.161.0 or later. It builds its styles with Hugo's built-in `css.Build`, so there's no need for the extended Hugo version, Sass, PostCSS or Node.js.

Note, if you want the older Sass-based version, replace the path with `github.com/bep/docuapi/v2`, and for the older jQuery-version, `github.com/bep/docuapi`.

See the [exampleSite](https://github.com/bep/docuapi/tree/master/exampleSite) and more specific its site [configuration](https://github.com/bep/docuapi/blob/master/exampleSite/hugo.toml) for the available options.

**Most notable:** This theme will use all the (non drafts) pages in the site and build a single-page API documentation. Using `weight` in the page front matter is the easiest way to control page order.

If you want a different page selection, please provide your own `layouts/index.html` template.

## Styling

You can customize the look-and-feel by setting CSS variables in your site configuration. These take precedence over the theme defaults:

```toml
[params.docuapi.cssvars]
  nav-bg      = "#254e70"
  examples-bg = "#002642"
  font-family = '"Source Sans 3", sans-serif'
```

See [assets/css/docuapi/vars.css](https://github.com/bep/docuapi/blob/master/assets/css/docuapi/vars.css) for the available variables. You can also set them in your own `assets/css/docuapi/vars-custom.css`, which is empty in the theme. For more extensive changes, use `partials/hook_head_end.html` to add your own stylesheet.

The code highlighting styles are generated from your `markup.highlight` configuration, so make sure CSS classes are enabled (`noClasses = false`, the default) and pick a [style](https://gohugo.io/quick-reference/syntax-highlighting-styles/):

```toml
[markup.highlight]
  noClasses = false
  style     = "monokai"
```

## Hooks

You can override the layouts by providing some custom partials:

* `partials/hook_head_end.html` is inserted right before the `head` end tag. Useful for additional styles etc.
* `partials/hook_body_end.html` which should be clear by its name.
* `partials/hook_left_sidebar_start.html` the start of the left sidebar
* `partials/hook_left_sidebar_end.html` the end of the left sidebar
* `partials/hook_left_sidebar_logo.html` the logo `img` source

The styles and Javascript import are also put in each partial and as such can be overridden if really needed:

* `partials/styles.html`
* `partials/js.html`

## Stargazers over time

[![Stargazers over time](https://starchart.cc/bep/docuapi.svg)](https://starchart.cc/bep/docuapi)
