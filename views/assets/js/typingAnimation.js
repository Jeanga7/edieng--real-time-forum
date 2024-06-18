let userNickname = null;
var txtDataToLoad = () => activeForm === "login" ? loginTxtsData : registerTxtsData
var otherBtn = () => activeForm === "login" ? document.querySelector("#register-btn") : document.querySelector("#login-btn")
let netcatTxt = (user, data) => {
    let res = [{
            txt: `$> nc --connect "${user}"`,
            c: "command",
            eol: true
        },
        {
            txt: "",
            c: "command",
            eol: true
        },
    ];
    if (Object.prototype.toString.call(data) === "[object Array]" && data.length > 0) {
        data.forEach(d => {
            let className = userNickname && d.sender === userNickname ? "bold-fluo" : "warning";
            res.push({
                txt: d.sender + ":",
                c: className,
                eol: false
            })
            res.push({
                txt: d.message,
                c: null,
                eol: true
            })
        })
    } else {
        res.push({
            txt: "No message.",
            c: "danger no-message",
            eol: true
        })
    }
    res.push({
        txt: userNickname + ":",
        c: "bold-fluo chat-input-name",
        eol: false
    })
    res.push({
        txt: " ",
        c: null,
        eol: false
    })
    return res;
};

window.onload = () => {
    LettersAnimation("#preloader h1", 0, undefined, false);
    if (!loggedIn) {
        authLayer();
    } else {
        feedsLayer();
    }
};

async function authLayer() {
    await loadDOM();
    btnBgPosition();
    await hidePreloader();
    await disableAllInputs();
    await startWelcomeTerminal();
    await startLogsTerminal();
    await shrinkLeftBox();
    await expandAuthRight();
    document.querySelector("#auth-right").style.flexGrow = "1";
    await startActiveTerminal(activeForm, txtDataToLoad());
    await enableAllInputs();
    firstFieldFocus(`#${activeForm}`);
    formToggle();
    document.addEventListener("keydown", authHandleEnter);
}

async function feedsLayer() {
    await loadDOM();
    await hidePreloader();
    await startFeedsTerminals();
}

function disableAllInputs() {
    return new Promise((resolve) => {
        document.querySelectorAll("input").forEach(e => {
            e.disabled = true;
        });
        resolve();
    });
}

function enableAllInputs() {
    return new Promise((resolve) => {
        document.querySelectorAll(".forum-input").forEach(e => {
            e.disabled = false;
        });
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
        setTimeout(async () => {
            document.querySelector("#preloader").classList.add("zoom-in");
            await new Promise((resolve) => {
                setTimeout(() => {
                    document.querySelector("#preloader").style.display = "none";
                    resolve();
                }, 500);
            })
            resolve();
        }, 2000);
    });
}

function loadDOM() {
    return new Promise((resolve) => {
        document.querySelector("#container").innerHTML = loggedIn ? feedsView : loginView;
        resolve();
    });
}

function startFeedsTerminals() {
    return new Promise((resolve) => {
        setTimeout(async () => {
            // Clear first Enter event listener
            document.removeEventListener("keydown", authHandleEnter);
            // Terminal Pop-up
            let w = document.querySelectorAll(".terminal-box").forEach(t => {
                t.style.opacity = "1";
                // t.style.height = "100%";
                t.style.width = "100%";
            });

            document.querySelector("#user-infos").style.right = "0";
            document.querySelector("#user-infos").style.top = "0";

            // Typing
            await new Promise((resolve) => {
                setTimeout(async () => {
                    typeInTxts("#user-infos-text", userInfosTxt, "ui", 10);
                    typeInTxts("#categories-text", categoriesTxt, "cat", 10);
                    typeInTxts("#actions-text", actionsTxt, "ac", 10);
                    typeInTxts("#all-users-text", allUsersTxt, "au", 10);
                    typeInTxts("#all-posts-text", allPostsTxt, "ap", 10);
                    resolve();
                }, 500)
            });

            // Categories List
            let cats = [];
            await new Promise((resolve) => {
                setTimeout(() => {
                    fetch(categoriesListRoute, {
                            method: 'GET'
                        })
                        .then(r => r.json())
                        .then(d => {

                            d.forEach(c => {
                                cats.push({
                                    txt: "   - " + c.name,
                                    c: "warning",
                                    eol: true
                                })
                            })
                            resolve();
                        })
                }, 600)
            });

            printInTerminal("#categories-text", cats, "cat", 10);

            AddActionsBtn();

            // Get User Data
            let userInfos = (await fetchUserData(userID))[0];
            userNickname = userInfos.nickname;

            // Get Posts
            let posts = await fetchPosts(userID);

            ShowUserInfos(userInfos);
            ShowAllPosts(posts);
            ListUsers();

            enableReactions();

            LettersAnimation("#user-fullname", 4000);
            LettersAnimation(".post-title", 5000, 50);
            LettersAnimation(".user-nickname");

            resolve();
        }, 500);
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
                await typeInTxts("#welcome-text", welcomeTxt, "wlc", 50);
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
                await typeInTxts("#logs-text", logsTxt, "logs");
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
        for (const d of txtData) {
            await typeInTxts(d.zone, d.txtsList, d.className);
        }

        resizeInput();
        firstFieldFocus(`#${terminal}`);

        await resetAuthForm();

        resolve();
    });
}

async function resetAuthForm() {
    document.removeEventListener("keydown", authHandleEnter);
    document.addEventListener("keydown", authHandleEnter);

    printed = {};
    checked = {};
    fieldLoaded = false;
}

async function typeInTxts(selector, txtsList, className, s = speed) {
    let elementsList = document.querySelectorAll(`${selector} .txt`);

    for (let x = 0; x < txtsList.length; x++) {
        await typing(elementsList[x], txtsList[x].txt, className, s);
    }
}

function resizeInput() {
    document.querySelectorAll(".forum-input").forEach(el => {
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
                        fieldLoaded = false;
                        resolve();
                    }, 750);
                })
            }
        });
    });
}

function btnBgPosition() {
    document.querySelectorAll(".btn-bg").forEach(el => {
        if (el.parentNode.id === "login-btn") {
            el.style.right = "0"
        } else {
            el.style.left = "0"
        }
    })
}

async function fetchUserData(userId) {
    return new Promise((resolve) => {
        let fd = new FormData();
        fd.append("id", userId);
        fetch(userDataRoute, {
                method: 'POST',
                body: fd
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            })
    });
}

async function fetchPosts(userId) {
    return new Promise((resolve) => {
        let fd = new FormData();
        fd.append("id", userId);
        fetch(postsRoute, {
                method: 'GET'
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            })
    });
}

async function fetchAllUsers() {
    return new Promise((resolve) => {
        fetch(allUsersRoute, {
                method: 'GET'
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            });
    });
}

async function fetchChat(chatUser) {
    return new Promise((resolve) => {
        let fd = new FormData();
        fd.append("type", "fetch");
        fd.append("with", chatUser);
        fetch(discussionsRoute, {
                method: 'POST',
                body: fd
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            });
    });
}

async function ShowUserInfos(userInfos) {
    let infosZone = document.createElement("div");
    infosZone.id = "user-infos-zone"

    // Profile picture
    let pp = document.createElement("div");
    pp.id = "user-profile-picture";
    pp.style.backgroundImage = "url('assets/img/profile/anonymous.png')";
    infosZone.appendChild(pp);

    // User Name
    let un = document.createElement("div");
    un.id = "user-name-zone";

    un.innerHTML = makeUserName(userInfos)
    infosZone.appendChild(un);

    // Statistics
    let st = document.createElement("div");
    st.id = "user-stats";
    st.innerHTML = userStats()
    infosZone.appendChild(st);

    document.querySelector("#user-infos .terminal-body").appendChild(infosZone);
    await new Promise((resolve) => {
        setTimeout(() => {
            infosZone.style.opacity = "1";
            resolve();
        }, 200)
    })
}

async function ShowAllPosts(posts) {
    let postsZone = document.createElement("div");
    postsZone.id = "posts-zone";

    posts.forEach((postData, i) => {
        // Post
        let post = document.createElement("div");
        post.className = "post";
        postImage = postData.PostPicture ? makePostImage(postData.image) : "";
        post.innerHTML = makePost(postData)
        postsZone.appendChild(post);
        setTimeout(() => {
            post.style.opacity = "1";
        }, 500 + i * 500);
    });

    document.querySelector("#all-posts .terminal-body").appendChild(postsZone);
}

async function ListUsers() {
    // Create user zone
    let usersZone = document.createElement("div");
    usersZone.id = "users-zone";
    document.querySelector("#all-users .terminal-body").appendChild(usersZone);

    await resetChatList();
    openWs();

    handleOpenChat();
    handleSendMsg();
}

async function resetChatList() {
    // Displaying users
    let allUsers = await fetchAllUsers();
    usersZone = document.querySelector("#users-zone");
    usersZone.innerHTML = "";
    var known = [];
    var unknown = [];

    const promises = allUsers.map(async (userData) => {
        // Seperate discussions
        let chat = await fetchChat(userData.Nickname);
        if (chat.length > 0) {
            known.push({
                data: userData,
                date: chat[0].updatedDate
            });
        } else {
            unknown.push({
                data: userData
            });
        }
    });

    await Promise.all(promises);

    known = known.sort((a, b) => CompareUpdatedDate(b.date, a.date));
    unknown = unknown.sort((a, b) => CompareNicknames(a.data.Nickname, b.data.Nickname));

    let listChats = (el, i) => {
        let user = document.createElement("div");
        user.className = "user";
        user.dataset.nickname = el.data.Nickname;
        user.innerHTML = makeUser(el.data);
        usersZone.appendChild(user);
        setTimeout(() => {
            user.style.opacity = "1";
        }, 500 + i * 500);
    }

    known.forEach(listChats);
    unknown.forEach(listChats);
}

let CompareUpdatedDate = (date1, date2) => (new Date(date1)) - (new Date(date2))
let CompareNicknames = (nickname1, nickname2) => nickname1.localeCompare(nickname2)

function AddActionsBtn() {
    let actionsZone = document.createElement("div");
    actionsZone.id = "actions-zone";
    actionsZone.innerHTML = actionsButtons
    document.querySelector("#actions .terminal-body").appendChild(actionsZone)
}

function LettersAnimation(selector, interval = 5000, speed = 30, loop = true) {
    let alphabet = 'abcdefghijklmnopqrstuvwxyz';
    let randomLetter = () => alphabet[Math.floor(Math.random() * 26)];

    document.querySelectorAll(selector).forEach(el => {
        el.dataset.name = el.textContent;

        const repeat = setInterval(() => {
            let iter = 0;
            const interval = setInterval(() => {
                el.textContent = el.textContent
                    .split("")
                    .map((letter, i) => {
                        if (i < iter) {
                            return el.dataset.name[i];
                        }
                        return randomLetter();
                    })
                    .join("");

                if (iter >= el.dataset.name.length) clearInterval(interval)

                iter += 1 / 3;
            }, speed);

            if (!loop) clearInterval(repeat);
        }, interval);
    });
}

function handleOpenChat() {
    document.querySelectorAll(".user").forEach(el => {
        el.addEventListener("click", async (e) => {
            let chatUser = e.target.closest(".user").dataset.nickname;
            let nc = document.querySelector("#netcat");

            if (chatWith === chatUser) {
                if (nc.classList.contains("active-chat")) nc.classList.remove("active-chat");
                chatWith = null
            } else {
                if (!nc.classList.contains("active-chat")) nc.classList.add("active-chat");
                chatWith = chatUser

                // Get chat
                let chat = await fetchChat(chatUser);
                let msgs = [];
                if (chat.length === 1) {
                    msgs = JSON.parse(chat[0].messages);
                }

                let netcat = document.querySelector("#netcat-text")
                netcat.innerHTML = "";

                await new Promise((resolve) => {
                    setTimeout(async () => {
                        await printInTerminal("#netcat-text", netcatTxt(chatUser, msgs), "nc", 10);
                        resolve();
                    }, 700)
                });

                await new Promise((resolve) => {
                    setTimeout(() => {
                        // Add chat input
                        let input = /*html*/ `
                            <input type="text" name="chat-input" id="chat-input" class="forum-input">
                        `;
                        let measurer = /*html*/ `
                            <span class="measurer"></span>
                        `;
                        let txts = netcat.querySelectorAll("span.bold-fluo");
                        txts[txts.length - 1].after(createElementFromHTML(input));
                        let chatInput = document.querySelector("#chat-input");
                        chatInput.after(createElementFromHTML(measurer));
                        chatInput.focus();
                        resizeInput();
                        resolve();
                    }, 1000);
                });

                await new Promise((resolve) => {
                    setTimeout(() => {
                        document.querySelector(".chat-input-name").addEventListener("click", (e) => {
                            document.querySelector("#chat-input").focus();
                        });
                        resolve();
                    }, 200);
                });
            }
        });
    });
}

function handleSendMsg() {
    document.addEventListener("keydown", async (e) => {
        let chatInput = document.querySelector("#chat-input");
        let chatInputName = document.querySelector(".chat-input-name");
        if (e.code === "Enter" && chatInput.value.length > 0) {
            e.preventDefault();
            let msg = chatInput.value;

            let fd = new FormData();
            fd.append("type", "send");
            fd.append("to", chatWith);
            fd.append("message", chatInput.value);

            fetch(discussionsRoute, {
                    method: 'POST',
                    body: fd
                })
                .then(r => r.json())
                .then(async d => {
                    if (d.value === "SUCCESS") {
                        // Clear field
                        chatInput.value = "";
                        chatInput.style.width = "0";
                        chatInput.parentNode.querySelector(".measurer").textContent = "";

                        // Insert new message
                        chatInputName.before(makeSpan("bold-fluo new-msg-name"));

                        let name = document.querySelectorAll(".new-msg-name");
                        name[name.length - 1].textContent = userNickname + ":"
                        name[name.length - 1].after(makeSpan("new-msg-txt"))

                        let txt = document.querySelectorAll(".new-msg-txt");
                        txt[txt.length - 1].textContent = msg;

                        txt[txt.length - 1].after(makeBr())
                    }
                });

        }
    });
}

function openWs() {
    ws = new WebSocket("ws://localhost:8989/api/discussions");
    ws.onopen = () => {
        console.log("Connected to WebSocket...");
    };
    ws.onmessage = (e) => {
        const message = JSON.parse(e);
        console.log("New message: ", message);
        // if (message.type === 'history') {
        //     loadHistory(message.messages.reverse());
        // } else if (message.type == "typing") {
        //     if (currentdiscussion === message.sender) {
        //         showTyping();
        //     }
        // } else if (message.type == "new_user") {
        //     handleNewUser(message.sender)
        // } else {
        //     handleMessage(message);
        // }
    }
}