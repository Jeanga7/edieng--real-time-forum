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
                    let logsZone = document.querySelector("#logs");
                    logsZone.scrollTop = logsZone.scrollHeight;
                    if (d.more) {                        
                        responseData = JSON.parse(d.more);
                        userID = responseData.id;
                    }
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
            let txtsList = [{
                    txt: "Enter your ",
                    c: null,
                    eol: false
                },
                {
                    txt: next,
                    c: "bold-fluo",
                    eol: false
                },
                {
                    txt: ": ",
                    c: null,
                    eol: false
                },
                {
                    txt: " ",
                    c: null,
                    eol: false
                },
            ]
            let typingClass = activeForm === "login" ? "lgn" : "rgst"
            await typeInTxts(rowSelector, txtsList, typingClass);
            printed[next] = true;
        }

        return await ValidateField(route, nextFieldSelector, next, data, checked, printed);
    }
}

async function switchToFeeds() {
    return new Promise((resolve) => {
        document.querySelector("#welcome").classList.add("slide-top-left");
        document.querySelector("#logs").classList.add("slide-bottom-left");
        document.querySelector("#auth-right").classList.add("slide-right");
        setTimeout(async () => {
            document.querySelector("#container").innerHTML = feedsView;
            await startFeedsTerminals();
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
                    await switchToFeeds();
                }
            }
        } else if (activeForm === "register") {
            let prev = ""
            let fieldset = ["firstname", "lastname", "nickname", "email", "password", "age", "gender"];

            for (let i = 0; i < fieldset.length; i++) {
                const field = fieldset[i];
                checked[field] = false;
                if (i === 0) {
                    [data, checked, printed] = await ValidateField(authRoute, `#register-${field}`, field, data, checked, printed);
                } else if (fieldLoaded && document.querySelector(`#register-${prev}`).value.length > 0) {
                    [data, checked, printed] = await ValidateNextField(authRoute, prev, field, `#register-${prev}`, `#register-${field}`, `#register-${field}-row`, data, checked, printed);
                }
                prev = field;
            }

            let allChecked = fieldset.reduce((res, check) => res && checked[check], true);

            if (allChecked) {
                await switchToFeeds();
            }
        }
    }
});