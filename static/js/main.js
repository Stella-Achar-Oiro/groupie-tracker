document.addEventListener('DOMContentLoaded', () => {
    const artistCards = document.querySelectorAll('.artist-card');
    artistCards.forEach(card => {
        card.addEventListener('click', () => {
            const artistId = card.dataset.id;
            window.location.href = `/artist/${artistId}`;
        });
    });

    const searchInput = document.getElementById('search-input');
    const searchButton = document.getElementById('search-button');

    searchButton.addEventListener('click', () => {
        const query = searchInput.value;
        fetch(`/api/search?q=${encodeURIComponent(query)}`)
            .then(response => response.json())
            .then(data => {
                updateArtistContainer(data);
            })
            .catch(error => console.error('Error:', error));
    });
});

function updateArtistContainer(artists) {
    const container = document.getElementById('artist-container');
    container.innerHTML = '';
    artists.forEach(artist => {
        const card = document.createElement('div');
        card.className = 'artist-card';
        card.dataset.id = artist.id;
        card.innerHTML = `
            <img src="${artist.image}" alt="${artist.name}">
            <h2>${artist.name}</h2>
            <p>Created: ${artist.creationDate}</p>
        `;
        container.appendChild(card);
    });
}