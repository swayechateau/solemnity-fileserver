package templates

var UploadForm = `
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Upload Media</title>
    <style>
    body {
        background-repeat: no-repeat;
        background: #2980b9;
        /* fallback for old browsers */
        background: -webkit-linear-gradient(to right, #2c3e50, #2980b9);
        /* Chrome 10-25, Safari 5.1-6 */
        background: linear-gradient(to right, #2c3e50, #2980b9);
        /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */
        color: #fff;
        height: 100%;
        width: 100%;
    }

    .file-details {
        background: rgba(255, 255, 255, .4);
        padding: 25px;
        border-radius: 20px;
    }

    </style>
</head>

<body>
    <div class="file-details">
        <form enctype="multipart/form-data" action="/upload" method="post">
            <label for="">Public</label>
            <select name="public" id="">
                <option value="true">Yes</option>
                <option value="false">No</option>
            </select>
            <input type="file" name="media" accept="*" multiple>
            <button type="submit">Upload</button>
        </form>
    </div>
</body>

</html>
`
var Documentation = `
<html>

<head>
    <meta charset='utf-8'>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width">

    <title>File Server API</title>

    <!-- Flatdoc -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
    <script src='https://cdn.rawgit.com/rstacruz/flatdoc/v0.9.0/legacy.js'></script>
    <script src='https://cdn.rawgit.com/rstacruz/flatdoc/v0.9.0/flatdoc.js'></script>

    <!-- Flatdoc theme -->
	<link  href='https://cdn.rawgit.com/rstacruz/flatdoc/v0.9.0/theme-white/style.css' rel='stylesheet'>
	<script src='https://cdn.rawgit.com/rstacruz/flatdoc/v0.9.0/theme-white/script.js'></script>
  

    <!-- Initializer -->
    <script>
    Flatdoc.run({
        fetcher: Flatdoc.github('swayechateau/solemnity-fileserver')
    });
    </script>
</head>

<body role='flatdoc'>

    <div class='header'>
        <div class='left'>
            <h1>File Server API</h1>
            <ul>
                <li><a href='https://scm.kitechsoftware.com/kin/microservices/file-server'>View on GitLab</a></li>
                <li><a href='https://scm.kitechsoftware.com/kin/microservices/file-server/issues'>Issues</a></li>
            </ul>
        </div>
        <div class='right'>

        </div>
    </div>

    <div class='content-root'>
        <div class='menubar'>
            <div class='menu section' role='flatdoc-menu'></div>
        </div>
        <div role='flatdoc-content' class='content'></div>
    </div>

</body>

</html>
`
