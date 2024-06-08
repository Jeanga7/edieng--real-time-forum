function generateSpans(t) {
    let res = "";
    t[0].txtsList.forEach(el=>{
        res += `<span class="txt ${el.c ?? ""}"></span>`;
        res += el.eol ? "<br>" : "";
    });
    return res;
}

let introRow = (title, txtsList) => `<!-- Intro row -->
<span id="${title}-intro">
    ${generateSpans(txtsList)}
</span>`

let terminalRow = (terminal, title, type) => `<!-- ${title} row -->
<span class="auth-row" id="${terminal}-${title.toLowerCase()}-row">
    <span class="txt"></span>
    <span class="txt bold-fluo"></span>
    <span class="txt"></span>
    <input type="${type}" name="${terminal}-${title.toLowerCase()}" id="${terminal}-${title.toLowerCase()}" class="auth-input">
    <span class="txt"></span>
    <span class="measurer"></span>
</span>
<br>`

let loginFields = [
    {name: "Nickname", type: "text"},
    {name: "Password", type: "password"}
];
let registerFields = [
    {name: "Firstname", type: "text"},
    {name: "Lastname", type: "text"},
    {name: "Nickname", type: "text"},
    {name: "Email", type: "email"},
    {name: "Password", type: "password"},
    {name: "Age", type: "number"},
    {name: "Gender", type: "text"}
];

function makeFields(terminal, data) {
    let dom = "";
    for (const r of data) {
        dom += terminalRow(terminal, r.name, r.type);
    }
    return dom;
}

const loginTerminal = `
<div id="login" class="terminal-box">
                <div class="terminal-head">&lt;LOGIN/&gt;</div>
                <div class="terminal-body">
                    <form id="login-form">
                        <p id="login-text" class="infos-txt">
                            <!-- Intro row -->
                            ${introRow("login", loginTxtsData)}
                            ${makeFields("login", loginFields)}
                        </p>
                    </form>
                </div>
            </div>
`

const registerTerminal = `
<div id="register" class="terminal-box not-displayed">
                <div class="terminal-head">&lt;REGISTER/&gt;</div>
                <div class="terminal-body">
                    <form id="register-form" enctype="multipart/form-data">
                        <p id="register-text" class="infos-txt">
                            <!-- Intro row -->
                            ${introRow("register", registerTxtsData)}
                            ${makeFields("register", registerFields)}
                        </p>
                    </form>
                </div>
            </div>
`

var loginView = `
        <div id="auth-left" class="auth-block">
            <div id="welcome" class="terminal-box">
                <div class="terminal-head">&lt;WELCOME/&gt;</div>
                <div class="terminal-body">
                    <p id="welcome-text">
                        <span class="txt"></span><span class="txt bold-fluo"></span><span class="txt"></span><br>
                        <span class="txt"></span><br>
                        <span class="txt bold-fluo"></span><span class="txt"></span>
                    </p>
                </div>
            </div>

            <div id="logs" class="terminal-box">
                <div class="terminal-head">&lt;LOGS/&gt;</div>
                <div class="terminal-body">
                    <p id="logs-text" class="infos-txt">
                        <span class="txt"></span><br>
                    </p>
                </div>
            </div>
        </div>

        <div id="auth-right" class="auth-block">
            <div id="terminal-zone">
                ${loginTerminal}
            </div>
            
            <div id="auth-btns">
                <span id="login-btn" class="auth-btn active"><span class="btn-bg"></span>LOGIN</span>
                <span id="register-btn" class="auth-btn"><span class="btn-bg"></span>REGISTER</span>
            </div>
        </div>
`