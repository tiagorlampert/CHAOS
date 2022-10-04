async function OpenUrl(address) {
    Swal.fire({
        title: 'Type the URL to open',
        input: 'text',
        reverseButtons: true,
        showCancelButton: true,
        confirmButtonText: 'Open',
        showLoaderOnConfirm: true,
        preConfirm: (url) => {
            Swal.fire({
                title: 'Opening...',
                onBeforeOpen: () => {
                    Swal.showLoading()
                }
            });

            sendOpenUrl(address, url)
                .then(response => {
                    if (!response.ok) {
                        return response.text().then(err => {
                            throw new Error(err);
                        });
                    }
                    return response.text();
                })
                .then((result) => {
                    Swal.close();
                    Swal.fire({
                        text: 'URL opened successfully!',
                        icon: 'success'
                    });
                })
                .catch(err => {
                    console.log('Error: ', err);
                    Swal.fire({
                        icon: 'error',
                        title: 'Ops...',
                        text: 'Error opening URL!',
                        footer: err
                    });
                })
        },
        allowOutsideClick: () => !Swal.isLoading()
    })
}

async function sendOpenUrl(address, urlToOpen) {
    let formData = new FormData();
    formData.append('address', address);
    formData.append('url', urlToOpen);

    const url = '/open-url';
    const initDetails = {
        method: 'POST',
        body: formData,
        mode: "cors"
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}