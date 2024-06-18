document.querySelectorAll(".like-dislike").forEach(el => {
    el.addEventListener("click", function () {
        if (el.classList.contains("like")) {
            like(el);                 
        } else if (el.classList.contains("dislike")) {
            dislike(el);
        }
    });
});

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

    fetch(route,{
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
    var like = parent.querySelector(".like");
    var dislike = parent.querySelector(".dislike");
    var count = el.querySelector(".reaction-count");
    var intCount = parseInt(count.innerHTML);
    var alterEgo = function(element) {
        if (element == like) {
            return dislike
        } else {
            return like
        }
    };

    const noReaction = !like.classList.contains("reaction-active") && !dislike.classList.contains("reaction-active");
    const oneReaction = (like.classList.contains("reaction-active") && !dislike.classList.contains("reaction-active")) || (!like.classList.contains("reaction-active") && dislike.classList.contains("reaction-active"));

    if (noReaction) {
        count.innerHTML = String(intCount + 1);
        el.classList.add("reaction-active");
    } else if (oneReaction) {
        var count = el.querySelector(".reaction-count");
        var intCount = parseInt(count.innerHTML);
        if (el.classList.contains("reaction-active")) { // Un seul bouton est actif
            // Lorsque c'est le bouton cliqué qui est actif,
            // on le désactive et on diminue son compteur
            if (intCount <= 1) {
                count.innerHTML = String(0);
            } else {
                count.innerHTML = String(intCount - 1);
            }
            el.classList.remove("reaction-active");
        } else {
            // Lorsque c'est le bouton cliqué qui est inactif
            // on désactive son alter ego et on diminue son compteur
            const y = alterEgo(el);
            var yCount = y.querySelector(".reaction-count");
            var yIntCount = parseInt(yCount.innerHTML);
            if (yIntCount <= 1) {
                yCount.innerHTML = String(0);
            } else {
                yCount.innerHTML = String(yIntCount - 1);
            }
            y.classList.remove("reaction-active");

            // Puis on active le bouton cliqué et on
            // incrément son compteur
            count.innerHTML = String(intCount + 1);
            el.classList.add("reaction-active");
        }
    }
    
    var count = ele.querySelector(".reaction-count");
    var intCount = parseInt(count.innerHTML);
    if (ele.classList.contains("reaction-active")) {
        if (intCount <= 1) {
            count.innerHTML = String(0);
        } else {
            count.innerHTML = String(intCount - 1);
        }
        ele.classList.remove("reaction-active");
    } else {
        count.innerHTML = String(intCount + 1);
        ele.classList.add("reaction-active");
    }
}