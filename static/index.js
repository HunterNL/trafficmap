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

const imageFactorForZoomLevel = (lvl) => lvl / 18
const opacityForZoomLevel = (lvl) => .5 + lvl/40

const setIconSize = (markers, imageFactor) => {
    markers.forEach(marker => {
        const iconOptions = marker.options.icon.options
        const anchor = iconOptions.iconAnchor
        const origSize = iconOptions.iconSizeOrig
        const realSizes = iconOptions.iconSize

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


const DEFAULT_ZOOM_LEVEL = 9

onReady(() => {
    const mapContainer = document.getElementById("map");
    if(!mapContainer) {
        throw new Error("Map element not found")
    }

    if(typeof L === "undefined") {
        throw new Error("Leaflet not found")
    }

    layerFactory(L) // Initalize leaflet plugin
    
    const map = L.map(mapContainer).setView([52.196665, 5.0811767], DEFAULT_ZOOM_LEVEL);
    const markerLayer = L.canvasIconLayer({
        opacity: opacityForZoomLevel(DEFAULT_ZOOM_LEVEL)
    }).addTo(map)
    const markers = [];

    console.log("Leaflet instance:",map)

    map.addEventListener("zoom",(e) => {
        setIconZoomEffect(map.getZoom(), markerLayer, markers)
    })

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    map.attributionControl.addAttribution('Data: <a href="http://opendata.ndw.nu/">opendata.ndw.nu/</a>');

    getData().then(d => {

        d.drips.forEach(drip => {
            const lat = parseFloat(drip.lat,10)
            const lon = parseFloat(drip.lon,10)

            const imgX = parseInt(drip.imageWidth)
            const imgY = parseInt(drip.imageHeight)

            if(Number.isNaN(lat) || Number.isNaN(lon)) {
                return
            }

            const icon = L.icon({
                iconUrl: "./images/"+drip.id+".png",
                iconSize: [imgX,imgY],
                iconSizeOrig: [imgX,imgY],
                iconAnchor: [imgX/2,imgY/2],
                // html: img
            })

            const marker = L.marker([lat,lon],{icon})

            markers.push(marker)            
        })


        setIconZoomEffect(map.getZoom(), markerLayer, markers)

        markerLayer.addLayers(markers)

        console.log("Added",markers.length,"displays to the map")
    })
})