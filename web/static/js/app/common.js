async function SendCommand(address, command) {
    let formData = new FormData();
    formData.append('address', address);
    formData.append('command', command);

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