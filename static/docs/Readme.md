FILE SERVER API
=======

The File Server API, only needs a html file input with the name of media for one file or media[] for multiple.

 * multi (boolean)
 * organisation (text)
 * owner (text)
 * type (text)
 * public (boolean)
 * date (date)
 * media (file) or media[] (files)
 
*Current version: [v0.1.0]*


Getting started
---------------

To get started eith link the file server to another api to extend its features or upload and copy your access codes here.

*In short: just download this file and upload it somewhere.*

The main JS and CSS files are also available in [view].

[Upload File >][uploader]

[uploader]: /upload
[view]: /view

### Via Browser

To upload via the browser, just create a form with ``enctype="multipart/form-data"`` and a file input named media.

``` html
<form enctype="multipart/form-data" action="http://fileserverhost.com/upload" method="post">

  <input type="file" name="media" accept="*">
            
  <button type="submit">Upload</button>
</form>
```

You may also use any or all of the additional fields.

``` html
<form enctype="multipart/form-data" action="http://fileserverhost.com/upload" method="post">
  <input type="text" name="multi" value="true" hidden>

  <label for="">Organisation</label>
  <input type="text" name="organisation" placeholder="organisation">

  <label for="">Owner</label>
  <input type="text" name="owner" placeholder="owner">

  <label for="">Type</label>
  <input type="text" name="type" placeholder="type">

  <label for="">Public</label>
  <select name="public" id="">
    <option value="1">Yes</option>
    <option value="0">No</option>
  </select>

  <label for="">Date</label>
  <input type="date" name="date" placeholder="date">

  <input type="file" name="media[]" accept="*" multiple>
            
  <button type="submit">Upload</button>
</form>
```


### Via a REST Client

You may also fetch a file. In this example, this fetches the file `Readme.md` in
the same folder as the HTML file.



Sample response.

``` json
{
  "access": {
    "organisation": "global",
    "owner": null,
    "type": null,
    "public": "1",
    "slug": "global3SFf0FE8AjahQ2Zy7HzI1",
    "share_code": "s1XXJ4WpNgcy",
    "access_code": "a1Xqzb5msNN9",
    "file_id": 1,
    "updated_at": "2020-09-07T21:51:00.000000Z",
    "created_at": "2020-09-07T21:51:00.000000Z",
    "id": 1
  },
  "file": {
    "name": "118747167_704659893448743_6996943649399192949_n.mp4",
    "extension": "mp4",
    "path": "\/_fileserver\/global\/video\/mp4\/2020\/Sep\/Mon",
    "full_path": "\/_fileserver\/global\/video\/mp4\/2020\/Sep\/Mon\/118747167_704659893448743_6996943649399192949_n.mp4",
    "size": 831358,
    "type": "video\/mp4",
    "updated_at": "2020-09-07T21:51:00.000000Z",
    "created_at": "2020-09-07T21:51:00.000000Z",
    "id": 1
  }
}
```

How it works
------------

Flatdoc is a hosted `.js` file (along with a theme and its assets) that you can
add into any page hosted anywhere.

#### All client-side

There are no build scripts or 3rd-party services involved. Everything is done in
the browser. Worried about performance? Oh, It's pretty fast.

Flatdoc utilizes the [GitHub API] to fetch your project's Readme files. You may
also configure it to fetch any arbitrary URL via AJAX.

#### Lightning-fast parsing

Next, it uses [marked], an extremely fast Markdown parser that has support for
GitHub flavored Markdown.

Flatdoc then simply renders *menu* and *content* DOM elements to your HTML
document. Flatdoc also comes with a default theme to style your page for you, or
you may opt to create your own styles.

Markdown extras
---------------

Flatdoc offers a few harmless, unobtrusive extras that come in handy in building
documentation sites.

#### Code highlighting

You can use Markdown code fences to make syntax-highlighted text. Simply
surround your text with three backticks. This works in GitHub as well.
See [GitHub Syntax Highlighting][fences] for more info.

    ``` html
    <strong>Hola, mundo</strong>
    ```

#### Blockquotes

Blockquotes show up as side figures. This is useful for providing side
information or non-code examples.

> Blockquotes are blocks that begin with `>`.

#### Smart quotes

Single quotes, double quotes, and double-hyphens are automatically replaced to
their typographically-accurate equivalent. This, of course, does not apply to
`<code>` and `<pre>` blocks to leave code alone.

> "From a certain point onward there is no longer any turning back. That is the
> point that must be reached."  
> --Franz Kafka

#### Buttons

If your link text has a `>` at the end (for instance: `Continue >`), they show
up as buttons.

> [View in GitHub >][project]

Customizing
===========

Basic
-----

### Theme options

For the default theme (*theme-white*), You can set theme options by adding
classes to the `<body>` element. The available options are:

#### big-h3
Makes 3rd-level headings bigger.

``` html
<body class='big-h3'>
```

#### no-literate
Disables "literate" mode, where code appears on the right and content text
appear on the left.

``` html
<body class='no-literate'>
```

#### large-brief
Makes the opening paragraph large.

``` html
<body class='large-brief'>
```

### Adding more markup

You have full control over the HTML file, just add markup wherever you see fit.
As long as you leave `role='flatdoc-content'` and `role='flatdoc-menu'` empty as
they are, you'll be fine.

Here are some ideas to get you started.

 * Add a CSS file to make your own CSS adjustments.
 * Add a 'Tweet' button on top.
 * Add Google Analytics.
 * Use CSS to style the IDs in menus (`#acknowledgements + p`).

### JavaScript hooks

Flatdoc emits the events `flatdoc:loading` and `flatdoc:ready` to help you make
custom behavior when the document loads.

``` js
$(document).on('flatdoc:ready', function() {
  // I don't like this section to appear
  $("#acknowledgements").remove();
});
```

Full customization
------------------

You don't have to be restricted to the given theme. Flatdoc is just really one
`.js` file that expects 2 HTML elements (for *menu* and *content*). Start with
the blank template and customize as you see fit.

[Get blank template >][template]

Misc
====

Inspirations
------------

The following projects have inspired Flatdoc.

 * [Backbone.js] - Jeremy's projects have always adopted this "one page
 documentation" approach which I really love.

 * [Docco] - Jeremy's Docco introduced me to the world of literate programming,
 and side-by-side documentation in general.

 * [Stripe] - Flatdoc took inspiration on the look of their API documentation.

 * [DocumentUp] - This service has the same idea but does a hosted readme 
 parsing approach.

Attributions
------------

[Photo](http://www.flickr.com/photos/doug88888/2953428679/) taken from Flickr,
licensed under Creative Commons.

Acknowledgements
----------------

© 2013, 2014, Rico Sta. Cruz. Released under the [MIT 
License](http://www.opensource.org/licenses/mit-license.php).

**Flatdoc** is authored and maintained by [Rico Sta. Cruz][rsc] with help from its 
[contributors][c].

 * [My website](http://ricostacruz.com) (ricostacruz.com)
 * [Github](http://github.com/rstacruz) (@rstacruz)
 * [Twitter](http://twitter.com/rstacruz) (@rstacruz)

[rsc]: http://ricostacruz.com
[c]:   http://github.com/rstacruz/flatdoc/contributors

[GitHub API]: http://github.com/api
[marked]: https://github.com/chjj/marked
[Backbone.js]: http://backbonejs.org
[dox]: https://github.com/visionmedia/dox
[Stripe]: https://stripe.com/docs/api
[Docco]: http://jashkenas.github.com/docco
[GitHub pages]: https://pages.github.com
[fences]:https://help.github.com/articles/github-flavored-markdown#syntax-highlighting
[DocumentUp]: http://documentup.com

[project]: https://github.com/rstacruz/flatdoc
[template]: https://github.com/rstacruz/flatdoc/raw/gh-pages/templates/template.html
[blank]: https://github.com/rstacruz/flatdoc/raw/gh-pages/templates/blank.html
[dist]: https://github.com/rstacruz/flatdoc/tree/gh-pages/v/0.9.0