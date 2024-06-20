// Initialization
const speed = 15;
var fieldLoaded = false;
var activeForm = "login";
var chatWith = null;
var chatChunks = [];

var printed = {};
var checked = {};
const authRoute = '/api/auth';
const userDataRoute = '/api/users';
const postsRoute = '/api/posts';
const categoriesListRoute = '/api/categories-list';
const allUsersRoute = '/api/all-users';
const discussionsRoute = '/api/discussions';
const port = ":8989"

let ws = new WebSocket(`ws://localhost${port}/messages`);
// let alphabet = 'abcdefghijklmnopqrstuvwxyz';

// let randomLetter = () => alphabet[Math.floor(Math.random() * 26)];

let capitalize = (str) => str.charAt(0).toUpperCase() + str.substr(1);

let formatDate = (d=undefined) => (d ? new Date(d) : new Date()).toLocaleString('fr-FR', {
    // year: '2-digit',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    hour12: false // Use 24-hour time format
});

// let reverseArray = (arr) => arr.length > 1
//                             ? [arr[]]

// function chunkList(arr, len) {
//     let newArr = [];
//     for (arr.length > len) {
//         newArr.push(arr.slice(len - arr.length))
//     }
// }

// let chunk = (arr, len) => [arr.length>=len ? ]

let loginTxtsData = [{
        zone: "#login-intro",
        className: "lgn",
        txtsList: [
            {txt: "$> help", c: "command", eol: true},
            {txt: "These infos are needed for login :", c: null, eol: true},
            {txt: " - nickname or email", c: "command", eol: false},
            {txt: " (lowercase)", c: "command warning", eol: true},
            {txt: " - password", c: "command", eol: true},
            {txt: "", c: null, eol: true},
            {txt: "User should be registered first before trying to register.", c: "warning", eol: true},
            {txt: "Check logs to see if inputed values are valid.", c: null, eol: true},
            {txt: "", c: null, eol: true},
            {txt: "", c: null, eol: true},
            {txt: "$> login", c: "command", eol: true},
        ]
    },
    {
        zone: "#login-nickname-row",
        className: "lgn",
        txtsList: [
            {txt: "Enter your ", c: null, eol: false},
            {txt: "nickname or email", c: "bold-fluo", eol: false},
            {txt: ": ", c: null, eol: false},
            {txt: " ", c: null, eol: false},
        ]
    }
];

let registerTxtsData = [{
        zone: "#register-intro",
        className: "rgst",
        txtsList: [
            {txt: "$> help", c: "command", eol: true},
            {txt: "These infos are needed for register :", c: null, eol: true},
            {txt: " - firstname", c: "command", eol: true},
            {txt: " - lastname", c: "command", eol: true},
            {txt: " - nickname", c: "command", eol: false},
            {txt: " (alphanumeric only)", c: "command warning", eol: true},
            {txt: " - email", c: "command", eol: true},
            {txt: " - password", c: "command", eol: false},
            {txt: " (+5 chars)", c: "command warning", eol: true},
            {txt: " - age", c: "command", eol: false},
            {txt: " (numeric)", c: "command warning", eol: true},
            {txt: " - gender", c: "command", eol: false},
            {txt: " (M or F)", c: "command warning", eol: true},
            {txt: "Check logs to see if inputed values are valid.", c: null, eol: true},
            {txt: "", c: null, eol: true},
            {txt: "", c: null, eol: true},
            {txt: "$> register", c: "command", eol: true},
        ]
    },
    {
        zone: "#register-firstname-row",
        className: "rgst",
        txtsList: [
            {txt: "Enter your ", c: null, eol: false},
            {txt: "first name", c: "bold-fluo", eol: false},
            {txt: ": ", c: null, eol: false},
            {txt: " ", c: null, eol: false},
        ]
    }
];

let welcomeTxt = [
    {txt: "Hello ", c: null, eol: false},
    {txt: "talent", c: "bold-fluo", eol: false},
    {txt: " !", c: null, eol: true},
    {txt: "Welcome to our", c: null, eol: true},
    {txt: "Real Time Forum", c: "bold-fluo", eol: false},
    {txt: ".", c: null, eol: false},
];

let logsTxt = [
    {txt: "$> Page loaded successfully", c: null, eol: true}
];

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

let userInfosTxt = [
    {txt: "$> rtf --user-infos", c: "command", eol: true}
];

let categoriesTxt = [
    {txt: "$> rtf --categories", c: "command", eol: true},
    {txt: "", c: null, eol: true},
    {txt: "List of categories:", c: null, eol: true},
];

let actionsTxt = [
    {txt: "$> rtf --actions", c: "command", eol: true},
];

let allUsersTxt = [
    {txt: "$> rtf --all-users", c: "command", eol: true},
];

let allPostsTxt = [
    {txt: "$> rtf --all-posts", c: "command", eol: true},
    {txt: "", c: null, eol: true},
];

async function authHandleEnter(e) {
    if (e.code === "Enter") {
        e.preventDefault();

        // Verify nickname
        let data = new FormData();
        data.append("form", activeForm);
        if (activeForm === "login") {
            [data, checked, printed] = await ValidateField(authRoute, "#login-nickname", "login", data, checked, printed);

            if (fieldLoaded) {
                // Verify password
                [data, checked, printed] = await ValidateNextField(authRoute, "login", "password", "#login-nickname", "#login-password", "#login-password-row", data, checked, printed);

                if (checked["login"] && checked["password"]) {
                    await switchToFeeds();
                }
            }
        } else if (activeForm === "register") {
            let prev = ""
            let fieldset = ["firstname", "lastname", "nickname", "email", "password", "age", "gender"];

            for (let i = 0; i < fieldset.length; i++) {
                const field = fieldset[i];
                checked[field] = false;
                if (i === 0) {
                    [data, checked, printed] = await ValidateField(authRoute, `#register-${field}`, field, data, checked, printed);
                } else if (fieldLoaded && document.querySelector(`#register-${prev}`).value.length > 0) {
                    [data, checked, printed] = await ValidateNextField(authRoute, prev, field, `#register-${prev}`, `#register-${field}`, `#register-${field}-row`, data, checked, printed);
                }
                prev = field;
            }

            let allChecked = fieldset.reduce((res, check) => res && checked[check], true);

            if (allChecked) {
                await switchToFeeds();
            }
        }
    } else if (e.code === "Tab") {
        e.preventDefault();
    }
}

// let limitWord = 10;
// let limitSentence = 50;
// document.querySelectorAll("input").forEach(el=>{
//     el.addEventListener("input", (e)=>{
//         let check = el.value.split(" ").reduce((prev, next)=>prev && length(next) <= limitWord, true)
//     })
// })