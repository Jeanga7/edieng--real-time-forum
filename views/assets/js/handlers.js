function handleWebsocket() {
    let ws = new WebSocket(`ws://localhost${port}/messages`);
    ws.onopen = () => {
        console.log("Connected to WebSocket...");
        let registerMsg = {
            "from": userNickname,
            "to": "@system",
            "type": "join",
            "message": "",
            "signature": sessionID
        }

        setTimeout(() => {
            ws.send(JSON.stringify(registerMsg));
        }, 500)
    };
    ws.onmessage = (e) => {
        let response = JSON.parse(e.data);

        if (Object.prototype.toString.call(response) === "[object String]") {
            response = JSON.parse(response);
        }

        switch (response.type) {
            case "success":
                console.log(response.message);
                break;

            case "message":
                if (chatWith === response.from) {
                    // Insert new message
                    let chatInputName = document.querySelector(".netcat-input-name");
                    insertMsg(chatInputName, false, "warning", response.time, response.from, response.message)
                    terminalScrollBottom("#netcat-input");
                } else {
                    // Notify
                    let user = document.querySelector(`.user[data-nickname="${response.from}"]`);
                    moveElementToParentStart(user);
                    let notifsContent = user.querySelector(".user-notifications").textContent;
                    let notifs = parseInt(notifsContent.replace("(", "").replace(")", ""));
                    user.querySelector(".user-notifications").textContent = `(${notifs+1})`;
                }
                break;

            case "online":
                let userToOnline = document.querySelector(`[data-nickname="${response.message}"]`);
                if (userToOnline) {
                    userToOnline.querySelector(".user-status").classList.add("status-online");
                    userToOnline.querySelector(".user-nickname").classList.add("txt-online");
                    userToOnline.querySelector(".user-notifications").classList.add("txt-online");
                }
                break;

            case "offline":
                let userToOffline = document.querySelector(`[data-nickname="${response.message}"]`);
                if (userToOffline) {
                    userToOffline.querySelector(".user-status").classList.remove("status-online");
                    userToOffline.querySelector(".user-nickname").classList.remove("txt-online");
                    userToOffline.querySelector(".user-notifications").classList.remove("txt-online");
                }
                break;
        }
    }
}

function handleLoadMoreMsg() {
    document.querySelector("#netcat").addEventListener("scroll", debounceScroll);
}

function handleOpenNewPost() {
    document.querySelector("#add-post-btn").addEventListener("click", async (e) => {
        let np = document.querySelector("#new-post");
        let npl = document.querySelector("#new-post-logs");
        if (np.classList.contains("active-new-post")) {
            document.querySelector("#new-post-text").textContent = "";
            document.querySelector("#new-post-logs-text").textContent = "";
            np.classList.remove("active-new-post");
            npl.classList.remove("active-new-post-logs");
        } else {
            np.classList.add("active-new-post");
            npl.classList.add("active-new-post-logs");
            await new Promise((resolve) => {
                setTimeout(async () => {
                    await printInTerminal("#new-post-text", newPostTxt, "np", 10);
                    await printInTerminal("#new-post-logs-text", newPostLogsTxt, "npl", 10);
                    resolve();
                }, 500);
            })

            await addInputBeforeCursor("new-post")
        }
    });
}