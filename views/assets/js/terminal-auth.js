// Login Check
let printed = {}
let checked = {}
document.addEventListener("keydown", async (e) => {
    if (e.code === "Enter") {
        e.preventDefault();

        // Verify username
        let username = document.querySelector("#login-username").value;
        if (username.length > 0 && !checked["username"]) {
            await new Promise((resolve) => {
                fetch('/api/users?login=' + username, {
                        method: 'POST'
                    })
                    .then(r => r.json())
                    .then(async (d) => {
                        let logTxt = document.querySelector("#logs-text");
                        let color = "success"
                        let res = "Username OK."
                        let found = d.length === 0
                        if (found) {
                            color = "danger"
                            res = "Bad username. Try again."
                        }
                        logTxt.innerHTML +=
                            `<span class="txt ${color}"></span><br>`;

                        let txts = logTxt.querySelectorAll(".txt");
                        await typing(txts[txts.length - 1], res, "logs");
                        if (!found) {
                            resolve();
                            printed["username"] = true
                            checked["username"] = true
                        }
                    })
            })
        }

        // Verify password
        document.querySelector("#login-username").readOnly = true

        if (!printed["password"]) {
            let txtsList = [
                "Enter your ",
                "password",
                ": ",
                " "
            ];
            await typeInTxts("#login-password-row", txtsList, "lgn");
            printed["password"] = true
        }

        // Verify username
        let pwd = document.querySelector("#login-password");
        pwd.focus();
        let password = pwd.value;
        if (password.length > 0 && !checked["password"]) {
            await new Promise((resolve) => {
                fetch(`/api/users?login=${username}&password=${password}`, {
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
                    })
            })
        }

        if (checked["username"] && checked["password"]) {
            document.querySelector("#welcome").classList.add("slide-top-left");
            document.querySelector("#logs").classList.add("slide-bottom-left");
            document.querySelector("#auth-right").classList.add("slide-right");
            await new Promise((resolve) => {
                setTimeout(()=>{
                    document.querySelector("#container").innerHTML = "";
                    resolve();
                }, 1800)
            })
        }
    }
});