const img = document.getElementById("preview-img")

// Get the rectangle for the given element absolute, i.e. relative to the whole document,
// rather than to the viewport. We need to do this because getBoundingClientRect() returns
// a rectangle that's relative to the viewport and therefore its coordinates change if
// the page scrolls.
// We only change x and y (and not top/bottom/left/right) because we only use those when
// selecting a preview.
function getAbsoluteRect(el) {
    let r = el.getBoundingClientRect();
    r.x += window.pageXOffset || document.documentElement.scrollLeft;
    r.y += window.pageYOffset || document.documentElement.scrollTop;
    return r;
}

// A point in the rectangle.
class Point {
    constructor(x, y) {
        this.x = x
        this.y = y
    }
}

// The rectangle drawn on top of the preview image.
class PreviewRect {
    constructor() {
        this.el = document.getElementById("preview-rect");
        this.resetBtn = document.getElementById("preview-rect-reset")
        this.resetBtn.onclick = this.reset.bind(this);

        // Get the overlay divs, which are used to display a greyed out area where the
        // rectangle isn't drawn.
        this.overlays = document.getElementsByClassName("preview-rect-overlay")
        if (this.overlays.length !== 4) {
            console.error("Bad number of overlay elements")
        }

        // Set the initial state of the rectangle.
        this.reset()
    }

    // Returns the rectangle's coordinates, or null if the rectangle hasn't been drawn yet
    // or reset.
    get coords() {
        if (this.el.hidden) {
            return null
        }

        const previewRect = getAbsoluteRect(this.el)
        const imageRect = getAbsoluteRect(img)

        return {
            x: previewRect.x - imageRect.x,
            y: previewRect.y - imageRect.y,
            height: previewRect.height,
            width: previewRect.width,
        }
    }

    // Draw the rectangle by changing the div's style.
    draw() {
        const minX = Math.min(this.origin.x, this.cursor.x);
        const minY = Math.min(this.origin.y, this.cursor.y);

        const maxX = Math.max(this.origin.x, this.cursor.x);
        const maxY = Math.max(this.origin.y, this.cursor.y);

        this.el.hidden = false;

        this.el.style.left = minX + 'px';
        this.el.style.top = minY + 'px';
        this.el.style.width = maxX - minX + 'px';
        this.el.style.height = maxY - minY + 'px';

        this.drawOverlays()
    }

    // Draw the overlays that grey out the image where the rectangle isn't.
    drawOverlays() {
        const imageRect = getAbsoluteRect(img);
        const previewRect = getAbsoluteRect(this.el);

        const imageCoords = {
            minX: imageRect.x,
            maxX: imageRect.x + imageRect.width,
            minY: imageRect.y,
            maxY: imageRect.y + imageRect.height,
        }

        const previewCoords = {
            minX: previewRect.x,
            maxX: previewRect.x + previewRect.width,
            minY: previewRect.y,
            maxY: previewRect.y + previewRect.height,
        }

        // Top overlay
        if (previewCoords.minY !== imageCoords.minY) {
            this.overlays[0].hidden = false;
            this.overlays[0].style.left = previewCoords.minX + "px";
            this.overlays[0].style.top = imageCoords.minY + "px";
            this.overlays[0].style.width = previewRect.width + "px";
            this.overlays[0].style.height = previewCoords.minY - imageCoords.minY + "px";
        } else {
            this.overlays[0].hidden = true;
        }

        // Left overlay
        if (previewCoords.minX !== imageCoords.minX) {
            this.overlays[1].hidden = false;
            this.overlays[1].style.left = imageCoords.minX + "px";
            this.overlays[1].style.top = imageCoords.minY + "px";
            this.overlays[1].style.width = previewCoords.minX - imageCoords.minX + "px";
            this.overlays[1].style.height = imageRect.height + "px";
        } else {
            this.overlays[1].hidden = true;
        }

        // Right overlay
        if (previewCoords.maxX !== imageCoords.maxX) {
            this.overlays[2].hidden = false;
            this.overlays[2].style.left = previewCoords.maxX + "px";
            this.overlays[2].style.top = imageCoords.minY + "px";
            this.overlays[2].style.width = imageCoords.maxX - previewCoords.maxX + "px";
            this.overlays[2].style.height = imageRect.height + "px";
        } else {
            this.overlays[2].hidden = true;
        }

        // Bottom overlay
        if (previewCoords.maxY !== imageCoords.maxY) {
            this.overlays[3].hidden = false;
            this.overlays[3].style.left = previewCoords.minX + "px";
            this.overlays[3].style.top = previewCoords.minY + previewRect.height + "px";
            this.overlays[3].style.width = previewRect.width + "px";
            this.overlays[3].style.height = imageCoords.maxY - previewCoords.maxY + "px";
        } else {
            this.overlays[3].hidden = true;
        }
    }

    // Hide all of the overlays.
    hideOverlays() {
        for (let overlay of this.overlays) {
            overlay.hidden = true;
        }
    }

    // Mark the initial drawing of the rectangle as done. The initial drawing is done
    // between the first mouse key down and mouse key up events. If the cursor moves
    // during this time the points of the rectangle is set directly on the cursor's
    // coordinates. If it moves after the first mouse key up event, only the edge that's
    // closest to the cursor (if any) is moved.
    finishInitialDrawing() {
        this.drawn = true;
        this.resetBtn.classList.remove("d-none");
    }

    // Set the rectangle to its initial state.
    reset() {
        this.drawn = false;
        this.resetBtn.classList.add("d-none");
        this.origin = new Point(0, 0);
        this.cursor = new Point(0, 0);
        this.draw();
        this.el.hidden = true;
        this.hideOverlays();
    }

    // Set the coordinates for the cursor point.
    setCursor(x, y) {
        this.cursor = this.correctedPoint(x, y);
        this.draw();
    }

    // Instantiate a Point which coordinates have been corrected to fit inside the image
    // if necessary.
    correctedPoint(x, y) {
        const imgRect = getAbsoluteRect(img);

        if (x < imgRect.x && y < imgRect.y) {
            return new Point(imgRect.x, imgRect.y)
        }

        if (x < imgRect.x) {
            return new Point(imgRect.x, y)
        }

        if (y < imgRect.y) {
            return new Point(x, imgRect.y)
        }

        if (x > imgRect.x + imgRect.width && y > imgRect.y + imgRect.height) {
            return new Point(imgRect.x + imgRect.width, imgRect.y + imgRect.height)
        }

        if (x > imgRect.x + imgRect.width) {
            return new Point(imgRect.x + imgRect.width, y)
        }

        if (y > imgRect.y + imgRect.height) {
            return new Point(x, imgRect.y + imgRect.height)
        }

        return new Point(x, y)
    }

    // Move the edge of the rectangle that's the closest to the provided coordinates,
    // if any.
    moveClosestEdge(x, y) {
        const cursor = this.correctedPoint(x, y)
        const clientRect = getAbsoluteRect(this.el);

        // Calculate the distances between the cursor and each edge.
        const distances = {
            left: Math.abs(cursor.x - clientRect.x),
            right: Math.abs(clientRect.x + clientRect.width - cursor.x),
            top: Math.abs(cursor.y - clientRect.y),
            bottom: Math.abs(clientRect.y + clientRect.height - cursor.y),
        }

        // Iterate over the distances and store the edge with the shortest distance.
        let closestEdge = null;
        for (const edge of Object.keys(distances)) {
            // If an edge is more than 70px away from the cursor, ignore it to avoid that
            // edge jumping across the image; the expected UX is that the user instead
            // taps/clicks close to the edge and drags it to the end position.
            if (distances[edge] > 70) {
                continue;
            }

            if (closestEdge === null || distances[edge] < distances[closestEdge]) {
                closestEdge = edge
            }
        }

        // If no close enough edge has been found, don't do anything else.
        if (closestEdge === null) {
            return;
        }

        // Update the rectangle's coordinates to move the edge that's the closest.
        switch (closestEdge) {
            case "left":
                // Moving the left edge means moving the x value of the point with the
                // lowest one.
                if (this.origin.x < this.cursor.x) {
                    this.origin.x = cursor.x
                } else {
                    this.cursor.x = cursor.x
                }
                break
            case "right":
                // Moving the right edge means moving the x value of the point with the
                // highest one.
                if (this.cursor.x > this.origin.x) {
                    this.cursor.x = cursor.x
                } else {
                    this.origin.x = cursor.x
                }
                break
            case "top":
                // Moving the top edge means moving the y value of the point with the
                // lowest one.
                if (this.origin.y < this.cursor.y) {
                    this.origin.y = cursor.y
                } else {
                    this.cursor.y = cursor.y
                }
                break
            case "bottom":
                // Moving the bottom edge means moving the y value of the point with the
                // highest one.
                if (this.origin.y < this.cursor.y) {
                    this.cursor.y = cursor.y
                } else {
                    this.origin.y = cursor.y
                }
                break
        }

        // Redraw the rectangle.
        this.draw();
    }
}

let rect = new PreviewRect();

function onMouseDown(e) {
    if (!rect.drawn) {
        if (e.type === "mousedown") {
            rect.origin.x = e.pageX
            rect.origin.y = e.pageY
        } else if (e.type === "touchstart") {
            rect.origin.x = e.touches[0].pageX
            rect.origin.y = e.touches[0].pageY
        }
    }
}

registerListeners("mousedown", onMouseDown)
registerListeners("touchstart", onMouseDown)

function followCursor(e) {
    // Only process mousemove events if the left button is down.
    if ((e.type === "mousemove" && e.buttons === 1) || e.type === "touchmove") {
        let x, y;
        switch (e.type) {
            case "mousemove":
                x = e.pageX;
                y = e.pageY;
                break;
            case "touchmove":
                x = e.touches[0].pageX
                y = e.touches[0].pageY
                break;
        }

        // If the rectangle isn't in its initial drawing phase anymore, try to find the
        // edge that's closest to where the cursor is and move it to the cursor's
        // position. Otherwise, set the rectangle's diagonal.
        if (rect.drawn) {
            rect.moveClosestEdge(x, y)
        } else {
            rect.setCursor(x, y)
        }
    }
}

registerListeners("mousemove", followCursor)
registerListeners("touchmove", followCursor)

function onMouseUp() {
    rect.finishInitialDrawing();
}

registerListeners("mouseup", onMouseUp);
registerListeners("touchend", onMouseUp);

function registerListeners(evType, fn) {
    const property = "on" + evType;

    img[property] = fn;
    rect.el[property] = fn;
    for (const overlay of rect.overlays) {
        overlay[property] = fn;
    }
}