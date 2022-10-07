function SelectFile() {
    let fileInput = document.getElementById('file-input');
    fileInput.click();
}

function UploadFile() {
    let file = document.getElementById("file-input").files[0];
    let formData = new FormData();
    formData.append("file", file);

    Swal.fire({
        title: 'Uploading ' + file.name + '...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    const url = '/upload';
    const initDetails = {
        method: 'post',
        body: formData,
        mode: "cors"
    }

    // Upload file to server
    fetch(url, initDetails).then(response => {
        if (!response.ok) {
            return response.text().then(err => {
                throw new Error(err);
            });
        }
        return response.text();
    }).then(response => {
        // Upload file to device
        SendToDevice(response);
    }).catch(err => {
        console.log('Error: ', err);
        Swal.fire({
            icon: 'error',
            title: 'Ops...',
            text: 'Error uploading file to server!',
            footer: err
        });
    });
}

function SendToDevice(filename) {
    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');
    let pathInput = document.getElementById('pathInput');
    let command = "upload";
    let filepath = pathInput.value + "/" + filename;

    // Say to device get file from server
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
            Swal.fire({
                text: 'File uploaded successfully!',
                icon: 'success'
            }).then(() => {
                Refresh();
            });
        }).catch(err => {
        console.log('Error: ', err);
        Swal.fire({
            icon: 'error',
            title: 'Ops...',
            text: 'Error uploading file to device!',
            footer: err
        });
    });
}