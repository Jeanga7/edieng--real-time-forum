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

                    let data = [{
                        txt: res,
                        c: color,
                        eol: true
                    }, ];

                    await printInTerminal("#logs-text", data, "logs", 10);

                    let logsZone = document.querySelector("#logs");
                    logsZone.scrollTop = logsZone.scrollHeight;
                    if (d.more) {
                        responseData = JSON.parse(d.more);
                        userID = responseData.id;
                        sessionID = responseData.idSession;
                        console.log(responseData);
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