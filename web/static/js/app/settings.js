function RefreshToken() {
    Swal.fire({
        title: 'Are you sure?',
        text: "You must restart the server for changes to take effect, and all connected devices will be disconnected.",
        icon: 'question',
        showCancelButton: true,
        reverseButtons: true,
        confirmButtonColor: '#ffc107',
        confirmButtonText: 'Refresh',
        cancelButtonText: 'Cancel'
    }).then((result) => {
        if (result.value) {
            refreshToken()
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

                    let secretKeyInput = document.getElementById('secretKeyInput');
                    secretKeyInput.value = response;

                    ShowNotification('success', 'Success!', 'Secret key refreshed successfully.');
                })
                .catch(err => {
                    console.log('Error: ', err);
                    ShowNotification('danger', 'Ops!', 'Failed refreshing secret key.\n' + JSON.parse(err.message).error);
                });
        }
    });
}

async function refreshToken() {
    const url = '/settings/refresh-token';
    const initDetails = {
        method: 'GET',
        mode: "cors"
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}