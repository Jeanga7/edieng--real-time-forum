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

function handleOpenComments() {
    document.querySelectorAll(".comment-btn").forEach(el => {
        el.addEventListener("click", async (e) => {
            let postId = el.parentElement.dataset.post;
            let ac = document.querySelector("#all-comments");
            document.querySelector("#all-comments-text").textContent = "";
            if (ac.classList.contains("active-all-comments") && postId === activePost) {
                ac.classList.remove("active-all-comments");
                activePost = null;
            } else {
                activePost = postId;

                // Fetch comments
                let data = new FormData();
                data.append("postId", postId);

                let comments = await new Promise((resolve) => {
                    fetch(commentsRoute, {
                            method: 'POST',
                            body: data
                        })
                        .then(r => r.json())
                        .then(d => {
                            resolve(d);
                        });
                });

                comments = comments.sort((a, b) => CompareUpdatedDate(a.creationDate, b.creationDate));

                // Open comments terminal
                ac.classList.add("active-all-comments");
                await new Promise((resolve) => {
                    setTimeout(async () => {
                        await printInTerminal("#all-comments-text", allCommentsTxt, "cmt", 10);
                        resolve();
                    }, 500);
                })

                // Add comment input
                await new Promise((resolve) => {
                    // Add comment input
                    let input = /*html*/ `
                        <input type="text" name="comment-input" id="comment-input" class="forum-input" autocomplete="false">
                    `;
                    let measurer = /*html*/ `
                        <span class="measurer"></span>
                    `;
                    let txts = document.querySelectorAll("span.comment-input-name");
                    txts[txts.length - 1].after(createElementFromHTML(input));
                    let commentInput = document.querySelector("#comment-input");
                    commentInput.after(createElementFromHTML(measurer));
                    commentInput.focus();
                    commentInput.classList.add("focused");
                    resizeInput();
                    resolve();
                });

                // Display comments
                let commentsZone = document.createElement("div");
                commentsZone.id = "comments-zone";
                document.querySelector(".comment-input-name").before(commentsZone);
                document.querySelector(".comment-input-name").before(document.createElement("br"));

                for (const c of comments) {
                    await new Promise((resolve) => {
                        setTimeout(() => {
                            commentsZone.appendChild(comment(c.creationDate, c.nickname, c.content))
                            resolve();
                        }, 100);
                    })
                }

                handleSendComment();
                terminalScrollBottom("#comment-input");
            }
        });
    })
}

function handleSendComment() {
    document.querySelector("#all-comments").addEventListener("keydown", async (e) => {
        let commentInput = document.querySelector("#comment-input");
        let cmtTxt = commentInput.value;
        if (e.code === "Enter" && cmtTxt.length > 0) {
            e.preventDefault();

            // Send comment text
            let data = new FormData();
            data.append("postId", activePost);
            data.append("sessionId", sessionID);
            data.append("commentTxt", cmtTxt);

            fetch(createCommentRoute, {
                    method: 'POST',
                    body: data
                })
                .then(r => r.json())
                .then(d => {
                    if (d.value === "SUCCESS") {
                        // Insert new comment
                        let commentsZone = document.querySelector("#comments-zone");
                        commentsZone.appendChild(comment(undefined, userNickname, cmtTxt))
                    }
                });

            // Reset input
            commentInput.value = "";
            commentInput.style.width = "0";
            commentInput.parentNode.querySelector(".measurer").textContent = "";

            // Update comments count number
            let el = document.querySelector(`.reaction-btns-zone[data-post="${activePost}"]`);
            let numComments = el.querySelector(".num-comments");
            numComments.textContent = String(parseInt(numComments.textContent) + 1);
            terminalScrollBottom("#comment-input");
        }
    });
}