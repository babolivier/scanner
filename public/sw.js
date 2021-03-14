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

// Implement a stub listener for fetch events to make the PWA installable.
// TODO: This will not work with Chrome 93 and later, but setting up an offline mode that
//  doesn't duplicate requests seems non-trivial so let's postpone that for later.
self.addEventListener("fetch", () => {});