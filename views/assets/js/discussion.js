document.querySelector(".discussion").addEventListener("click", function () {
    var name = document.querySelector(".discussion").textContent.split("-")[0];
    console.log(name);

    // var popup = document.createElement("div");
    // popup.classList.add("layer");
    // document.body.appendChild(popup);
    // popup.innerHTML = `
    //     <div class="popup">
    //         <span id="close_popup">X</span>
    //     </div>
    // `;
    
    let layer = document.querySelector(".layer");
    layer.style.display = "flex";
});

let closePopup = document.querySelector("#close_popup");
if (closePopup) {
    closePopup.addEventListener("click", () => {
        var layer = document.querySelector(".layer");
        layer.style.display = "none";
    });
}