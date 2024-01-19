package main

var uploadForm = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bootstrap Drag and Drop File Upload</title>
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet">
    <style>
        .drag-area {
            border: 2px dashed #ddd;
            padding: 30px;
            text-align: center;
            margin: 10px;
        }
        .drag-area.highlight {
            border-color: #28a745;
            background-color: rgba(40, 167, 69, 0.2);
        }
    </style>
</head>
<body>
    <div class="container py-5">
        <h2 class="text-center mb-4">Bootstrap File Upload</h2>
        <form id="fileUploadForm" action="http://localhost:8080/upload" method="post" enctype="multipart/form-data">
            <div class="drag-area highlight" id="dragArea">
                <p class="mb-2">Drag and drop files here</p>
                <button type="button" class="btn btn-primary" id="buttonSelectFiles">Or select files</button>
                <input type="file" multiple hidden id="fileInput">
            </div>
        </form>
    </div>

    <script>
        const dragArea = document.getElementById('dragArea');
        const fileInput = document.getElementById('fileInput');
        const buttonSelectFiles = document.getElementById('buttonSelectFiles');

        buttonSelectFiles.addEventListener('click', () => fileInput.click());

        fileInput.addEventListener('change', (event) => {
            handleFiles(event.target.files);
        });

        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dragArea.addEventListener(eventName, preventDefaults, false);
        });

        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }

        ['dragenter', 'dragover'].forEach(eventName => {
            dragArea.addEventListener(eventName, highlight, false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dragArea.addEventListener(eventName, unhighlight, false);
        });

        function highlight() {
            dragArea.classList.add('highlight');
        }

        function unhighlight() {
            dragArea.classList.remove('highlight');
        }

        dragArea.addEventListener('drop', handleDrop, false);

        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            handleFiles(files);
        }

        function handleFiles(files) {
            // Here you can handle the file list, display thumbnails, or directly upload them using AJAX.
            console.log(files);
        }
    </script>
</body>
</html>
`
var documentation = `
<html>

<head>
    <meta charset='utf-8'>
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width">

    <title>File Server API</title>

    <!-- Flatdoc -->
    <script src="http://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js"></script>
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
