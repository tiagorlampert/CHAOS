function ShowNotification(type, title, message) {
    $.notify({
        title: "<strong>" + title + "</strong>",
        message: message
    }, {
        type: type
    });
}