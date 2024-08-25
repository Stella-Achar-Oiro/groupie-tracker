const apiBaseUrl = '/api';
let artists = [];
let locations = [];
let dates = [];
let relations = [];

const artistGrid = document.getElementById('artist-grid');
const modal = document.getElementById('artist-modal');
const modalContent = document.getElementById('modal-content');
const closeBtn = document.getElementsByClassName('close')[0];
const searchInput = document.getElementById('search-input');

async function fetchData() {
    try {
        const response = await axios.get(`${apiBaseUrl}/data`);
        const data = response.data;
        artists = data.artists;
        locations = data.locations;
        dates = data.dates;
        relations = data.relations;

        renderArtists(artists);
    } catch (error) {
        console.error('Error fetching data:', error);
    }
}

function renderArtists(artistsToRender) {
    artistGrid.innerHTML = '';
    artistsToRender.forEach((artist, index) => {
        const artistCard = document.createElement('div');
        artistCard.className = 'artist-card';
        artistCard.style.animationDelay = `${index * 0.1}s`;
        artistCard.innerHTML = `
            <img src="${artist.image}" alt="${artist.name}" class="artist-image">
            <div class="artist-info">
                <div class="artist-name">${artist.name}</div>
                <div class="artist-details">
                    <p>Creation: ${artist.creationDate}</p>
                    <p>First Album: ${artist.firstAlbum}</p>
                </div>
            </div>
        `;
        artistCard.addEventListener('click', () => showArtistDetails(artist));
        artistGrid.appendChild(artistCard);
    });
}

function showArtistDetails(artist) {
    const artistRelation = relations[0].index.find(r => r.id === artist.id);
    modalContent.innerHTML = `
        <h2>${artist.name}</h2>
        <img src="${artist.image}" alt="${artist.name}" style="width: 200px; height: 200px; object-fit: cover; border-radius: 10px; box-shadow: 0 5px 15px rgba(0,0,0,0.3);">
        <p><strong>Members:</strong> ${artist.members.join(', ')}</p>
        <p><strong>Creation Date:</strong> ${artist.creationDate}</p>
        <p><strong>First Album:</strong> ${artist.firstAlbum}</p>
        <h3>Concert Locations and Dates</h3>
        <ul id="concert-list"></ul>
    `;

    const concertList = document.getElementById('concert-list');
    for (const [location, dates] of Object.entries(artistRelation.datesLocations)) {
        const listItem = document.createElement('li');
        listItem.innerHTML = `<strong>${location}:</strong> ${dates.join(', ')}`;
        concertList.appendChild(listItem);
    }

    modal.style.display = 'block';
    initMap(artist, artistRelation.datesLocations);
}

function initMap(artist, datesLocations) {
    if (window.map) {
        window.map.remove();
    }
    
    window.map = L.map('map').setView([20, 0], 2);
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors'
    }).addTo(window.map);

    const locations = Object.keys(datesLocations);
    const geocodePromises = locations.map(location => {
        const [city, country] = location.split(', ');
        return fetch(`https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(city)},${encodeURIComponent(country)}`)
            .then(response => response.json())
            .then(data => {
                if (data.length > 0) {
                    const { lat, lon } = data[0];
                    return { location, lat, lon, dates: datesLocations[location] };
                }
                return null;
            });
    });

    Promise.all(geocodePromises)
        .then(results => {
            results.filter(result => result !== null).forEach(({ location, lat, lon, dates }) => {
                L.marker([lat, lon]).addTo(window.map)
                    .bindPopup(`<b>${location}</b><br>${dates.join(', ')}`);
            });
        })
        .catch(error => console.error('Error fetching location data:', error));
}

closeBtn.onclick = () => {
    modal.style.display = 'none';
};

window.onclick = (event) => {
    if (event.target == modal) {
        modal.style.display = 'none';
    }
};

searchInput.addEventListener('input', async (e) => {
    const searchTerm = e.target.value.toLowerCase();
    if (searchTerm.length > 0) {
        try {
            const response = await axios.get(`${apiBaseUrl}/search?q=${searchTerm}`);
            renderArtists(response.data);
        } catch (error) {
            console.error('Error searching artists:', error);
        }
    } else {
        renderArtists(artists);
    }
});

fetchData();