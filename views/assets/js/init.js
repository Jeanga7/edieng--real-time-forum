// Initialization
const speed = 15;
let fieldLoaded = false;
var activeForm = "login";

let printed = {};
let checked = {};
let authRoute = '/api/auth';
let userDataRoute = '/api/users';
let postsRoute = '/api/posts';
let categoriesListRoute = '/api/categories-list';

let capitalize = (str) => str.charAt(0).toUpperCase() + str.substr(1);

let loginTxtsData = [{
        zone: "#login-intro",
        className: "lgn",
        txtsList: [
            {txt: "$> help", c: "command", eol: true},
            {txt: "These infos are needed for login :", c: null, eol: true},
            {txt: " - nickname or email", c: "command", eol: true},
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
            {txt: " (alphanumeric only)", c: "warning", eol: true},
            {txt: " - email", c: "command", eol: true},
            {txt: " - password", c: "command", eol: false},
            {txt: " (+5 chars)", c: "warning", eol: true},
            {txt: " - age", c: "command", eol: false},
            {txt: " (numeric)", c: "warning", eol: true},
            {txt: " - gender", c: "command", eol: false},
            {txt: " (M or F)", c: "warning", eol: true},
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

let netcatTxt = [
    {txt: `$> nc --connect "Tony Stark"`, c: "command", eol: true},
];