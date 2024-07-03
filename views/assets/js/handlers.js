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
                // console.log(response.message);
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
                    userToOnline.querySelector(".user-typing").classList.add("txt-online");
                }
                break;
            case "offline":
                let userToOffline = document.querySelector(`[data-nickname="${response.message}"]`);
                if (userToOffline) {
                    userToOffline.querySelector(".user-status").classList.remove("status-online");
                    userToOffline.querySelector(".user-nickname").classList.remove("txt-online");
                    userToOffline.querySelector(".user-notifications").classList.remove("txt-online");
                    userToOffline.querySelector(".user-typing").classList.remove("txt-online");
                }
                break;
            case "typing":
                displayTyping(response);
        }
    }
}

function handleOpenChat() {
    document.querySelectorAll(".user").forEach(el => {
        el.addEventListener("click", throttle(async (e) => {
            let netcat = document.querySelector("#netcat-text")

            let chatUser = e.target.closest(".user").dataset.nickname;
            let nc = document.querySelector("#netcat");

            await new Promise((resolve) => {
                netcat.innerHTML = "";
                resolve();
            });

            if (chatWith === chatUser) {
                if (nc.classList.contains("active-chat")) nc.classList.remove("active-chat");
                chatWith = null;
                chatChunks = [];
            } else {
                if (!nc.classList.contains("active-chat")) nc.classList.add("active-chat");
                chatWith = chatUser;

                // Get chat
                let chat = await fetchChat(chatUser);
                let msgs = [];
                if (chat.length === 1) {
                    msgs = JSON.parse(chat[0].messages);
                }
                chatChunks = chunkList(msgs, 10);
                msgs = chatChunks.shift();

                await new Promise((resolve) => {
                    setTimeout(async () => {
                        await printInTerminal("#netcat-text", netcatTxt(chatUser, msgs), "nc", 5);
                        resolve();
                    }, 700)
                });

                await addInputBeforeCursor("netcat")

                await new Promise((resolve) => {
                    setTimeout(() => {
                        focusInput();
                        terminalScrollBottom("#netcat-input");
                        resolve();
                    }, 200);
                });

                await handleResetNotifications();

                handleLoadMoreMsg();
            }
        }, "open-chat", 3000));
    });
}

function handleResetNotifications() {
    return new Promise((resolve) => {
        // Reset notifications DOM
        ResetNotifsNumber();

        // Reset notifications request
        let data = new FormData();
        data.append("sender", chatWith);
        data.append("dest", userNickname);

        fetch(resetNotifsRoute, {
                method: 'POST',
                body: data
            })
            .then(r => r.json())
            .then(d => {
                resolve();
            })
    })
}

function handleSendMsg() {
    document.querySelector("#netcat").addEventListener("keydown", async (e) => {
        let chatInput = document.querySelector("#netcat-input");
        if (chatInput && chatWith !== null) {
            if (chatInput.value.length > 0) {
                if (e.code === "Enter") {
                    e.preventDefault();
                    ws.send(JSON.stringify({
                        "from": userNickname,
                        "to": chatWith,
                        "type": "message",
                        "message": chatInput.value,
                        "signature": sessionID
                    }))
                    // Reset input
                    let msg = chatInput.value;
                    chatInput.value = "";
                    chatInput.style.width = "0";
                    chatInput.parentNode.querySelector(".measurer").textContent = "";
                    // Insert new message
                    let chatInputName = document.querySelector(".netcat-input-name");
                    insertMsg(chatInputName, false, "bold-fluo", undefined, userNickname, msg)
                    terminalScrollBottom("#netcat-input");
                    let user = document.querySelector(`.user[data-nickname="${chatWith}"]`);
                    moveElementToParentStart(user);
                    await handleResetNotifications();
                } else {
                    throttle(sendTypingNotif, "sendTyping", 1200)();
                }
            }
        }
    });
}

function sendTypingNotif() {
    ws.send(JSON.stringify({
        "from": userNickname,
        "to": chatWith,
        "type": "typing",
        "message": "",
        "signature": sessionID
    }))
}

function displayTyping(response) {
    let userBlock = document.querySelector(`[data-nickname="${response.from}"]`);
    if (userBlock) {
        userBlock.querySelector(".user-typing").classList.remove("not-displayed");
        clearTimeout(timers["typing"]);
        timers["typing"] = setTimeout(() => {
            userBlock.querySelector(".user-typing").classList.add("not-displayed");
        }, 2500);
    }
}

function handleLoadMoreMsg() {
    document.querySelector("#netcat").addEventListener("scroll", debounceScroll);
}

function handleOpenNewPost() {
    document.querySelector("#add-post-btn").addEventListener("click", throttle(async (e) => {
        await emptyNewPost();
        let np = document.querySelector("#new-post");
        let npl = document.querySelector("#new-post-logs");
        if (np.classList.contains("active-new-post")) {
            np.classList.remove("active-new-post");
            npl.classList.remove("active-new-post-logs");
        } else {
            np.classList.add("active-new-post");
            npl.classList.add("active-new-post-logs");
            await new Promise((resolve) => {
                setTimeout(async () => {
                    await printInTerminal("#new-post-text", newPostTxtIntro, "np", 10);
                    insertNewPostFieldRow("title")
                    await printInTerminal("#new-post-title-row", newPostTxtField, "np", 10);
                    await printInTerminal("#new-post-logs-text", newPostLogsTxt, "npl", 10);
                    resolve();
                }, 500);
            })

            await addInputBeforeCursor("new-post", undefined, "title");
            handleAddPost();
        }
    }, "open-new-post", 500));
}

function handleOpenComments() {
    document.querySelectorAll(".comment-btn").forEach(el => {
        el.addEventListener("click", throttle(async (e) => {
            await emptyComments();
            let postId = el.parentElement.dataset.post;
            let ac = document.querySelector("#all-comments");
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
        }, "open-comments", 500));
    })
}

function emptyComments() {
    return new Promise((resolve) => {
        document.querySelector("#all-comments-text").textContent = "";
        resolve();
    })
}

function emptyNewPost() {
    return new Promise((resolve) => {
        document.querySelector("#new-post-text").textContent = "";
        document.querySelector("#new-post-logs-text").textContent = "";
        resolve();
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

function handleLogout() {
    document.querySelector("#logout-btn").addEventListener("click", async (e) => {
        ws.close();
        await new Promise((resolve) => {
            setTimeout(() => {
                window.location.href = '/logout';
                resolve();
            }, 100);
        })
    });
}

function handleAddPost() {
    printed = {};
    checked = {};
    document.querySelector("#new-post").addEventListener("keydown", async (e) => {
        if (e.code === "Enter") {
            e.preventDefault();
            let data = new FormData();
            let finalData = new FormData();
            let fieldset = ["author", "title", "content", "category"];

            [data, checked, printed] = await Validation(sessionID, createPostRoute, "author", data, checked, printed, "#new-post-logs-text", "new-post-logs");
            let iniData = data;

            for (let i = 1; i < fieldset.length; i++) {
                var field = fieldset[i];
                var prev = fieldset[i - 1];
                var next = i < fieldset.length ? fieldset[i + 1] : undefined;

                if (!printed[prev] || !checked[prev]) {
                    return
                }

                [data, checked, printed] = await ValidateField(createPostRoute, `#new-post-${field}`, field, iniData, checked, printed, "#new-post-logs-text", "new-post-logs");
                if (!checked[field]) {
                    return
                } else {
                    document.querySelector("#new-post-" + field).classList.remove("focused");
                }

                if (next && printed[field] && !printed[next]) {
                    let fieldElement = document.querySelector("#new-post-" + field);
                    fieldElement.setAttribute("contenteditable", false);
                    await printAddPostField(next)
                    printed[next] = true
                }
            }

            let allChecked = fieldset.reduce((res, check) => res && checked[check], true);

            finalData.append("author", sessionID);
            for (const field of fieldset.slice(-3)) {
                finalData.append(field, document.querySelector("#new-post-" + field).textContent)
            }

            if (allChecked) {
                // alert("All checked.")
                // Display the values
                for (const pair of finalData.entries()) {
                    console.log(pair[0], pair[1]);
                }

                fetch(createPostRoute, {
                    method: 'POST',
                    body: finalData
                })
                .then(r=>r.json())
                .then(async d=>{
                    let output = [{
                        txt: "Post published successfully!",
                        c: "bold-fluo command",
                        eol: true
                    }, ];
                    await printInTerminal("#new-post-logs-text", output, "new-post-logs", 10);
                })

            }
        }
    });
}

async function printAddPostField(fieldName) {
    await new Promise((resolve) => {
        setTimeout(async () => {
            let newPostTxt = [{
                    txt: "Enter the ",
                    c: null,
                    eol: false
                },
                {
                    txt: fieldName,
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
            ];

            insertNewPostFieldRow(fieldName);
            await printInTerminal(`#new-post-${fieldName}-row`, newPostTxt, "np", 10);

            await addInputBeforeCursor("new-post", undefined, fieldName);
            moveToEndOfSpan("#new-post-" + fieldName);

            resolve();
        }, 50);
    })

}