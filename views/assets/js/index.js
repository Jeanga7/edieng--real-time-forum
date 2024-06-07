var modal = document.querySelector("#new-post-modal");

document.querySelectorAll("#new-post-btn, #new-post-close-btn").forEach(el => {
    el.addEventListener("click", function(){
        toggleModal();
    });
});

document.addEventListener("click", function(event){
    if (event.target.matches("#new-post-close-btn") || event.target.matches("#new-post-modal")) {
        toggleModal();
    }
});

document.querySelectorAll(".ellipsis").forEach((el) => {
    el.addEventListener("click", function(){
        if (el.classList.contains("ellipsis")) {
            el.classList.remove("ellipsis");
        } else {
            el.classList.add("ellipsis");
        }
    });
});

document.querySelector("#new-post-file-upload").addEventListener("click", function(){
    document.querySelector("#new-post-image").click();
});

document.querySelector("#new-post-image").addEventListener("input", function(e){
    // Vérifiez si un fichier a été sélectionné
    if (e.target.files.length > 0) {
        // Récupérez le nom du premier fichier sélectionné
        var nomFichier = e.target.files[0].name;
        
        // Affichez le nom du fichier dans l'élément <p>
        document.querySelector("#upload-txt").innerHTML = nomFichier;
    }
});

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

modal.querySelectorAll(".category-bubble").forEach((el) => {
    el.addEventListener("click", function(){
        var chbx = el.querySelector(".post-category-checkbox");
        chbx.checked = !chbx.checked

        if (el.classList.contains("category-active")) {
            el.classList.remove("category-active");
        } else {
            el.classList.add("category-active");
        }
    });
});

document.querySelector("#new-post-submit-btn").addEventListener("click", function(event){
    event.preventDefault();

    // Sélectionnez le formulaire, pas le bouton
    const form = document.querySelector("#new-post-modal-body");
    const formData = new FormData(form);
    const route = form.getAttribute("action");

    // Envoi asynchrone des données du formulaire
    fetch(route, {
        method: 'POST',
        body: formData
    })
    .then(response => {
        // Vérifie si la réponse est au format JSON
        const contentType = response.headers.get("content-type");
        if (contentType && contentType.includes("application/json")) {
            return response.json();
        } else {
            throw new Error("La réponse n'est pas au format JSON");
        }
    })
    .then((data) => {
        postErrorHandler(data.value);
    })
    .catch((error) => {
        console.error('Error:', error);
        postErrorHandler(data.value);
    });
});


document.querySelectorAll(".comment-btn").forEach(el => {
    var comments = el.parentElement.parentElement.querySelector(".comments");
    
    el.addEventListener("click", function() {
        if (comments.style.height == "0" || comments.style.height == "" || comments.style.height == "0px") {
            comments.style.height = "auto";
        } else {
            comments.style.height = "0";
        }
    });
});

document.querySelectorAll(".add-comment-form").forEach(el => {
    el.addEventListener("submit", function(event){
        event.preventDefault();

        var formData = new FormData(event.target);
        var route = el.getAttribute("action");

        fetch(route, {
            method: "POST",
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.value == "SUCCESS") {
                toggleModal();
                successShow("comment");
            } else {
                toggleModal();
                failShow(data.value);
            }
        })
        .catch((error) => {
            console.log('Error: ', error);
        });
    });
});

document.querySelector("#new-category").addEventListener("click", function(){
    var parent = this.parentNode;
    
    // Créez un nouvel élément à partir de la chaîne de caractères HTML
    var newElement = document.createElement('input');
    newElement.type = "text";
    newElement.name = "categories[]";
    newElement.className = "comment-txt";
    newElement.id = "new-cat-input";
    
    // Remplacez l'élément existant par le nouvel élément
    parent.replaceChild(newElement, this);
});







function toggleModal() {
    // Initialiser la propriété visibility si elle n'est pas définie
    if (modal.style.visibility === "") {
        modal.style.visibility = "hidden";
    }

    if (modal.style.visibility == "hidden") {
        modal.style.visibility = "visible";
        modal.style.opacity = "1";
        modal.style.backdropFilter = "blur(4px)";
    } else {
        modal.style.opacity = "0";
        modal.style.visibility = "hidden";
        modal.style.backdropFilter = "blur(0)";
    }
}

var hidden = document.querySelector("#hidden");
var successAnimation = hidden.querySelector("#success-animation");

var modalBody = document.querySelector("#new-post-modal-body");

function successShow(val) {
    modalBody.style.display = "flex";
    modalBody.style.justifyContent = "center";
    modalBody.innerHTML = `<div class="add-animation" id="success-animation">
    <lottie-player src="https://lottie.host/495eacba-f26b-4e22-83d6-acf4f6294784/KAIPSUEkAg.json" background="transparent" speed="1" style="width: 300px; height: 300px" loop autoplay direction="1" mode="normal"></lottie-player>
    <h3 style="color:green" class="res-txt">New `+val+` added !</h3>
    </div>`;
    successAnimation.style.margin = "auto";
    setTimeout(function(){
        window.location.href = "/";
    }, 2000);
}

function failShow(val) {
    modalBody.style.display = "flex";
    modalBody.style.justifyContent = "center";
    modalBody.innerHTML = `<div class="add-animation" id="failure-animation">
    <lottie-player src="https://lottie.host/13886175-6124-444a-8716-946ada521fa6/XcTHWAJd48.json" background="transparent" speed="1" style="width: 300px; height: 300px" loop autoplay direction="1" mode="normal"></lottie-player>
    <h3 style="color:red" class="res-txt">`+val+`</h3>
    </div>`;
    successAnimation.style.margin = "auto";
    setTimeout(function(){
        window.location.href = "/";
    }, 2000);
}

function postErrorHandler(value) {
    switch (value) {
        case "SUCCESS":
            // Afficher le modal de résultat avec comme paramètre échec
            successShow("post");
            return;
        default:
            failShow(value);
    }
  }


  function createCommentHTML(authorName, authorNickname, commentContent) {
    return `
        <div class="comment">
            <div class="comment-author">
                <span style="background-image: url(assets/img/profile/spidey.png)" class="comment-author-pic"></span>
                <div class="comment-author-infos">
                    <span class="comment-author-name">${authorName}</span>
                    <span class="comment-author-nickname">${authorNickname}</span>
                </div>
            </div>
            <div class="comment-content ellipsis">
                ${commentContent}
            </div>
        </div>
    `;
}
