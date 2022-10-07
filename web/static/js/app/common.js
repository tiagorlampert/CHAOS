async function SendCommand(address, command, parameter) {
    let formData = new FormData();
    formData.append('address', address);
    formData.append('command', command);
    formData.append('parameter', parameter);

    const url = '/command';
    const initDetails = {
        method: 'post',
        body: formData,
        mode: "cors"
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}

function HandleError(err){
    if (err.message === "unsupported platform") {
        Swal.fire({
            icon: 'warning',
            title: 'Ops...',
            text: 'Error processing command!',
            footer: err
        });
    } else {
        Swal.fire({
            icon: 'error',
            title: 'Ops...',
            text: 'Error processing command!',
            footer: err
        });
    }
}