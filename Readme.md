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

To get started either link the file server to another api to extend its features or upload and copy your access codes here.


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