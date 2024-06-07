window.onload = (e) => {
  var loginForm = document.querySelector("#login-form");
  var registerForm = document.querySelector("#register-form");

  // Switch between login & register form
  document.querySelectorAll(".change-auth").forEach((el) => {
    el.addEventListener("click", function () {

      if (loginForm.classList.contains("active")) {
        loginForm.classList.remove("active");
        registerForm.classList.add("active");
      } else {
        registerForm.classList.remove("active");
        loginForm.classList.add("active");
      }
    });
  });

  // Send form data
  document.querySelectorAll("#login-form, #register-form").forEach((el) => {
    el.addEventListener("submit", (event) => {
      event.preventDefault();

      // Récupération des données du formulaire
      var formData = new FormData(event.target);
      var route = el.getAttribute("action");

      // Envoi asynchrone des données du formulaire
      fetch(route, {
        method: 'POST',
        body: formData
      })
      .then(response => response.json())
      .then((data) => {
        fieldErrorHandler(data.value, el);
      })
      .catch((value) => {
        console.log('Error:', value);
      });
    });
  });
};

function fieldErrorHandler(value, el) {
  var errorZone = el.querySelectorAll(".error-zone");
  if (errorZone.length > 0) {
    errorZone.forEach(x => {
      x.style.display = "block";
    });
  }
  switch (value) {
    case "INVALID EMAIL":
      errorField(el, "email");
      break;
    case "INVALID NICKNAME":
      errorField(el, "nickname");
      break;
    case "NICKNAME OR EMAIL ALREADY USED":
        errorField(el, "nickname");
        break;
    case "USER NOT FOUND":
        errorField(el, "nickname");
        break;
    case "INVALID PASSWORD":
      errorField(el, "password");
      break;
    case "INVALID NAME":
      errorField(el, "name");
      break;
    case "SUCCESS":
        el.querySelector(".error-zone").remove();
        window.location.href = "/";
        return;
  }
  addErrorTimes();
  var errorHTML = `<span class="error-notice">
  <span>&times;</span>
  <span class="error-text">`+value+`</span>
</span>`;
  el.querySelector(".error-zone").innerHTML += errorHTML;
}

function addErrorTimes() {
  document.querySelectorAll(".field-error").forEach(el => {
    el.insertAdjacentHTML('afterend', '<span class="error-times">&times;</span>');
  });
}

function errorField(el, className) {
  el.querySelectorAll("input").forEach(e => {
    e.classList.remove("field-error");
  });
  el.querySelectorAll(".error-times").forEach(e => {
    e.remove();
  });
  el.querySelectorAll("."+className).forEach(e => {
    e.classList.add("field-error");
  });

  el.querySelector(".error-zone").innerHTML = "";
}