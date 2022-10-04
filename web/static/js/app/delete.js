function DeleteFile(filename) {
    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');
    let pathInput = document.getElementById('pathInput');
    let filepath = "delete " + pathInput.value + "/" + filename;

    Swal.fire({
        title: 'Are you sure?',
        text: "The file '" + filename + "' will be deleted permanently.",
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonColor: '#d64130',
        confirmButtonText: 'Delete',
        cancelButtonText: 'Cancel'
    }).then((result) => {
        if (result.value) {
            Swal.fire({
                title: 'Deleting ' + filename + '...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

            SendCommand(address, filepath)
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
                    Swal.fire({
                        text: 'Deleted successfully!',
                        icon: 'success'
                    }).then(() => {
                        Refresh();
                    });
                }).catch(err => {
                console.log('Error: ', err);
                Swal.fire({
                    icon: 'error',
                    title: 'Ops...',
                    text: 'Error deleting file!',
                    footer: err
                });
            });
        }
    });
}