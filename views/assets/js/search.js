document.getElementById('search-field').addEventListener('keydown', function (event) {
    // Vérifier si la touche appuyée est Entrée
    if (event.key === 'Enter') {
        var searchValue = event.target.value;

        var redirectUrl = '/?search=' + encodeURIComponent(searchValue);

        window.location.href = redirectUrl;
    }
});