function enableReactions() {
    document.querySelectorAll(".like-dislike").forEach(el => {
        el.addEventListener("click", function (e) {
            let target = e.target.closest(".reaction");
            if (target.classList.contains("like-btn")) {
                like(target);
            } else if (target.classList.contains("dislike-btn")) {
                dislike(target);
            }
        });
    });
}

function like(el) {
    sendData(el, "/like")
}

function dislike(el) {
    sendData(el, "/dislike")
}

function sendData(el, route) {
    var data = el.parentElement.getAttribute("data-post");
    var typeReacted = "post";
    if (el.parentElement.classList.contains("comments-reactions")) {
        typeReacted = "comment";
    }
    var formData = new FormData();
    formData.append("post", data);
    formData.append("action", route);
    formData.append("typeReacted", typeReacted);
    fetch(route, {
            method: 'POST',
            body: formData
        })
        .then(response => response.json())
        .then(data => {
            if (data.value == "SUCCESS") {
                countUpdate(el);
            }
        })
        .catch(error => console.log("Error: ", error));
}

function countUpdate(el) {
    var parent = el.parentElement;
    var like = parent.querySelectorAll(".like-btn")[0];
    var dislike = parent.querySelectorAll(".dislike-btn")[0];
    var count = el.querySelectorAll(".reaction-count")[0];
    var intCount = parseInt(count.textContent);
    var alterEgo = function (element) {
        if (element == like) {
            return dislike
        } else {
            return like
        }
    };
    var count = el.querySelector(".reaction-count");
    var intCount = parseInt(count.textContent);
    const y = alterEgo(el);
    var yCount = y.querySelector(".reaction-count");
    var yIntCount = parseInt(yCount.textContent);
    const noReaction = !like.classList.contains("reaction-active") && !dislike.classList.contains("reaction-active");
    const oneReaction = (like.classList.contains("reaction-active") && !dislike.classList.contains("reaction-active")) ||
        (!like.classList.contains("reaction-active") && dislike.classList.contains("reaction-active"));

    if (noReaction) {
        count.textContent = String(intCount + 1);
        el.classList.toggle("reaction-active");
        el.querySelector(".reaction-check").textContent = "X";
    } else if (oneReaction) {
        if (el.classList.contains("reaction-active")) { // Un seul bouton est actif
            // Lorsque c'est le bouton cliqué qui est actif,
            // on le désactive et on diminue son compteur
            if (intCount <= 1) {
                count.textContent = String(0);
            } else {
                count.textContent = String(intCount - 1);
            }
            el.classList.toggle("reaction-active");
            el.querySelector(".reaction-check").textContent = " ";
        } else {
            // Lorsque c'est le bouton cliqué qui est inactif
            // on désactive son alter ego et on diminue son compteur
            if (yIntCount <= 1) {
                yCount.textContent = String(0);
            } else {
                yCount.textContent = String(yIntCount - 1);
            }
            y.classList.toggle("reaction-active");
            y.querySelector(".reaction-check").textContent = " ";
            // Puis on active le bouton cliqué et on
            // incrément son compteur
            count.textContent = String(intCount + 1);
            el.classList.toggle("reaction-active");
            el.querySelector(".reaction-check").textContent = "X";
        }
    }
}