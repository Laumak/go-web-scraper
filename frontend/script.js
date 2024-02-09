const mapsScript = document.createElement("script");
mapsScript.src = `https://maps.googleapis.com/maps/api/js?key=${MAPS_API_KEY}&map_ids=${MAP_ID}&callback=initMap`;
document.getElementsByTagName("body")[0].appendChild(mapsScript);

function initMap() {
  map = new google.maps.Map(document.getElementById("map"), {
    center: { lat: 60.22962504913607, lng: 24.932131085949063 },
    zoom: 12,
    mapId: MAP_ID,
  });

  const condos = [
    {
      title: "Viputie 11",
      content: "asd 1",
      lat: 60.23216947591634,
      lng: 24.718600640875724,
    },
    {
      title: "Juvanpuistonkuja 2",
      content: "asd 2",
      lat: 60.2756474,
      lng: 24.75063556394958,
    },
  ];

  condos.forEach(({ title, content, lat, lng }) => {
    // Set up info window for marker click
    const infowindow = new google.maps.InfoWindow({
      content: content,
      ariaLabel: title,
    });

    // Create a marker on the map
    const marker = new google.maps.Marker({
      position: { lat, lng },
      map,
      title,
    });

    // Open the info window when associated marker is clicked
    marker.addListener("click", () => {
      infowindow.open({
        anchor: marker,
        map,
      });
    });
  });
}
