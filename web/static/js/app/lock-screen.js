function LockScreen(address) {
    Swal.fire({
        title: 'Are you sure?',
        text: "The device screen will be locked.",
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonText: 'Lock',
        cancelButtonText: 'Cancel'
    }).then((result) => {
        if (result.value) {
            Swal.fire({
                title: 'Processing command...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

            SendCommand(address, "lock")
                .then(response => {
                    if (!response.ok) {
                        throw Error(response.statusText);
                    }
                    return response.text();
                })
                .then(response => {
                    Swal.close();
                    Swal.fire({
                        text: 'Command send successfully!',
                        icon: 'success'
                    });
                }).catch(err => {
                console.log('Error: ', err);
                Swal.fire({
                    icon: 'error',
                    title: 'Ops...',
                    text: 'Error processing command!',
                    footer: err
                });
            });
        }
    });
}