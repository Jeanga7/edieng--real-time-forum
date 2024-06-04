const currentUser = document.querySelector('#user-full-name').textContent;

// let Ischatopened = false
let currentdiscussion

let ws;
let currentDestinataire;

document.addEventListener('DOMContentLoaded', (event) => {
    openWs();
    document.querySelectorAll('.discussion').forEach(user => {
        // console.log(user, Ischatopened);
        const username = user.textContent.split(' - ')[0].trim();
        user.addEventListener('click', () => {
            user.setAttribute("Ischatopened", true)
            // Ischatopened = true
            // console.log(username);
            currentdiscussion = username
            startChat(username);
        });
    });

    document.querySelector("#send-message").addEventListener("click", () => {
        // Ischatopened = true
        sendMessage();
    });

    document.querySelector("#close-chat").addEventListener('click', () => {
        // var divs = document.getElementById(currentdiscussion)
        // divs.style.backgroundColor = "green"
        document.querySelectorAll('.discussion').forEach(user => {
            user.setAttribute("Ischatopened", false)
        });
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
        var sender
        if (message.type === 'history') {
            loadHistory(message.messages);
        } else {
            if (message.sender.includes(" ")) {
                sender = message.sender.replace(" ", "_")
            } else {
                sender = message.sender
            }

            // const document.querySelector
            // if (Ischatopened && message.destinataire === currentUser) {
            console.log("NOOOOOOOOOOOOOOON");
            displayMessage(message.message, message.date, message.destinataire);
            var divs = document.getElementById(sender)
            divs.style.backgroundColor = "red"
            console.log(divs);
            // } else {
            // }
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
    if (message != "") {
        const chatMessages = document.getElementById("chat-messages");
        const messageDiv = document.createElement("div");
        const timeDiv = document.createElement("div");

        timeDiv.classList.add("smsDate");
        timeDiv.textContent = date;
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
        displayMessage(msg.message, msg.date, msg.destinataire);
    });
}

function closeChat() {
    document.getElementById('chat-container').style.display = 'none';
    // Ischatopened = false
    document.getElementById("chat-messages").innerHTML = '';
}
