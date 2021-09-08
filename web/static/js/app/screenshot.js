function TakeScreenshot(address) {
    Swal.fire({
        title: 'Processing screenshot...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    SendCommand(address, "screenshot")
        .then(response => {
            if (!response.ok) {
                throw Error(response.statusText);
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
            text: 'Error processing screenshot!',
            footer: err
        });
    });
}