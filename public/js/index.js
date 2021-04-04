const offlineMsg = "Pas de connexion";

function getPreview() {
    const btn = document.querySelector("#preview button");

    // Don't do anything if we're already getting a preview.
    if (btn.disabled) {
        return;
    }

    const spinner = document.querySelector("#preview-spinner");
    const tip = document.querySelector("#preview-tip");
    const img = document.querySelector("#preview-img");
    const errMsg = document.querySelector("#preview-err");

    // When waiting for a preview, only show the spinner, and don't allow asking for
    // another preview until the current one has been generated.
    btn.disabled = true;
    tip.classList.add("d-none");
    img.classList.add("d-none");
    errMsg.classList.add("d-none");
    spinner.classList.remove("d-none");

    function showErr() {
        // Show the error message and reset the button.
        spinner.classList.add("d-none");
        btn.disabled = false;
        errMsg.classList.remove("d-none");
    }

    // Request the preview.
    fetch("/preview.jpg")
        .then(response => {
            if (response.status === 200) {
                // Otherwise, if the request was a success, turn the image bytes
                // into a data URL.
                response.blob()
                    .then(dataURLForBlob)
                    .then(dataURL => {
                        // Hide the spinner and allow clicking the button again.
                        spinner.classList.add("d-none");
                        btn.disabled = false;
                        // Set the data URL as the source of the img element, and
                        // show it.
                        img.setAttribute("src", dataURL);
                        img.classList.remove("d-none");
                    })
            } else {
                // Show an user-readable error and log what actually went wrong.
                showErr();
                response.text().then(console.error);
            }
        })
        .catch((err) => {
            // Show an user-readable error and log what actually went wrong.
            showErr();
            console.error(err);
        });
}

function scan() {
    const btn = document.querySelector("#scan button");

    // Don't do anything if we're already getting a preview.
    if (btn.disabled) {
        return;
    }

    const formatSelect = document.querySelector("#scan select");
    const spinner = document.querySelector("#scan-spinner");
    const scanFormatErr = document.querySelector("#scan-format-err");
    const scanErr = document.querySelector("#scan-err");
    const scanSuccess = document.querySelector("#scan-success");
    const scanFilename = document.querySelector("#scan-filename")

    // When scanning, only show the spinner, and don't allow asking for another scan
    // until the current one has completed.
    btn.disabled = true;
    spinner.classList.remove("d-none");
    scanFormatErr.classList.add("d-none");
    scanErr.classList.add("d-none");
    scanSuccess.classList.add("d-none");

    function showElement(element) {
        // Show the given element and reset the button.
        spinner.classList.add("d-none");
        btn.disabled = false;
        element.classList.remove("d-none");
    }

    // Check that a format has been set.
    const format = formatSelect.value;
    if (format === "default") {
        showElement(scanFormatErr);
        return;
    }

    // Trigger the scan with the desired format.
    let url = "/scan?format=" + format;

    // If a rectangle has been drawn on top of the preview, only scan what's in it.
    const coords = rect.coords;
    if (coords !== null) {
        url += `&x=${Math.trunc(coords.x)}&y=${Math.trunc(coords.y)}`
        url += `&width=${coords.width}&height=${coords.height}`
    }

    fetch(url).
        then(response => {
            if (response.status === 200) {
                response.text()
                    .then(text => {
                        if (text === offlineMsg) {
                            // If the text from the response is the one sent by the
                            // service worker when the app is offline, show a
                            // user-readable error.
                            showElement(scanErr);
                        } else {
                            // Otherwise, if the request was a success, show the name of
                            // the newly generated file.
                            scanFilename.innerHTML = text;
                            showElement(scanSuccess);
                        }
                    })
            } else if (response.status === 400) {
                // If the response has a 400 status code, it means the format is either
                // missing or unsupported.
                showElement(scanFormatErr);
                response.text().then(console.error);
            } else {
                // Otherwise, show a standard user-readable error.
                showElement(scanErr);
                response.text().then(console.error);
            }
        })
        .catch((err) => {
            // Show an user-readable error and log what actually went wrong.
            showElement(scanErr);
            console.error(err);
        });
}

function dataURLForBlob(blob){
    // Generate a data URL from the given bytes, using the FileReader API.
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.onloadend = () => resolve(reader.result);
        reader.onerror = reject;
        reader.readAsDataURL(blob);
    })
}

// Register the event handlers.
document.querySelector("#preview button").onclick = getPreview;
document.querySelector("#scan button").onclick = scan;

// Register the service worker if supported by the browser.
if ('serviceWorker' in navigator) {
    navigator.serviceWorker.register('/sw.js', {scope: "."})
        .catch(err => {
            console.error('Failed to register service worker: ', err);
        })
}