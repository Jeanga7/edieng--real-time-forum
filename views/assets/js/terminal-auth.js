// Login Check
let printed = {}
let checked = {}
document.addEventListener("keydown", async (e) => {
    if (e.code === "Enter") {
        e.preventDefault();

        // Verify nickname
        let data = new FormData()
        let nickname = document.querySelector("#login-nickname").value;
        data.append("login", nickname)
        
        if (nickname.length > 0 && !checked["nickname"]) {
            await new Promise((resolve) => {
                fetch('/api/users', {
                        method: 'POST',
                        body: data
                    })
                    .then(r => r.json())
                    .then(async (d) => {
                        let logTxt = document.querySelector("#logs-text");
                        let color = "success";
                        let res = "Nickname OK.";
                        let found = d.length === 0;
                        if (found) {
                            color = "danger";
                            res = "Bad nickname. Try again.";
                        }
                        logTxt.innerHTML +=
                            `<span class="txt ${color}"></span><br>`;

                        let txts = logTxt.querySelectorAll(".txt");
                        await typing(txts[txts.length - 1], res, "logs");
                        if (!found) {
                            printed["nickname"] = true;
                            checked["nickname"] = true;
                            resolve();
                        }
                    });
            });
        }

        // Verify password
        if (checked["nickname"]) {
            if (!printed["password"]) {
                document.querySelector("#login-nickname").readOnly = true;
                let txtsList = [
                    "Enter your ",
                    "password",
                    ": ",
                    " "
                ];
                await typeInTxts("#login-password-row", txtsList, "lgn");
                printed["password"] = true;
            }

            let pwd = document.querySelector("#login-password");
            pwd.focus();
            let password = pwd.value;
            if (password.length > 0 && !checked["password"]) {
                await new Promise((resolve) => {
                    fetch(`/api/users?login=${nickname}&password=${password}`, {
                            method: 'POST'
                        })
                        .then(r => r.json())
                        .then(async (d) => {
                            let logTxt = document.querySelector("#logs-text");
                            let color = "success"
                            let res = "Password OK."
                            let found = d.length === 0
                            if (found) {
                                color = "danger"
                                res = "Bad password. Try again."
                            }
                            logTxt.innerHTML +=
                                `<span class="txt ${color}"></span><br>`;

                            let txts = logTxt.querySelectorAll(".txt");
                            await typing(txts[txts.length - 1], res, "logs");
                            if (!found) {
                                resolve();
                                printed["password"] = true
                                checked["password"] = true
                            }
                        });
                });
            }
        }

        if (checked["nickname"] && checked["password"]) {
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
    }
});