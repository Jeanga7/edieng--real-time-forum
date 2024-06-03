var Chatopened = false
document.addEventListener('DOMContentLoaded', (event) => {
    openWs();

    document.querySelectorAll('.discussion').forEach(user => {
        user.addEventListener('click', () => {
            const username = user.textContent.split(' - ')[0].trim();
            startChat(username);
        });
    });

    document.querySelector("#send-message").addEventListener("click", () => {
        sendMessage();
    });

    document.querySelector("#close-chat").addEventListener('click', () => {
        closeChat();
    });
});

let ws;
let currentUser;

function openWs() {
    ws = new WebSocket("ws://localhost:8989/messages");

    ws.onopen = () => {
        console.log("Connected to WebSocket...");
    };

    ws.onmessage = (e) => {
        const message = JSON.parse(e.data);
        console.log("message", message);
        if (message.type === 'history') {
            loadHistory(message.messages);
        } else {
            displayMessage(message.mess, message.date,"sms-received");
        }
    };

    ws.onerror = (error) => {
        console.log("WebSocket Error: ", error);
    };

    ws.onclose = () => {
        console.log("WebSocket connection closed");
    };

   console.log(getCookie("sessionID"))
}

function getCookie(name) {
    const value = '; ' + document.cookie;
    const parts = value.split('; ' + name + '=');
    console.log(parts);
    if (parts.length === 2) return parts.pop().split(';').shift();
  }

function startChat(username) {
    Chatopened = true
    currentUser = username;
    document.getElementById('chat-user-name').textContent = username;
    document.getElementById('chat-container').style.display = 'block';
    document.getElementById('chat-messages').innerHTML = ''; 
    sendInitialMessage(username);
}

function sendInitialMessage(username) {
    const data = {
        mess: "",
        uname: username,
        date: "",
        type: "fetch"
    };
    ws.send(JSON.stringify(data));
}

function sendMessage() {
    const messageInput = document.getElementById('chat-input');
    const message = messageInput.value;
    messageInput.value = "";
    if (message != "") {
        const date = new Date();
        const dateTime = `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()} ${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
        const data = {
            mess: message,
            uname: currentUser,
            date: dateTime,
            type: "message"
        };

        console.log("je suis here",data.uname);
        ws.send(JSON.stringify(data));
        displayMessage(message, dateTime,"sms");
    }
}

function displayMessage(message, date,type) {
        const chatMessages = document.getElementById("chat-messages");
        const messageDiv = document.createElement("div");
        const timeDiv = document.createElement("div");

        timeDiv.classList.add("dayeu");
        timeDiv.textContent = date;

        messageDiv.classList.add(type);
        messageDiv.textContent = message;

        chatMessages.appendChild(timeDiv);
        chatMessages.appendChild(messageDiv);
}


function loadHistory(messages) {
    const chatMessages = document.getElementById("chat-messages");
    chatMessages.innerHTML = '';
    messages.forEach(msg => {
        displayMessage(msg.mess, msg.date,"sms-received");
    });
}

function closeChat() {
    document.getElementById('chat-container').style.display = 'none';
    document.getElementById("chat-messages").innerHTML = '';
}
