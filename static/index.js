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


onReady(() => {

    const mapContainer = document.getElementById("map");
    if(!mapContainer) {
        throw new Error("Map element not found")
    }
    
    var map = L.map(mapContainer).setView([52.196665, 5.0811767], 9);

    console.log("Leaflet instance:",map)

    mapContainer.dataset.zoomlevel = map.getZoom()

    map.addEventListener("zoom",(e) => {
        mapContainer.dataset.zoomlevel = map.getZoom()
    })

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    map.attributionControl.addAttribution('Data: <a href="http://opendata.ndw.nu/">opendata.ndw.nu/</a>');


    getData().then(d => {
        let dripcount = 0;

        d.drips.forEach(drip => {
            const lat = parseFloat(drip.lat,10)
            const lon = parseFloat(drip.lon,10)

            if(Number.isNaN(lat) || Number.isNaN(lon)) {
                return
            }

            const img = document.createElement("img")

            img.src = "./images/"+drip.id+".png"

            img.addEventListener("error", e => {
                console.warn("Image not found for display",drip)
                e.target.parentElement.style.display = "none"
            })

            img.classList.add("drip_img")


            const icon = L.divIcon({
                html: img
            })

            const marker = L.marker([lat,lon],{icon})
            marker.addTo(map)
            dripcount++
        })

        console.log("Added",dripcount,"displays to the map")
    })
})