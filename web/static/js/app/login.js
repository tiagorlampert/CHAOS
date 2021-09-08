// Add event when `enter` key was pressed on username
window.addEventListener("load", function () {
    let input = document.getElementById("inputUsername");
    input.addEventListener("keyup", function (event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            document.getElementById("loginButton").click();
        }
    });
});

// Add event when `enter` key was pressed on password
window.addEventListener("load", function () {
    let input = document.getElementById("inputPassword");
    input.addEventListener("keyup", function (event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            document.getElementById("loginButton").click();
        }
    });
});

function Login() {
    let username = document.getElementById('inputUsername');
    let password = document.getElementById('inputPassword');
    if (!username.value || !password.value) {
        return
    }

    auth(username.value, password.value)
        .then(response => {
            if (response.status !== 200 && response.status !== 401) {
                console.log('Error! Status Code: ' + response.status);
                return;
            }
            return response.json();
        })
        .then(response => {
            window.location.href = "/";
        })
        .catch(err => {
            console.log('Error: ', err);
            ShowNotification('danger', 'Ops!', 'Invalid username or password.');
        });
}

async function auth(username, password) {
    let formData = new FormData();
    formData.append('username', username);
    formData.append('password', password);

    const url = '/auth';
    const initDetails = {
        method: 'post',
        body: formData,
        mode: "cors"
    }

    let response = await fetch(url, initDetails);
    let data = await response;
    return data;
}