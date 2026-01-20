package handlers

// This is for testing or for when the user decides to just use the backend as a headless content bucket
// This 404 response is the default page that will be shown when the user doesn't set frontend mode
const Page_404 = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>Not Found</title>
    </head>
    <body>
        <h1 style="text-align: center;">404</h1>
        <div style="text-align: center; margin-bottom: 20px;">Not Found</div>
        <hr />
    </body>
</html>
>
`
