const currentUser = document.querySelector('#user-full-name').textContent;

let ws;
let currentDestinataire;

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


function openWs() {
    ws = new WebSocket("ws://localhost:8989/messages");

    ws.onopen = () => {
        console.log("Connected to WebSocket...");
    };

    ws.onmessage = (e) => {
        const message = JSON.parse(e.data);
        // console.log("destinataire", currentDestinataire);
        if (message.type === 'history') {
            console.log("destinataire",message);
            loadHistory(message.messages);
        } else {
            console.log("desti",message.destinataire);
            displayMessage(message.message, message.date, message.destinataire);
        }
    };

}

function startChat(username) {
    currentDestinataire = username;
    document.getElementById('chat-user-name').textContent = username;
    document.getElementById('chat-container').style.display = 'block';
    document.getElementById('chat-messages').innerHTML = '';
    sendInitialMessage(username);
}

function sendInitialMessage(username) {
    const data = {
        message: "",
        sender: currentUser,
        destinataire: username,
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
            message: message,
            sender: currentUser,
            destinataire: currentDestinataire,
            date: dateTime,
            type: "message"
        };

        // console.log("je suis here", data.destinataire);
        ws.send(JSON.stringify(data));
        displayMessage(message, dateTime, currentDestinataire);
    }
}

function displayMessage(message, date, destinataire) {
    // class sms-received 
    // console.log("messageReçu", message);
    if (message != "") {
        const chatMessages = document.getElementById("chat-messages");
        const messageDiv = document.createElement("div");
        const timeDiv = document.createElement("div");

        timeDiv.classList.add("smsDate");
        timeDiv.textContent = date;
        console.log(";",currentUser);
        console.log("ù",destinataire);
        if (currentUser !== destinataire) {
            messageDiv.classList.add("sms");
        } else {
            messageDiv.classList.add("sms-received");
        }
        messageDiv.textContent = message;

        chatMessages.appendChild(timeDiv);
        chatMessages.appendChild(messageDiv);
    }
}


function loadHistory(messages) {
    const chatMessages = document.getElementById("chat-messages");
    chatMessages.innerHTML = '';
    messages.forEach(msg => {
        console.log("msg", msg);
        displayMessage(msg.message, msg.date, msg.destinataire);
    });
}

function closeChat() {
    document.getElementById('chat-container').style.display = 'none';
    document.getElementById("chat-messages").innerHTML = '';
}
