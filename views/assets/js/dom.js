// FUNCTIONS
function loopSpans(txtsList) {
    let res = "";
    txtsList.forEach(el => {
        res += `<span class="txt ${el.c ?? ""}"></span>`;
        res += el.eol ? "<br>" : "";
    });
    return res;
}

function makeSpan(cl) {
    let span = document.createElement("span");
    span.className = "txt "+cl;
    return span;
}

let makeBr = () => document.createElement("br")

async function printInTerminal(selector, data, typingClass, s=speed) {
    let t = document.querySelector(selector);
    for (let i = 0; i < data.length; i++) {
        const d = data[i];
        t.innerHTML += `<span class="txt ${d.c ?? ""}"></span>` + (d.eol ? "<br>" : "");
        let txts = t.querySelectorAll(".txt");
        await typing(txts[txts.length - 1], d.txt, typingClass, s);
        terminalScrollBottom(selector);       
    }
}

function createElementFromHTML(htmlString) {
    var div = document.createElement('div');
    div.innerHTML = htmlString.trim();
  
    // Change this to div.childNodes to support multiple top-level nodes.
    return div.firstChild;
  }

let generateSpans = (t) => loopSpans(t[0].txtsList);

let cliRow = (title, txtsList, section) => `
<!-- ${capitalize(section)} row -->
<span id="${title}-${section}">
    ${generateSpans(txtsList)}
</span>`;

let terminalRow = (terminal, title, type) => `
<!-- ${title} row -->
<span class="auth-row" id="${terminal}-${title.toLowerCase()}-row">
    <span class="txt"></span>
    <span class="txt bold-fluo"></span>
    <span class="txt"></span>
    <input type="${type}" name="${terminal}-${title.toLowerCase()}" id="${terminal}-${title.toLowerCase()}" class="forum-input" autocomplete="false">
    <span class="txt"></span>
    <span class="measurer"></span>
</span>
<br>`;

function makeFields(terminal, data) {
    let dom = "";
    for (const r of data) {
        dom += terminalRow(terminal, r.name, r.type);
    }
    return dom;
}

// cliRow(name, data, "cli")
let makeTerminal = (name, rows, pClass) => `
<div id="${name}" class="terminal-box">
<div class="terminal-head">&lt;${name.toUpperCase()}/&gt;</div>
<div class="terminal-body">
    <p id="${name}-text" class="${pClass}">
    ${rows}
    </p>
</div>
</div>`;

let loginRows = cliRow("login", loginTxtsData, "intro") + makeFields("login", loginFields);
const loginTerminal = makeTerminal("login", loginRows, "infos-txt")

let registerRows = cliRow("register", registerTxtsData, "intro") + makeFields("register", registerFields)
const registerTerminal = makeTerminal("register", registerRows, "infos-txt")

let welcomeRows = loopSpans(welcomeTxt);
const welcomeTerminal = makeTerminal("welcome", welcomeRows, "welcome-text")

let logsRows = loopSpans(logsTxt);
const logsTerminal = makeTerminal("logs", logsRows, "infos-txt")

var loginView = `
<div id="auth-left" class="auth-block">
    ${welcomeTerminal}

    ${logsTerminal}
</div>

<div id="auth-right" class="auth-block">
    <div id="terminal-zone">
        ${loginTerminal}
    </div>
    
    <div id="auth-btns">
        <span id="login-btn" class="auth-btn active"><span class="btn-bg"></span>LOGIN</span>
        <span id="register-btn" class="auth-btn"><span class="btn-bg"></span>REGISTER</span>
    </div>
</div>`


let userInfosRows = loopSpans(userInfosTxt);
const userInfosTerminal = makeTerminal("user-infos", userInfosRows, "infos-txt")

let categoriesRows = loopSpans(categoriesTxt);
const categoriesTerminal = makeTerminal("categories", categoriesRows, "infos-txt")

let actionsRows = loopSpans(actionsTxt);
const actionsTerminal = makeTerminal("actions", actionsRows, "infos-txt")

let allUsersRows = loopSpans(allUsersTxt);
const allUsersTerminal = makeTerminal("all-users", allUsersRows, "infos-txt")

let allPostsRows = loopSpans(allPostsTxt);
const allPostsTerminal = makeTerminal("all-posts", allPostsRows, "infos-txt")

const newPostTerminal = makeTerminal("new-post", "", "infos-txt")

const allCommentsTerminal = makeTerminal("all-comments", "", "infos-txt")

const netcatTerminal = makeTerminal("netcat", "", "infos-txt")

var feedsView = /*html*/ `
<div id="feeds-grid">
    <div id="feeds-left" class="feeds-col">
        ${userInfosTerminal}
        ${categoriesTerminal}
        ${actionsTerminal}
    </div>
    <div id="feeds-middle" class="feeds-col">
        ${newPostTerminal}
        ${allPostsTerminal}
        ${allCommentsTerminal}
    </div>
    <div id="feeds-right" class="feeds-col">
        ${allUsersTerminal}
        ${netcatTerminal}
    </div>
</div>
`

let actionsButtons = /*html*/`
<span id="add-post-btn">Add post</span>
<span id="logout-btn" onclick="window.location.href='/logout'">Logout</span>
`

let makePost = (postData) => /*html*/`
<hr class="post-separator">
<span class="post-title">${postData.Title}</span>
<span class="post-time-author">[${formatDate(postData.CreationDate)}] by ${postData.AuthorName}</span>
<p class="post-text">${postData.Content}</p>
${postImage}
<div class="reaction-btns-zone" data-post="${postData.Id}">
    <span class="reaction like-btn like-dislike ${postData.Liked ? "reaction-active" : ""}">[<span class="reaction-check">${postData.Liked ? "X" : " "}</span>]Like(<span class="reaction-count">${postData.Likes}</span>)</span>
    <span class="reaction dislike-btn like-dislike ${postData.Disliked ? "reaction-active" : ""}">[<span class="reaction-check">${postData.Disliked ? "X" : " "}</span>]Dislike(<span class="reaction-count">${postData.Dislikes}</span>)</span>
    <span class="reaction comment-btn">$>Comment(${postData.NumComments})</span>
</div>
`

let makeUser = (userData) => /*html*/`
    <span class="user-status ${userData.Online ? 'status-online' : ''}"></span>
    <span class="user-nickname ${userData.Online ? 'txt-online' : ''}">${userData.Nickname}</span>
    <span class="user-notifications ${userData.Online ? 'txt-online' : ''}">(${userData.NbrNotifications})</span>
`

let makePostImage = (img) => /*html*/`
<div class="post-image" style="background-image:url('assets/img/posts/${img}');"></div>
`

let userStats = () => /*html*/ `
<span id="stats-posted">0 posted</span>
<span id="stats-reacted">0 reacted</span>
<span id="stats-replied">0 replied</span>`

let makeUserName = (userInfos) => /*html*/ `
<span id="user-fullname">${userInfos.firstname + " " + userInfos.lastname}</span>
<span id="user-nickname">@${userInfos.nickname}</span>`