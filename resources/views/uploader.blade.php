<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Document</title>
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
        <form enctype="multipart/form-data" action="upload" method="post">
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
            <input type="text" name="multi" value="true" hidden>
            <button type="submit">Upload</button>
        </form>
    </div>


</body>

</html>