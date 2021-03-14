# Scanner

This project provides a simple UI to control a scanner.

![Screenshot](/screenshots/screenshot1.jpg)

It currently offers two features:

* displaying a preview from the scanner
* triggering a scan and uploading the resulting file to a WebDAV server

The motivation behind it is that my parents own an HP printer which can
also scan documents. While printing over the network works pretty well,
the only way they can _scan_ over the network is by using the proprietary
HP app, which is non-free, annoying to use and _requires_ the creation of
an account on their platform to interact using your _own_ printer on your
own network.

Given my parents are mostly using their tablets and smartphones while at
home, and they want to get rid of their 15-year-old desktop computer, which
main use nowadays is for scanning documents, I figured I'd write a simple
app they can use to trigger scans.

The project comes in two parts, a PWA (progressive web app) and a server.
The PWA lives on my parents' devices and interacts with the server when
they want to preview or scan a document. The server runs on a Raspberry Pi
4, which is connected to the printer via USB, and interacts with it using
[SANE](http://www.sane-project.org/).

Once a scan is complete, the resulting file is uploaded to my Nextcloud
instance using WebDAV. My parents can then use the Nextcloud app on their
devices to retrieve it.

This project has been built specifically for this use case. The app's UI
is entirely in French, and some features specific to HP printers might be
hardcoded in the code. So use it at your own risks.

Note that building this project requires the libsane development headers to
be installed on the system (which can be done  by installing the `libsane-devel`
package on Debian-based systems).