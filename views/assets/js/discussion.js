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
        if (message.type === 'history') {
            loadHistory(message.messages);
        } else {
            displayMessage(message.mess, message.date);
        }
    };

    ws.onerror = (error) => {
        console.log("WebSocket Error: ", error);
    };

    ws.onclose = () => {
        console.log("WebSocket connection closed");
    };
}

function startChat(username) {
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
        ws.send(JSON.stringify(data));
        displayMessage(message, dateTime);
    }
}

function displayMessage(message, date) {
    if (message != "") {
        const chatMessages = document.getElementById("chat-messages");
        const messageDiv = document.createElement("div");
        const timeDiv = document.createElement("div");

        timeDiv.classList.add("dayeu");
        timeDiv.textContent = date;

        messageDiv.classList.add("sms");
        messageDiv.textContent = message;

        chatMessages.appendChild(timeDiv);
        chatMessages.appendChild(messageDiv);
    }
}


function loadHistory(messages) {
    const chatMessages = document.getElementById("chat-messages");
    chatMessages.innerHTML = '';
    messages.forEach(msg => {
        displayMessage(msg.mess, msg.date);
    });
}

function closeChat() {
    document.getElementById('chat-container').style.display = 'none';
    document.getElementById("chat-messages").innerHTML = '';
}
