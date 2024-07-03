let userNickname = null;
var txtDataToLoad = () => activeForm === "login" ? loginTxtsData : registerTxtsData
var otherBtn = () => activeForm === "login" ? document.querySelector("#register-btn") : document.querySelector("#login-btn")


window.onload = async () => {
    LettersAnimation("#preloader h1", 0, undefined, false);
    if (!loggedIn) {
        await authLayer();
    } else {
        await feedsLayer();
    }
    focusInput();
    disableTab();
};

function disableTab() {
    document.addEventListener("keydown", (e) => {
        if (e.code === "Tab") {
            e.preventDefault();
        }
    });
}

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
            let w = document.querySelectorAll(".terminal-box, #new-post-zone");
            for (const t of w) {
                await new Promise((resolve) => {
                    setTimeout(() => {
                        t.style.opacity = "1";
                        // t.style.height = "100%";
                        t.style.width = "100%";
                        resolve();
                    }, 100);
                })
            }

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
            handleOpenNewPost();
            handleOpenComments();

            ListUsers();

            enableReactions();

            LettersAnimation("#user-fullname", 4000);
            LettersAnimation(".post-title", 5000, 50);
            LettersAnimation(".user-nickname");

            handleLogout();

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



function resizeInput() {
    document.querySelectorAll(".forum-input").forEach(el => {
        el.parentNode.addEventListener("click", (e) => {
            el.focus();
        });

        el.addEventListener("input", (e) => {
            resetInputWidth(el);
        });
    });
}

function resetInputWidth(el) {
    let m = el.parentNode.querySelector(".measurer");
    m.textContent = el.value;
    let mWidth = m.offsetWidth;
    el.style.width = (mWidth + 15) + "px";
}

function firstFieldFocus(selector) {
    removeAllFocusedClass();
    console.log(document.querySelector(selector));
    let firstField = document.querySelector(selector).querySelectorAll(`.form-input`)[0];
    firstField.focus();
    firstField.classList.add("focused");
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



async function ListUsers() {
    // Create user zone
    let usersZone = document.createElement("div");
    usersZone.id = "users-zone";
    document.querySelector("#all-users .terminal-body").appendChild(usersZone);

    await resetChatList();

    handleWebsocket();
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

function AddActionsBtn() {
    let actionsZone = document.createElement("div");
    actionsZone.id = "actions-zone";
    actionsZone.innerHTML = actionsButtons
    document.querySelector("#actions .terminal-body").appendChild(actionsZone)
}

function addInputBeforeCursor(terminal, type = "input", fieldName = "input") {
    return new Promise((resolve) => {
        // Add chat input
        let input = /*html*/ `
                        <span type="text" name="${terminal}-${fieldName}" id="${terminal}-${fieldName}" class="forum-input" autocomplete="false" contenteditable></span>
                    `;

        let space = /*html*/`
            <span style="width:8px;height:8px;display:inline-block;"></span>
        `

        let cursor = document.querySelector(`#${terminal}-text`).querySelector(".blinking");
        if (cursor) {
            cursor.parentElement.before(createElementFromHTML(input));
            cursor.parentElement.before(createElementFromHTML(space));
            let inputElement = document.querySelector(`#${terminal}-${fieldName}`);
            inputElement.focus();
            removeAllFocusedClass();
            inputElement.classList.add("focused");
        }
        resolve();
    });
}

function ResetNotifsNumber() {
    let user = document.querySelector(`[data-nickname="${chatWith}"]`);
    user.querySelector(".user-notifications").textContent = `(0)`;
}

function focusInput() {
    document.querySelectorAll(".terminal-box").forEach(el => {
        el.addEventListener("click", (e) => {
            let focused = el.querySelectorAll(".focused");
            if (focused.length > 0) {
                let field = focused[0]
                if (field.tagName.toLowerCase() === "span" && field.isContentEditable) {
                    moveToEndOfSpan(".focused");
                } else {
                    field.focus();
                }
            }
        });
    });
}

let debounceScroll = (e) => debounce(loadOldMsg, "chat", 800)(e);

async function loadOldMsg(e) {
    if (e.target.scrollTop < 100 && chatChunks.length > 0) {
        let msgs = reverseArray(chatChunks.shift());

        let oldMsgs = () => document.querySelectorAll(".old-msg-name");
        if (oldMsgs().length > 0) {
            for (const msg of msgs) {
                await new Promise((resolve) => {
                    setTimeout(() => {
                        let className = msg.sender !== userNickname ? "warning" : "bold-fluo"
                        insertMsg(oldMsgs()[0], true, className, msg.date, msg.sender, msg.message);
                        resolve();
                    }, 100);
                })
            }
        }
        e.target.scrollTop = e.target.scrollHeight;
    }
}