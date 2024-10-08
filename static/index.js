
const DEFAULT_ZOOM_LEVEL = 9
const MIN_DRAG_TO_DISMISS = 200

function onReady(f) {
    if (document.readyState == "complete" || document.readyState == "interactive") {
        f()
    } else {
        document.addEventListener("DOMContentLoaded", f)
    }
}

async function getData() {
    return fetch("./data.json").then(r => r.json())
}

function setSidebarVisibility(bool) {
    if(bool) {
        document.getElementById("sidebar")?.classList.add("visible")
    } else {
        document.getElementById("sidebar")?.classList.remove("visible")
    }
    
}

function imageForDripId(id) {
    return "./images/" + id + ".png"
}


const dripDb = new Map()

function formatOffset(n) {
    return ("" + n / 1000).replace(".", ",")
}

function formatSide(side) {
    if( side == "L") {
        return "Li"
    }
    if( side == "R") {
        return "Re"
    }
    return ""
}

function setRoadStyle(element, style) {
    element.classList.remove("road_a","road_n","road_s")
    element.classList.add("road_"+style.toLowerCase());
}

function renderDripToSidebar(sidebarElement, drip) {
    const writeToElem = (elemClass,content) => sidebarElement.querySelector("."+elemClass).textContent = content
    const roadStyleElement = sidebarElement.querySelector(".hecto_road")
    const hectoElem = sidebarElement.querySelector(".hecto_container");

    writeToElem("drip_name",drip.name)
    writeToElem("drip_org",drip.organization)


    const childElements = sidebarElement.querySelectorAll(".drip_text_line")

    if(Array.isArray(drip.text)) {
        for (let index = 0; index < drip.text.length; index++) {
            if(index>3) break;
            childElements[index].innerText = drip.text[index]
        }
    } else {
        for (let index = 0; index < 4; index++) {
            childElements[index].innerText = ""
        }
    }

    // writeToElem("drip_text",drip.text ? drip.text.join() : "");
    

    if(drip.roadId != "" && drip.roadOffset >= 0) {
        hectoElem.style.display = "block"
        writeToElem("hecto_road",drip.roadId)
        writeToElem("hecto_bottom",formatOffset(drip.roadOffset))
        writeToElem("hecto_side",formatSide(drip.roadSide))
        setRoadStyle(roadStyleElement, drip.roadId[0])
    } else {
        hectoElem.style.display = "none"
    }

    


    const img = sidebarElement.querySelector("img")
    img.src = imageForDripId(drip.id)
}

function onMarkerClick(event,data) {
    const dripId = data[0].data.dripId
    const drip = dripDb.get(dripId)

    sidebar = document.getElementById("sidebar")

    renderDripToSidebar(sidebar, drip)

    setSidebarVisibility(true)
}

const imageFactorForZoomLevel = (lvl) => lvl / 18
const opacityForZoomLevel = (lvl) => .5 + lvl/40

const setIconSize = (markers, imageFactor) => {
    markers.forEach(marker => {
        const iconOptions = marker.options.icon.options
        const anchor = iconOptions.iconAnchor
        const origSize = iconOptions.iconSizeOrig
        const realSizes = iconOptions.iconSize

        if(!origSize) return;

        realSizes[0] = origSize[0] * imageFactor
        realSizes[1] = origSize[1] * imageFactor

        anchor[0] = realSizes[0] / 2
        anchor[1] = realSizes[1] / 2
    })
}

function setIconZoomEffect(zoomLevel, markerLayer, markers) {
    const opacity = opacityForZoomLevel(zoomLevel);
    const imageFactor = imageFactorForZoomLevel(zoomLevel);

    markerLayer.setOptions({opacity})
    setIconSize(markers, imageFactor)

    // For use in CSS
    markerLayer._map._container.dataset.zoomlevel = zoomLevel
}


function addSwipeListener(element) {
    if(!(element instanceof Element)) {
        throw new Error("Given argument is not an Element")
    }

    let touchStartY = 0

    function onStart(e) {
        touchStartY = e.targetTouches[0].clientY

        element.style.transition = "transform .015s"
    }

    function onMove(e) {
        const currentY = e.targetTouches[0].clientY
        element.style.transform = "translate3d(0,"+Math.max(0,currentY - touchStartY)+"px,0)"
    }

    function onEnd(e) {
        const currentY = e.changedTouches[0].clientY
        const difference = currentY-touchStartY
        
        element.style.transform = null
        element.style.transition = "transform .2s"

        if(difference > MIN_DRAG_TO_DISMISS) {
            element.classList.remove('visible')
        }
    }

    element.addEventListener("touchstart", onStart,{passive:true})
    element.addEventListener("touchmove", onMove,{passive:true})
    element.addEventListener("touchend", onEnd,{passive:true})
}

function setupSidebarSwipe() {
    const     sidebar = document.getElementById("sidebar");
    if(!sidebar) {
        throw new Error("Sidebar element not found")
    }

    addSwipeListener(sidebar)
}

onReady(() => {
    setupMap()
    setupSidebarSwipe()

    document.getElementById("close-button")?.addEventListener("click", () => setSidebarVisibility(false))
})

function createIcon(drip) {
    if(drip.imageWidth && drip.imageHeight) {
        const imgX = parseInt(drip.imageWidth)
        const imgY = parseInt(drip.imageHeight)

        return L.icon({
            iconUrl: imageForDripId(drip.id),
            iconSize: [imgX, imgY],
            iconSizeOrig: [imgX, imgY],
            iconAnchor: [imgX / 2, imgY / 2],
        })
    } 
    throw new Error("Drip doesn't have an image");

}

function setupMap() {
    const mapContainer = document.getElementById("map")
    if (!mapContainer) {
        throw new Error("Map element not found")
    }

    if (typeof L === "undefined") {
        throw new Error("Leaflet not found")
    }

    layerFactory(L) // Initalize leaflet plugin

    const map = L.map(mapContainer).setView([52.196665, 5.0811767], DEFAULT_ZOOM_LEVEL)
    const markerLayer = L.canvasIconLayer({
        opacity: opacityForZoomLevel(DEFAULT_ZOOM_LEVEL)
    }).addTo(map)
    const markers = []

    console.log("Leaflet instance:", map)

    markerLayer.addOnClickListener(onMarkerClick)

    map.addEventListener("zoom", (e) => {
        setIconZoomEffect(map.getZoom(), markerLayer, markers)
    })

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map)

    map.attributionControl.addAttribution('Data: <a href="http://opendata.ndw.nu/">opendata.ndw.nu/</a>')


    getData().then(d => {

        d.drips.forEach(drip => {
            dripDb.set(drip.id, drip)
        })

        d.drips.forEach(drip => {
            const lat = parseFloat(drip.lat, 10)
            const lon = parseFloat(drip.lon, 10)

            if (Number.isNaN(lat) || Number.isNaN(lon)) {
                return
            }

            
            if(drip.imageHeight == 0 || drip.imageWidth == 0 ) {
                return
            }
            const marker = L.marker([lat, lon], { icon:createIcon(drip) })
            
            marker.dripId = drip.id

            markers.push(marker)
        })


        setIconZoomEffect(map.getZoom(), markerLayer, markers)

        markerLayer.addLayers(markers)

        console.log("Added", markers.length, "displays to the map")
    })
}
