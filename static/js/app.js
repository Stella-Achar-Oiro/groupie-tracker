function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

mapboxgl.accessToken = 'pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw';

const searchInput = document.getElementById('search-input');
const suggestionsContainer = document.getElementById('suggestions');
const creationYearSlider = document.getElementById('creation-year');
const creationYearDisplay = document.getElementById('creation-year-display');
const firstAlbumYearSlider = document.getElementById('first-album-year');
const firstAlbumYearDisplay = document.getElementById('first-album-year-display');
const memberCheckboxes = document.getElementById('member-checkboxes');
const locationCheckboxes = document.getElementById('location-checkboxes');
let map;
let allArtists = [];
let allLocations = new Set();

searchInput.addEventListener('input', () => {
    const query = searchInput.value;
    if (query.length > 1) {
        fetch(`/api/suggestions?q=${encodeURIComponent(query)}`)
            .then(response => response.json())
            .then(suggestions => displaySuggestions(suggestions))
            .catch(error => console.error('Error:', error));
    } else {
        suggestionsContainer.innerHTML = '';
    }
});

creationYearSlider.addEventListener('input', updateCreationYearDisplay);
firstAlbumYearSlider.addEventListener('input', updateFirstAlbumYearDisplay);

function updateCreationYearDisplay() {
    creationYearDisplay.textContent = creationYearSlider.value;
    applyFilters();
}

function updateFirstAlbumYearDisplay() {
    firstAlbumYearDisplay.textContent = firstAlbumYearSlider.value;
    applyFilters();
}

function displaySuggestions(suggestions) {
    suggestionsContainer.innerHTML = '';
    suggestions.forEach(suggestion => {
        const div = document.createElement('div');
        div.className = 'suggestion-item';
        div.textContent = `${suggestion.text} (${suggestion.type})`;
        div.onclick = () => {
            searchInput.value = suggestion.text;
            suggestionsContainer.innerHTML = '';
            searchArtists(suggestion.text);
        };
        suggestionsContainer.appendChild(div);
    });
}

function searchArtists(query) {
    showLoading();
    const filters = getFilterValues();
    fetch(`/api/search?q=${encodeURIComponent(query)}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(filters)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        displayResults(data.artists);
        hideLoading();
    })
    .catch(error => {
        console.error('Error:', error);
        showError('An error occurred while searching for artists. Please try again later.');
        hideLoading();
    });
}

function getFilterValues() {
    return {
        creationYearMin: parseInt(document.getElementById('creation-year').value),
        creationYearMax: 2023,
        firstAlbumYearMin: parseInt(document.getElementById('first-album-year').value),
        firstAlbumYearMax: 2023,
        members: Array.from(document.querySelectorAll('#member-checkboxes input:checked')).map(cb => parseInt(cb.value)),
        locations: Array.from(document.querySelectorAll('#location-checkboxes input:checked')).map(cb => cb.value)
    };
}

function initializeFilters() {
    const earliestCreationYear = Math.min(...allArtists.map(artist => artist.creationDate));
    const latestCreationYear = Math.max(...allArtists.map(artist => artist.creationDate));
    const earliestAlbumYear = Math.min(...allArtists.map(artist => parseInt(artist.firstAlbum.split('-')[2])));
    const latestAlbumYear = Math.max(...allArtists.map(artist => parseInt(artist.firstAlbum.split('-')[2])));
    
    creationYearSlider.min = earliestCreationYear;
    creationYearSlider.max = latestCreationYear;
    creationYearSlider.value = earliestCreationYear;
    creationYearDisplay.textContent = earliestCreationYear;

    firstAlbumYearSlider.min = earliestAlbumYear;
    firstAlbumYearSlider.max = latestAlbumYear;
    firstAlbumYearSlider.value = earliestAlbumYear;
    firstAlbumYearDisplay.textContent = earliestAlbumYear;

    const maxMembers = Math.max(...allArtists.map(artist => artist.members.length));
    memberCheckboxes.innerHTML = '';
    for (let i = 1; i <= maxMembers; i++) {
        const label = document.createElement('label');
        label.innerHTML = `<input type="checkbox" value="${i}"> ${i}`;
        memberCheckboxes.appendChild(label);
    }

    allLocations = new Set();
    allArtists.forEach(artist => {
        artist.locations.forEach(location => {
            allLocations.add(location.split('-')[0].trim());
        });
    });

    locationCheckboxes.innerHTML = '';
    Array.from(allLocations).sort().forEach(location => {
        const label = document.createElement('label');
        label.innerHTML = `<input type="checkbox" value="${location}"> ${location}`;
        locationCheckboxes.appendChild(label);
    });

    // Add event listeners to checkboxes
    memberCheckboxes.querySelectorAll('input').forEach(checkbox => {
        checkbox.addEventListener('change', applyFilters);
    });
    locationCheckboxes.querySelectorAll('input').forEach(checkbox => {
        checkbox.addEventListener('change', applyFilters);
    });
}

function applyFilters() {
    const creationYear = parseInt(creationYearSlider.value);
    const firstAlbumYear = parseInt(firstAlbumYearSlider.value);
    const selectedMembers = Array.from(memberCheckboxes.querySelectorAll('input:checked')).map(cb => parseInt(cb.value));
    const selectedLocations = Array.from(locationCheckboxes.querySelectorAll('input:checked')).map(cb => cb.value);

    const filteredArtists = allArtists.filter(artist => {
        const artistAlbumYear = parseInt(artist.firstAlbum.split('-')[2]);
        const artistLocations = artist.locations.map(loc => loc.split('-')[0].trim());
        return artist.creationDate >= creationYear &&
               artistAlbumYear >= firstAlbumYear &&
               (selectedMembers.length === 0 || selectedMembers.includes(artist.members.length)) &&
               (selectedLocations.length === 0 || artistLocations.some(loc => selectedLocations.includes(loc)));
    });

    displayResults(filteredArtists);
}

function displayResults(artists) {
    const container = document.getElementById('results-container');
    container.innerHTML = '';
    artists.forEach(artist => {
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.innerHTML = `
            <img src="placeholder.jpg" data-src="${artist.image}" alt="${artist.name}" class="lazy-image">
            <h3>${artist.name}</h3>
            <p><i class="fas fa-calendar-alt"></i> Created: ${artist.creationDate}</p>
            <p><i class="fas fa-compact-disc"></i> First Album: ${artist.firstAlbum}</p>
        `;
        card.onclick = () => fetchArtistDetails(artist.id);
        container.appendChild(card);
    });
    lazyLoadImages();
}

function lazyLoadImages() {
    const images = document.querySelectorAll('.lazy-image');
    const options = {
        root: null,
        rootMargin: '0px',
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.classList.remove('lazy-image');
                observer.unobserve(img);
            }
        });
    }, options);

    images.forEach(img => observer.observe(img));
}

function transitionToPage(callback) {
    const mainContent = document.querySelector('.container');
    mainContent.classList.add('page-transition', 'fade-out');
    setTimeout(() => {
        callback();
        mainContent.classList.remove('fade-out');
    }, 300);
}

// Use this when fetching artist details
function fetchArtistDetails(id) {
    showLoading();
    transitionToPage(() => {
        fetch(`/api/artist/${id}`)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            displayArtistDetails(data);
            hideLoading();
        })
        .catch(error => {
            console.error('Error:', error);
            showError('An error occurred while fetching artist details. Please try again later.');
            hideLoading();
        });
    });
}

let favorites = JSON.parse(localStorage.getItem('favorites')) || [];

function toggleFavorite(artistId) {
    const index = favorites.indexOf(artistId);
    if (index === -1) {
        favorites.push(artistId);
    } else {
        favorites.splice(index, 1);
    }
    localStorage.setItem('favorites', JSON.stringify(favorites));
    updateFavoriteButton(artistId);
}

function updateFavoriteButton(artistId) {
    const button = document.getElementById(`favorite-${artistId}`);
    if (button) {
        button.innerHTML = favorites.includes(artistId) 
            ? '<i class="fas fa-star"></i> Remove from Favorites' 
            : '<i class="far fa-star"></i> Add to Favorites';
    }
}

function shareArtist(artist) {
    if (navigator.share) {
        navigator.share({
            title: artist.name,
            text: `Check out ${artist.name} on Groupie Tracker!`,
            url: window.location.href
        }).then(() => console.log('Successful share'))
        .catch((error) => console.log('Error sharing:', error));
    } else {
        alert('Web Share API is not supported in your browser. You can copy the URL to share.');
    }
}

function displayArtistDetails(details) {
    const container = document.getElementById('artist-details');
    container.style.display = 'block';
    
    console.log('Artist details:', details);

    container.innerHTML = `
        <h2>${details.artist.name}</h2>
        <img src="${details.artist.image}" alt="${details.artist.name}" style="max-width: 300px;">
        <p><i class="fas fa-users"></i> Members: ${details.artist.members.join(', ')}</p>
        <p><i class="fas fa-calendar-alt"></i> Creation Date: ${details.artist.creationDate}</p>
        <p><i class="fas fa-compact-disc"></i> First Album: ${details.artist.firstAlbum}</p>
        <h3><i class="fas fa-map-marker-alt"></i> Locations:</h3>
        <ul>${details.locations.map(loc => `<li>${loc.address}</li>`).join('')}</ul>
        <h3><i class="fas fa-calendar-check"></i> Dates:</h3>
        <ul>${details.dates.map(date => `<li>${date}</li>`).join('')}</ul>
        <h3><i class="fas fa-link"></i> Relations:</h3>
        <ul>${Object.entries(details.relations).map(([loc, dates]) => `
            <li>${loc}: ${dates.join(', ')}</li>
        `).join('')}</ul>
        <div class="action-buttons">
            <button id="favorite-${details.artist.id}" onclick="toggleFavorite(${details.artist.id})">
                ${favorites.includes(details.artist.id) ? '<i class="fas fa-star"></i> Remove from Favorites' : '<i class="far fa-star"></i> Add to Favorites'}
            </button>
            <button onclick="shareArtist(${JSON.stringify(details.artist)})">
                <i class="fas fa-share-alt"></i> Share
            </button>
        </div>
    `;

    console.log('Locations for map:', details.locations);
    displayMap(details.locations);
}
function displayMap(locations, concertDates) {
    console.log('Locations received:', locations);

    if (map) {
        map.remove();
    }

    map = new mapboxgl.Map({
        container: 'map',
        style: 'mapbox://styles/mapbox/streets-v11',
        center: [0, 0],
        zoom: 1
    });

    const bounds = new mapboxgl.LngLatBounds();

    locations.forEach(location => {
        console.log('Creating marker for:', location);
        
        if (!location.lon || !location.lat) {
            console.error('Invalid coordinates for location:', location);
            return;
        }

        const el = document.createElement('div');
        el.className = 'custom-marker';
        el.innerHTML = '<i class="fas fa-map-marker-alt"></i>';
        el.style.color = '#FF0000';
        el.style.fontSize = '24px';

        el.addEventListener('click', () => {
            console.log('Clicked location:', location);
            const popup = document.createElement('div');
            popup.className = 'custom-popup';
            popup.innerHTML = `
                <h3 style="color: black; margin: 0; padding: 5px 0;">${location.address}</h3>
                <button class="popup-close">&times;</button>
            `;
            popup.style.position = 'absolute';
            popup.style.backgroundColor = 'white';
            popup.style.padding = '10px';
            popup.style.borderRadius = '4px';
            popup.style.boxShadow = '0 2px 4px rgba(0,0,0,0.2)';
            popup.style.zIndex = '1000';
            popup.style.minWidth = '100px';
            popup.style.textAlign = 'center';
        
            const existingPopup = document.querySelector('.custom-popup');
            if (existingPopup) {
                existingPopup.remove();
            }
        
            map.getCanvasContainer().appendChild(popup);
        
            const point = map.project([location.lon, location.lat]);
            popup.style.left = `${point.x - (popup.offsetWidth / 2)}px`;
            popup.style.top = `${point.y - popup.offsetHeight - 10}px`;

            const closeButton = popup.querySelector('.popup-close');
            closeButton.addEventListener('click', (e) => {
                e.stopPropagation();
                popup.remove();
            });
        });

        new mapboxgl.Marker({ element: el })
            .setLngLat([location.lon, location.lat])
            .addTo(map);

        bounds.extend([location.lon, location.lat]);
    });

    if (!bounds.isEmpty()) {
        map.fitBounds(bounds, { padding: 50 });
    } else {
        console.warn('No valid locations to display on the map');
    }

    map.on('move', () => {
        const popup = document.querySelector('.custom-popup');
        if (popup) {
            popup.remove();
        }
    });
}

document.getElementById('creation-year').addEventListener('change', () => searchArtists(document.getElementById('search-input').value));
document.getElementById('first-album-year').addEventListener('change', () => searchArtists(document.getElementById('search-input').value));
document.querySelectorAll('#member-checkboxes input, #location-checkboxes input').forEach(checkbox => {
    checkbox.addEventListener('change', () => searchArtists(document.getElementById('search-input').value));
});

function showError(message) {
    const errorElement = document.getElementById('error-message');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
    setTimeout(() => {
        errorElement.style.display = 'none';
    }, 5000);
}

// Initial search on page load
window.addEventListener('load', () => {
    searchArtists('');
});