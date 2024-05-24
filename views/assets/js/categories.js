// Créer un objet URL à partir de l'URL actuelle
var url = new URL(window.location.href);

// Obtenir les paramètres de l'URL
var params = url.searchParams;

// Récupérer la valeur du paramètre 'param'
var categoryActive = params.get('category');

var categoriesCard = document.querySelector("#categories-card");

if (categoryActive == null) {
    categoriesCard.querySelector("#category-all").classList.add("category-active");
} else {
    categoriesCard.querySelectorAll(".category-bubble").forEach(el => {
        if (el.innerHTML == categoryActive) {
            el.classList.add("category-active");
        }
    });
}


document.querySelectorAll(".post-categories").forEach(el => {
    var cat = el.innerHTML.split(",");
    var catHTML = "";
    for (let index = 0; index < cat.length; index++) {
        const name = cat[index];
        catHTML += `<span class="category-badge">` + name + `</span>`;
    }
    el.innerHTML = catHTML;
});