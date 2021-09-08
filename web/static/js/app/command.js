// Add event when `enter` key was pressed on input
window.addEventListener("load", function () {
    let input = document.getElementById("sendCommandInput");
    input.addEventListener("keyup", function (event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            document.getElementById("sendCommandBtn").click();
        }
    });
});

function sendingProgressBtn() {
    let sendCommandBtn = $("#sendCommandBtn");
    // Disable button
    sendCommandBtn.prop("disabled", true);
    // Add spinner to button
    sendCommandBtn.html(
        `<span class="spinner-border spinner-border-sm" id="sendSpan"></span> Sending`
    );
}

function defaultProgressBtn() {
    let sendCommandBtn = $("#sendCommandBtn");
    // Disable button
    sendCommandBtn.prop("disabled", false);
    // Remove progress
    let sendSpan = $("#sendSpan");
    sendSpan.remove();
    sendCommandBtn.text('Send');
}

function PerformCommand() {
    // Get query parameter
    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');

    // Get input value
    let sendCommandInput = document.getElementById('sendCommandInput');
    if (!sendCommandInput.value) {
        return
    }

    // Change button state
    sendingProgressBtn();

    SendCommand(address, sendCommandInput.value)
        .then(response => {
            if (response.status !== 200 && response.status !== 503) {
                console.log('Error! Status Code: ' + response.status);
                return;
            }
            return response.text();
        })
        .then(response => {
            // Add response to list
            let listGroup = $('.list-group');
            listGroup.append("<li class='border-bottom'><pre><br>" + response + "</pre></li>");
            listGroup.animate({scrollTop: listGroup.prop("scrollHeight")}, 400);
            sendCommandInput.value = "";

            // Reset send button to default state
            defaultProgressBtn();
        })
        .catch(err => {
            console.log('Error: ', err);
        });
}