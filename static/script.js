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
        const response = await fetch('/api/data');
        const data = await response.json();
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
    const fragment = document.createDocumentFragment();
    artistsToRender.forEach((artist, index) => {
        const artistCard = document.createElement('div');
        artistCard.className = 'artist-card';
        artistCard.style.animationDelay = `${index * 0.05}s`;
        artistCard.innerHTML = `
            <img src="${artist.image}" alt="${artist.name}" class="artist-image" loading="lazy">
            <div class="artist-info">
                <div class="artist-name">${artist.name}</div>
                <div class="artist-details">
                    <p>Creation: ${artist.creationDate}</p>
                    <p>First Album: ${artist.firstAlbum}</p>
                </div>
            </div>
        `;
        artistCard.addEventListener('click', () => showArtistDetails(artist));
        fragment.appendChild(artistCard);
    });
    artistGrid.appendChild(fragment);
}

function showArtistDetails(artist) {
    const artistRelation = relations[0].index.find(r => r.id === artist.id);
    modalContent.innerHTML = `
        <h2>${artist.name}</h2>
        <img src="${artist.image}" alt="${artist.name}" style="width: 200px; height: 200px; object-fit: cover; border-radius: 10px; box-shadow: 0 5px 15px rgba(0,0,0,0.3);">
        <p><strong>Members:</strong> ${artist.members.join(', ')}</p>
        <p><strong>Creation Date:</strong> ${artist.creationDate}</p>
        <p><strong>First Album:</strong> ${artist.firstAlbum}</p>
        <h3>Concert Information</h3>
        <p>
            <a href="#" id="show-locations">Locations</a> | 
            <a href="#" id="show-dates">Dates</a> |
            <a href="#" id="show-relations">Relations</a>
        </p>
        <div id="concert-info"></div>
    `;

    document.getElementById('show-locations').addEventListener('click', (e) => {
        e.preventDefault();
        showLocations(artistRelation.datesLocations);
    });

    document.getElementById('show-dates').addEventListener('click', (e) => {
        e.preventDefault();
        showDates(artistRelation.datesLocations);
    });

    document.getElementById('show-relations').addEventListener('click', (e) => {
        e.preventDefault();
        showRelations(artistRelation.datesLocations);
    });

    showRelations(artistRelation.datesLocations);

    modal.style.display = 'block';
    initMap(artist, artistRelation.datesLocations);
}

function showRelations(datesLocations) {
    const concertInfo = document.getElementById('concert-info');
    const fragment = document.createDocumentFragment();
    const ul = document.createElement('ul');
    ul.id = 'concert-list';
    for (const [location, dates] of Object.entries(datesLocations)) {
        const li = document.createElement('li');
        li.innerHTML = `<strong>${location}:</strong> ${dates.join(', ')}`;
        ul.appendChild(li);
    }
    fragment.appendChild(ul);
    concertInfo.innerHTML = '';
    concertInfo.appendChild(fragment);
}

function showLocations(datesLocations) {
    const concertInfo = document.getElementById('concert-info');
    concertInfo.innerHTML = '<ul>' + 
        Object.keys(datesLocations).map(location => `<li>${location}</li>`).join('') +
    '</ul>';
}

function showDates(datesLocations) {
    const concertInfo = document.getElementById('concert-info');
    const allDates = Object.values(datesLocations).flat();
    concertInfo.innerHTML = '<ul>' + 
        allDates.map(date => `<li>${date}</li>`).join('') +
    '</ul>';
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

const debounce = (func, delay) => {
    let timeoutId;
    return (...args) => {
        clearTimeout(timeoutId);
        timeoutId = setTimeout(() => func.apply(null, args), delay);
    };
};

searchInput.addEventListener('input', debounce(async (e) => {
    const searchTerm = e.target.value.toLowerCase();
    if (searchTerm.length > 0) {
        const filteredArtists = artists.filter(artist => 
            artist.name.toLowerCase().includes(searchTerm) ||
            artist.members.some(member => member.toLowerCase().includes(searchTerm))
        );
        renderArtists(filteredArtists);
    } else {
        renderArtists(artists);
    }
}, 300));

closeBtn.onclick = () => {
    modal.style.display = 'none';
};

window.onclick = (event) => {
    if (event.target == modal) {
        modal.style.display = 'none';
    }
};

fetchData();