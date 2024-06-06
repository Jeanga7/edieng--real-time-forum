// Initialization
let printed = {};
let checked = {};
let loaded = false;

let capitalize = (str) => str.charAt(0).toUpperCase() + str.substr(1);

async function ValidateField(selector, varName, data, checked, printed) {
    let field = document.querySelector(selector);
    field.focus();
    let val = field.value;
    data.append(varName, val)
    if (val.length > 0 && !checked[varName]) {
        loaded = true;
        await new Promise((resolve) => {
            fetch('/api/users', {
                    method: 'POST',
                    body: data
                })
                .then(r => r.json())
                .then(async (d) => {
                    let logTxt = document.querySelector("#logs-text");
                    let color = "success";
                    let res = `${capitalize(varName)} OK.`;
                    let found = d.length === 0;
                    if (found) {
                        color = "danger";
                        res = `Bad ${varName}. Try again.`;
                    }
                    logTxt.innerHTML +=
                        `<span class="txt ${color}"></span><br>`;

                    let txts = logTxt.querySelectorAll(".txt");
                    await typing(txts[txts.length - 1], res, "logs");
                    if (!found) {
                        printed[varName] = true;
                        checked[varName] = true;
                        resolve();
                    }
                });
        });
    }
    return [data, checked, printed]
}


async function ValidateNextField(prev, next, prevFiedSelector, nextFieldSelector, rowSelector, typingClass, data, checked, printed) {
    if (checked[prev]) {
        if (!printed[next]) {
            document.querySelector(prevFiedSelector).readOnly = true;
            let txtsList = [
                "Enter your ",
                next,
                ": ",
                " "
            ];
            await typeInTxts(rowSelector, txtsList, typingClass);
            printed[next] = true;
        }

        return await ValidateField(nextFieldSelector, next, data, checked, printed);
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
        [data, checked, printed] = await ValidateField("#login-nickname", "login", data, checked, printed);

        if (loaded) {
            // Verify password
            [data, checked, printed] = await ValidateNextField("login", "password", "#login-nickname", "#login-password", "#login-password-row", "lgn", data, checked, printed);

            if (checked["login"] && checked["password"]) {
                hideAuth();
            }
        }
    }
});