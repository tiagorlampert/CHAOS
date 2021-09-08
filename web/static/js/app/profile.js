function UpdatePassword() {
    let username = document.getElementById('inputUsername');
    let oldPassword = document.getElementById('inputOldPassword');
    let newPassword = document.getElementById('inputNewPassword');
    let confirmNewPassword = document.getElementById('inputConfirmNewPassword');

    if (!username.value || !oldPassword.value || !newPassword.value || !confirmNewPassword.value) {
        return
    }

    if (newPassword.value !== confirmNewPassword.value) {
        ShowNotification('danger', 'Ops!', 'Passwords didn\'t match. Try again.');
        return
    }

    update(username.value, oldPassword.value, newPassword.value)
        .then(response => {
            if (!response.ok) {
                return response.text().then(err => {
                    throw new Error(err);
                });
            }
            return response;
        })
        .then(response => {
            ShowNotification('success', 'Success!', 'User password updated successfully.');
        })
        .catch(err => {
            console.log('Error: ', err);
            ShowNotification('danger', 'Ops!', 'Failed updating user password.\n' + JSON.parse(err.message).error);
        });
}

async function update(username, oldPassword, newPassword) {
    let formData = new FormData();
    formData.append('username', username);
    formData.append('old-password', oldPassword);
    formData.append('new-password', newPassword);

    const url = '/user/password';
    const initDetails = {
        method: 'PUT',
        body: formData,
        mode: "cors"
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}