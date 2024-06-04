const speed = 15;
var activeForm = "login";
let loginTxtsData = [{
        zone: "#login-intro",
        className: "lgn",
        txtsList: [
            "$> help",
            "These infos are needed for login :",
            " - nickname or email",
            " - password",
            "User should be registered first before trying to login.",
            "$> login"
        ]
    },
    {
        zone: "#login-username-row",
        className: "lgn",
        txtsList: [
            "Enter your ",
            "nickname or email",
            ": ",
            " "
        ]
    }
];

let registerTxtsData = [{
        zone: "#register-intro",
        className: "rgst",
        txtsList: [
            "$> help",
            "These infos are needed for register :",
            " - nickname or email",
            " - password",
            "User should be registered first before trying to register.",
            "$> register"
        ]
    },
    {
        zone: "#register-username-row",
        className: "rgst",
        txtsList: [
            "Enter your ",
            "nickname or email",
            ": ",
            " "
        ]
    }
];

var txtDataToLoad = () => activeForm === "login" ? loginTxtsData : registerTxtsData
var otherBtn = () => activeForm === "login" ? document.querySelector("#register-btn") : document.querySelector("#login-btn")

window.onload = async () => {
    await loadLoginDOM();
    btnBgPosition();
    await hidePreloader();
    await disableAllInputs();
    await startWelcomeTerminal();
    await startLogsTerminal();
    await shrinkLeftBox();
    await expandAuthRight();
    await startActiveTerminal(activeForm, txtDataToLoad());
    await enableAllInputs();
    resizeInput();
    firstFieldFocus("#login-form");
    formToggle();
};

function disableAllInputs() {
    return new Promise((resolve) => {
        document.querySelectorAll("input").forEach(e => {
            e.disabled = true;
        })
        resolve();
    });
}

function enableAllInputs() {
    return new Promise((resolve) => {
        document.querySelectorAll(".auth-input").forEach(e => {
            e.disabled = false;
        })
        resolve();
    });
}

function typing(element, txt, className, s = speed) {
    return new Promise(async (resolve) => {
        for (let i = 0; i < txt.length; i++) {
            await typeChar(s, element, txt, className, i);
        }
        resolve();
    });
}

function typeChar(delay, element, txt, className, i) {
    return new Promise((resolve) => {
        setTimeout(() => {
            let prevUnderscore = document.querySelector("#" + className + "-underscore");
            if (prevUnderscore) {
                prevUnderscore.remove();
            }
            element.innerHTML += txt.charAt(i);
            let newUnderscore = document.createElement("span");
            newUnderscore.textContent = "_";
            newUnderscore.id = className + "-underscore";
            newUnderscore.className = "cursor blinking";
            element.appendChild(newUnderscore);
            resolve();
        }, delay);
    });
}

function hidePreloader() {
    return new Promise((resolve) => {
        setTimeout(() => {
            document.querySelector("#preloader").style.display = "none";
            resolve();
        }, 2000);
    });
}

function loadLoginDOM() {
    return new Promise((resolve) => {
        document.querySelector("#container").innerHTML = loginView;
        resolve();
    });
}

function startWelcomeTerminal() {
    return new Promise(async (resolve) => {
        // Terminal Pop-up
        let w = document.querySelector("#welcome");
        w.style.opacity = "1";
        w.style.height = "100%";
        w.style.width = "100%";

        // Typing
        await new Promise((resolve) => {
            setTimeout(async () => {
                let txtsList = ["Hello ", "talent", " !", "Welcome to our", "Real Time Forum", "."];
                await typeInTxts("#welcome-text", txtsList, "wlc", 50);
                resolve();
            }, 500)
        });
        resolve();
    });
}

function startLogsTerminal() {
    return new Promise(async (resolve) => {
        // Shrink welcome terminal
        let w = document.querySelector("#welcome");
        w.style.height = "48%";

        document.querySelector("#welcome-text").style.fontSize = "2em";

        // Pop-up Logs terminal
        await new Promise((resolve) => {
            setTimeout(() => {
                let l = document.querySelector("#logs");
                l.style.opacity = "1";
                l.style.height = "48%";
                l.style.width = "100%";
                resolve()
            }, 200);
        });

        // Typing
        await new Promise((resolve) => {
            setTimeout(async () => {
                let txtsList = ["$> Page loaded successfully"];
                await typeInTxts("#logs-text", txtsList, "logs");
                resolve();
            }, 500)
        });
        resolve();
    });
}

function shrinkLeftBox() {
    return new Promise((resolve) => {
        let leftBox = document.querySelector("#auth-left")
        leftBox.style.width = "calc(50% - 2vh)";
        leftBox.style.marginLeft = "0";
        resolve();
    });
}

function expandAuthRight() {
    return new Promise((resolve) => {
        setTimeout(() => {
            document.querySelector("#auth-right").style.width = "calc(50% - 2vh)";
            resolve();
        }, 400);
    });
}

function startActiveTerminal(terminal, txtData) {
    return new Promise(async (resolve) => {
        // Terminal Pop-up
        let w = document.querySelector("#" + terminal);
        w.style.overflow = "hidden";

        await new Promise((resolve) => {
            setTimeout(async () => {
                w.style.opacity = "1";
                w.style.height = "100%";
                w.style.width = "100%";
                resolve();
            }, 200)
        });


        // Typing
        await new Promise((resolve) => {
            setTimeout(async () => {
                for (const d of txtData) {
                    await typeInTxts(d.zone, d.txtsList, d.className);
                }
                resolve();
            }, 500)
        });
        resolve();
    });
}

async function typeInTxts(selector, txtsList, className, s = speed) {
    let elementsList = document.querySelectorAll(`${selector} .txt`);

    for (let x = 0; x < elementsList.length; x++) {
        await typing(elementsList[x], txtsList[x], className, s);
    }
}

function resizeInput() {
    document.querySelectorAll(".auth-input").forEach(el => {
        el.parentNode.addEventListener("click", (e) => {
            el.focus();
        });

        el.addEventListener("input", (e) => {
            let m = el.parentNode.querySelector(".measurer");
            m.textContent = el.value;
            let mWidth = m.offsetWidth;
            el.style.width = (mWidth + 15) + "px";
        });
    });
}

function firstFieldFocus(selector) {
    document.querySelectorAll(`${selector} input`)[0].focus();
}

function formToggle() {
    document.querySelectorAll("#register-btn,#login-btn").forEach(el => {
        el.addEventListener("click", async (e) => {
            if (!e.target.classList.contains("active")) {
                let authBox = document.querySelector("#" + activeForm);
                await new Promise((resolve) => {
                    authBox.classList.add("instant-slide-top");
                    resolve();
                })
                await new Promise((resolve) => {
                    setTimeout(async () => {
                        authBox.remove();
                        document.querySelector("#terminal-zone").innerHTML = activeForm === "login" ? registerTerminal : loginTerminal;
                        activeForm = activeForm === "login" ? "register" : "login";
                        document.querySelector("#" + activeForm).classList.remove("not-displayed");
                        e.target.classList.add("active")
                        otherBtn().classList.remove("active")
                        await startActiveTerminal(activeForm, txtDataToLoad());
                        resolve();
                    }, 750);
                })
            }
        });
    });
}

function btnBgPosition() {
    document.querySelectorAll(".btn-bg").forEach(el=>{
        if (el.parentNode.id === "login-btn") {
            el.style.right = "0"
        } else {
            el.style.left = "0"
        }
    })
}