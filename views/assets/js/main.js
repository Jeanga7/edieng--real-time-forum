enableMode();
document.querySelector("#light").addEventListener("click", function(){
    toggleMode();
});

document.querySelector("#logo").addEventListener("click", function(){
    window.location.href = "/";
});


function setCookie(cname, cvalue, exdays) {
    var d = new Date();
    d.setTime(d.getTime() + (exdays*24*60*60*1000));
    var expires = "expires="+ d.toUTCString();
    // Ajout de l'attribut SameSite avec la valeur 'Lax' ou 'Strict' selon vos besoins
    // 'Lax' permet au cookie d'être envoyé dans les requêtes cross-site, mais seulement pour les requêtes GET
    // 'Strict' empêche le cookie d'être envoyé dans les requêtes cross-site
    var sameSite = "SameSite=Strict";
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/;" + sameSite;
}


function getCookie(cname) {
    var name = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i <ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function enableMode() {
    var mode = getCookie("lightMode");
    if (mode == "light" || mode == "dark") {
        switchLightsTo(mode);
    } else {
        setCookie('lightMode', 'dark', 7);
    }
}

function toggleMode() {
    var mode = getCookie("lightMode");
    if (mode == "light") {
        switchLightsTo("dark");
    } else if (mode == "dark") {
        switchLightsTo("light");
    }
}

function switchLightsTo(mode) {
    if (mode == "dark") {
        document.documentElement.style.setProperty('--main-bg', '#000');
        document.documentElement.style.setProperty('--black', 'white');
        document.documentElement.style.setProperty('--white', '#1d1d1d');
        document.documentElement.style.setProperty('--txt-black', 'white');
        setCookie('lightMode', 'dark', 7);
    } else if (mode == "light") {
        document.documentElement.style.setProperty('--main-bg', '#ddd');
        document.documentElement.style.setProperty('--black', 'black');
        document.documentElement.style.setProperty('--white', 'white');
        document.documentElement.style.setProperty('--txt-black', 'black');
        setCookie('lightMode', 'light', 7);
    }
}