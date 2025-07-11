import { EventsOn } from "@wailsio/runtime";

// Initialize map
const map = L.map("map").setView([20, 0], 2);
L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
  attribution: "&copy; OpenStreetMap contributors",
}).addTo(map);

// Marker groups
const vatsimLayer = L.layerGroup().addTo(map);
const realLayer = L.layerGroup().addTo(map);

const vatsimIcon = L.AwesomeMarkers.icon({
  icon: "plane",
  markerColor: "blue",
  prefix: "fa",
});

const realIcon = L.AwesomeMarkers.icon({
  icon: "plane",
  markerColor: "red",
  prefix: "fa",
});

L.control.layers(null, {
  "VATSIM Traffic": vatsimLayer,
  "Real Flights": realLayer,
}).addTo(map);

function updateVatsim(data) {
  vatsimLayer.clearLayers();
  data.forEach((a) => {
    const m = L.marker([a.Latitude, a.Longitude], {
      icon: vatsimIcon,
      rotationAngle: a.Heading || 0,
      rotationOrigin: "center center",
    }).bindPopup(`
      <b>${a.Callsign}</b><br>
      Alt: ${a.Altitude} ft<br>
      SPD: ${a.Groundspeed} kts<br>
      ${a.PlannedDep || ""} → ${a.PlannedDest || ""}
    `);
    m.addTo(vatsimLayer);
  });
}

function updateReal(data) {
  realLayer.clearLayers();
  data.forEach((f) => {
    const m = L.marker([f.Latitude, f.Longitude], {
      icon: realIcon,
      rotationAngle: f.Heading || 0,
      rotationOrigin: "center center",
    }).bindPopup(`
      <b>${f.Callsign}</b><br>
      Alt: ${Math.round(f.Altitude)} m<br>
      SPD: ${Math.round(f.Speed)} m/s
    `);
    m.addTo(realLayer);
  });
}

EventsOn("vatsim-data", updateVatsim);
EventsOn("real-data", updateReal);
