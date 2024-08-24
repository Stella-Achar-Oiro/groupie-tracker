// Step 1: Initialize communication with the platform
var platform = new H.service.Platform({
    apikey: "UxDC5OOIhOC4BEmxOmwUKAa8q413CJS311fmOzt3UeE"
});
var defaultLayers = platform.createDefaultLayers();

// Step 2: Initialize a map - centered over Europe
var map = new H.Map(document.getElementById('map'),
    defaultLayers.vector.normal.map, {
    center: { lat: 20, lng: 10 },
    zoom: 2, // Increased zoom level for better initial visibility
    pixelRatio: window.devicePixelRatio || 1
});

// Add a resize listener to make sure that the map occupies the whole container
window.addEventListener('resize', () => map.getViewPort().resize());

// Step 3: Make the map interactive
var behavior = new H.mapevents.Behavior(new H.mapevents.MapEvents(map));

// Create the default UI components
var ui = H.ui.UI.createDefault(map, defaultLayers);

// Step 4: Read locations and add markers
window.onload = function () {
    let locations = document.getElementsByClassName('cities');
    Array.from(locations).forEach(location => {
        let cityName = location.textContent.trim();
        fetch(`https://geocoder.ls.hereapi.com/6.2/geocode.json?searchtext=${encodeURIComponent(cityName)}&gen=9&apiKey=UxDC5OOIhOC4BEmxOmwUKAa8q413CJS311fmOzt3UeE`)
            .then(response => response.json())
            .then(data => {
                let position = data.Response.View[0].Result[0].Location.DisplayPosition;
                let marker = new H.map.Marker({
                    lat: position.Latitude,
                    lng: position.Longitude
                });
                map.addObject(marker);
            })
            .catch(error => console.error(`Error fetching geolocation for ${cityName}:`, error));
    });
}
