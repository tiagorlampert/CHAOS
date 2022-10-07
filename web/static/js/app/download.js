function DownloadFile(filename) {
    Swal.fire({
        title: 'Downloading ' + filename + '...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');
    let pathInput = document.getElementById('pathInput');
    let command = "download";
    let filepath = pathInput.value + "/" + filename;

    SendCommand(address, command, filepath)
        .then(response => {
            if (!response.ok) {
                return response.text().then(err => {
                    throw new Error(err);
                });
            }
            return response.text();
        })
        .then(response => {
            Swal.close();
            window.location.href = 'download/' + response;
        }).catch(err => {
        console.log('Error: ', err);
        Swal.fire({
            icon: 'error',
            title: 'Ops...',
            text: 'Error downloading file!',
            footer: err
        });
    });
}