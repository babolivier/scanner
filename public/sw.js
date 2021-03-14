const cacheName = "bzh.abolivier.scanner.v1";
const cachedFiles = [
    "/css/bootstrap.min.css",
    "/css/index.css",
    "/js/bootstrap.bundle.min.js",
    "/js/jquery-3.6.0.min.js",
    "/js/index.js",
    "/index.html",
    "/misc/offline-msg.txt"
];

// Cache static files on install.
self.addEventListener("install", event => {
    event.waitUntil(
        caches.open(cacheName)
            .then(cache => {
                return cache.addAll(cachedFiles);
            })
    )
});

// Implement support for a very basic offline mode.
self.addEventListener("fetch", event => {
    const url = new URL(event.request.url)

    event.respondWith(
        new Promise((resolve, reject) => {
            if (url.pathname.endsWith("preview.jpg") || url.pathname.endsWith("scan")) {
                // Never cache results from the server's endpoints as these are expected
                // to change for each request.
                resolve(fetch(event.request))
            }

            // Otherwise, try to fetch the resouce from the cache, or from the network if
            // the cache yielded no result.
            caches.match(event.request)
                .then(response => {
                    if (response) {
                        resolve(response);
                    }
                    resolve(fetch(event.request));
                }).catch((err) => {
                    console.log(1)
                    console.log(event.request.url)
                    return caches.match("misc/offline-msg.txt");
                });
        }).catch(() => {
            console.log(2)
            console.log(event.request.url)
            // Serve a basic string to tell the app we're offline.
            return caches.match("misc/offline-msg.txt")
        })
    )
});