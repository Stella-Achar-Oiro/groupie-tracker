document.addEventListener('DOMContentLoaded', function() {
    fetchArtists();
});

function fetchArtists() {
    fetch('/artists')
        .then(response => response.json())
        .then(data => {
            const artistsContainer = document.getElementById('artists');
            data.Artists.forEach(artist => {
                const artistCard = createArtistCard(artist);
                artistsContainer.appendChild(artistCard);
            });
        })
        .catch(error => console.error('Error:', error));
}

function createArtistCard(artist) {
    const card = document.createElement('div');
    card.className = 'artist-card';
    card.innerHTML = `
        <img src="${artist.Image}" alt="${artist.Name}">
        <h2>${artist.Name}</h2>
        <p>Created: ${artist.CreationDate}</p>
        <p>First Album: ${artist.FirstAlbum}</p>
    `;
    card.onclick = () => showArtistDetails(artist);
    return card;
}

function showArtistDetails(artist) {
    fetch(`/events?artist=${artist.ID}`)
        .then(response => response.json())
        .then(events => {
            const detailsContainer = document.getElementById('artist-details');
            detailsContainer.innerHTML = `
                <h2>${artist.Name}</h2>
                <img src="${artist.Image}" alt="${artist.Name}">
                <p>Members: ${artist.Members.join(', ')}</p>
                <p>Creation Date: ${artist.CreationDate}</p>
                <p>First Album: ${artist.FirstAlbum}</p>
                <h3>Upcoming Events:</h3>
                <ul>
                    ${Object.entries(events.DatesLocations).map(([location, dates]) => `
                        <li>${location}: ${dates.join(', ')}</li>
                    `).join('')}
                </ul>
            `;
            document.getElementById('details').style.display = 'block';
        })
        .catch(error => console.error('Error:', error));
}

document.querySelector('.close').onclick = function() {
    document.getElementById('details').style.display = 'none';
}

window.onclick = function(event) {
    const modal = document.getElementById('details');
    if (event.target == modal) {
        modal.style.display = 'none';
    }
}