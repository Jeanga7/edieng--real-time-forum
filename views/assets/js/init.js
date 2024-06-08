const speed = 15;
let fieldLoaded = false;
var activeForm = "login";
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
    {txt: "$> Page loaded successfully", c: null, eol: false}
];