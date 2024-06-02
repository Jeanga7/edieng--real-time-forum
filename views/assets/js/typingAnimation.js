const speed = 15;

window.onload = async () => {
    await hidePreloader();
    await disableAllInputs();
    await startWelcomeTerminal();
    await startLogsTerminal();
    await shrinkLeftBox();
    await expandAuthRight();
    await startLoginTerminal();
    await enableAllInputs();
    resizeInput();
    console.log("OK");
};

function disableAllInputs() {
    return new Promise((resolve) => {
        document.querySelectorAll("input").forEach(e=>{
            e.disabled = true;
        })
        resolve();
    })
}

function enableAllInputs() {
    return new Promise((resolve) => {
        document.querySelectorAll(".auth-input").forEach(e=>{
            e.disabled = false;
        })
        resolve();
    })
}

function typing(element, txt, className, s=speed) {
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
            let prevU = document.querySelector("#" + className + "-underscore");
            if (prevU) {
                prevU.remove();
            }
            element.innerHTML += txt.charAt(i);
            let newU = document.createElement("span");
            newU.textContent = "_";
            newU.id = className + "-underscore";
            newU.className = "cursor blinking"
            element.appendChild(newU);
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
                let elementsList = document.querySelectorAll("#welcome-text .txt");
                let txtsList = ["Hello ", "talent", " !", "Welcome to our", "Real Time Forum", "."];

                for (let x = 0; x < elementsList.length; x++) {
                    await typing(elementsList[x], txtsList[x], "wlc", 50);
                }
                resolve()
            }, 500)
        })

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
            }, 200)
        })

        // Typing
        await new Promise((resolve) => {
            setTimeout(async () => {
                let elementsList = document.querySelectorAll("#logs-text .txt");
                let txtsList = ["$> Page loaded successfully"];

                for (let x = 0; x < elementsList.length; x++) {
                    await typing(elementsList[x], txtsList[x], "logs");
                }
                resolve()
            }, 500)
        })

        resolve()
    })
}

function shrinkLeftBox() {
    return new Promise((resolve) => {
        let leftBox = document.querySelector("#auth-left")
        leftBox.style.width = "calc(50% - 2vh)";
        leftBox.style.marginLeft = "0";

        resolve()
    })
}

function expandAuthRight() {
    return new Promise((resolve) => {
        setTimeout(() => {
            document.querySelector("#auth-right").style.width = "calc(50% - 2vh)";
            resolve();
        }, 400);
    })
}

function startLoginTerminal() {
    return new Promise(async (resolve) => {
        // Terminal Pop-up
        let w = document.querySelector("#login");
        w.style.opacity = "1";
        w.style.height = "95%";
        w.style.width = "100%";

        // Typing
        await new Promise((resolve) => {
            setTimeout(async () => {
                let elementsList = document.querySelectorAll("#login-text .txt");
                let txtsList = [
                    "$> help",
                    "These infos are needed for login :",
                    " - nickname or email",
                    " - password",
                    "User should be registered first before trying to login.",
                    "$> login",
                    "Enter your ",
                    "nickname or email",
                    ": ",
                    " "
                ];

                for (let x = 0; x < elementsList.length; x++) {
                    await typing(elementsList[x], txtsList[x], "lgn");
                }
                resolve()
            }, 500)
        })

        resolve();
    });
}

function resizeInput() {
    document.querySelectorAll(".auth-input").forEach(el=>{
        el.parentNode.addEventListener("click", (e)=>{
            el.focus();
        });
        
        el.addEventListener("input", (e)=>{
            let m = el.parentNode.querySelector(".measurer");
            m.textContent = el.value;
            let mWidth = m.offsetWidth;
            el.style.width = (mWidth+15) + "px";
        });
    })
}