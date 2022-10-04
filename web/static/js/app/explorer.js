// Add event when `enter` key was pressed on input
window.addEventListener("load", function () {
    let input = document.getElementById("pathInput");
    input.addEventListener("keyup", function (event) {
        if (event.keyCode === 13) {
            event.preventDefault();
            document.getElementById("refreshBtn").click();
        }
    });
});

function Home() {
    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');
    ExploreDirectory(address)
}

function Refresh() {
    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');

    let pathInput = document.getElementById('pathInput');
    if (!pathInput.value) {
        return
    }

    let encodedPath = btoa(encodeURI(pathInput.value));

    ExploreDirectory(address, encodedPath)
}

function OpenFolder(directory) {
    console.log("dir", directory)
    let urlParams = new URLSearchParams(window.location.search);
    let address = urlParams.get('address');
    let pathInput = document.getElementById('pathInput');
    if (!pathInput.value) {
        return
    }

    let newPath = pathInput.value + "/" + directory;

    console.log("new path", newPath)
    const encodedURI = encodeURI(newPath);
    let encodedPath = btoa(encodedURI);

    ExploreDirectory(address, encodedPath)
}

function ExploreDirectory(address, path) {
    Swal.fire({
        title: 'Loading...',
        onBeforeOpen: () => {
            Swal.showLoading()
        }
    });

    let uri = "/explorer?address=" + address
    if (path) {
        uri = uri.concat("&path=", path)
    }
    console.log("uri", uri)
    window.location.replace(uri);
}