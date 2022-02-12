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
                        return response.text().then(err => {
                            throw new Error(err);
                        });
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
                HandleError(err);
            });
        }
    });
}