// Initialization
let printed = {};
let checked = {};
let authRoute = '/api/auth';

let capitalize = (str) => str.charAt(0).toUpperCase() + str.substr(1);

async function ValidateField(route, selector, varName, data, checked, printed) {
    let field = document.querySelector(selector);
    field.focus();
    let val = field.value;
    data.append(varName, val)
    if (val.length > 0 && !checked[varName]) {
        await new Promise((resolve) => {
            fetch(route, {
                    method: 'POST',
                    body: data
                })
                .then(r => r.json())
                .then(async (d) => {
                    let logTxt = document.querySelector("#logs-text");
                    let res = `${capitalize(varName)} OK.`;
                    let color = "success";
                    let found = d.value === "SUCCESS";

                    if (!found) {
                        color = "danger";
                        res = d.value;
                    } else {
                        printed[varName] = true;
                        checked[varName] = true;
                        fieldLoaded = true;
                    }

                    logTxt.innerHTML +=
                        `<span class="txt ${color}"></span><br>`;
                    let txts = logTxt.querySelectorAll(".txt");
                    await typing(txts[txts.length - 1], res, "logs");

                    resolve();
                });
        });
    }
    return [data, checked, printed]
}


async function ValidateNextField(route, prev, next, prevFiedSelector, nextFieldSelector, rowSelector, data, checked, printed) {
    if (checked[prev]) {
        if (!printed[next]) {
            document.querySelector(prevFiedSelector).readOnly = true;
            let txtsList = [
                "Enter your ",
                next,
                ": ",
                " "
            ];
            let typingClass = activeForm === "login" ? "lgn" : "rgst"
            await typeInTxts(rowSelector, txtsList, typingClass);
            printed[next] = true;
        }

        return await ValidateField(route, nextFieldSelector, next, data, checked, printed);
    }
}

async function hideAuth() {
    document.querySelector("#welcome").classList.add("slide-top-left");
    document.querySelector("#logs").classList.add("slide-bottom-left");
    document.querySelector("#auth-right").classList.add("slide-right");
    await new Promise((resolve) => {
        setTimeout(() => {
            document.querySelector("#container").innerHTML = "";
            resolve();
        }, 1800)
    });
}


document.addEventListener("keydown", async (e) => {
    if (e.code === "Enter") {
        e.preventDefault();

        // Verify nickname
        let data = new FormData();
        data.append("form", activeForm);
        if (activeForm === "login") {
            [data, checked, printed] = await ValidateField(authRoute, "#login-nickname", "login", data, checked, printed);

            if (fieldLoaded) {
                // Verify password
                [data, checked, printed] = await ValidateNextField(authRoute, "login", "password", "#login-nickname", "#login-password", "#login-password-row", data, checked, printed);

                if (checked["login"] && checked["password"]) {
                    hideAuth();
                }
            }
        } else if (activeForm === "register") {
            [data, checked, printed] = await ValidateField(authRoute, "#register-firstname", "firstname", data, checked, printed);

            if (fieldLoaded) {
                // Verify lastname
                [data, checked, printed] = await ValidateNextField(authRoute, "firstname", "lastname", "#register-firstname", "#register-lastname", "#register-lastname-row", data, checked, printed);

                if (checked["firstname"] && checked["lastname"]) {
                    hideAuth();
                }
            }
        }
    }
});