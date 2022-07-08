async function GenerateBinary() {
    let address = document.getElementById('address');
    let port = document.getElementById('port');
    let osTarget = document.getElementById('os_target');
    let filename = document.getElementById('filename');
    let runHidden = document.getElementById('run_hidden');

    if (!address.value || !osTarget.value) {
        ShowNotification('warning', 'Ops!', 'You should fill all the required fields.');
        return
    }

    Swal.fire({
        title: 'Building...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    generate(address.value, port.value, osTarget.value, filename.value, runHidden.checked)
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
        })
        .catch(err => {
            console.log('Error: ', err);
            Swal.close();
            ShowNotification('danger', 'Ops!', 'Failed building client binary.\n' + JSON.parse(err.message).error)
        });
}

async function generate(address, port, osTarget, filename, runHidden) {
    event.preventDefault();
    let formData = new FormData();
    formData.append('address', address);
    formData.append('port', port);
    formData.append('os_target', osTarget);
    formData.append('filename', filename);
    formData.append('run_hidden', runHidden);

    const url = '/generate';
    const initDetails = {
        method: 'POST',
        body: formData,
        mode: "cors",
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}