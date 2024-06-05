const currentUser = document.querySelector('#user-full-name').textContent;

// let Ischatopened = false
let currentdiscussion;

let ws;
let currentDestinataire;

document.addEventListener('DOMContentLoaded', (event) => {
    openWs();
    
    document.querySelectorAll('.discussion').forEach(user => {
        const username = user.textContent.split(' - ')[0].trim();
        user.addEventListener('click', () => {
            user.setAttribute("Ischatopened", true)
            currentdiscussion = username
            startChat(username);
        });
    });

    document.querySelector("#send-message").addEventListener("click", () => {
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
        AskforNotification()
    };



    ws.onmessage = (e) => {
        const message = JSON.parse(e.data);
        if (message.type === 'history') {
            loadHistory(message.messages);
        } else {
            console.log("NOOOOOOOOOOOOOOON");
            
            if (currentdiscussion == message.sender) {
                displayMessage(message.message, message.date, message.destinataire);
            }
            console.log(message);
            var div = document.getElementById(message.sender+"_discussion")
            var t = div.getAttribute("ischatopened")
            console.log(t);
            if (t === "false"){
                updateNotification(message.sender, message.Notif.notification)
            }
        }
    };

}

function updateNotification(sender, count) {
    if (sender.includes(" ")) {
        sender = sender.replace(" ", "_")
    }

    let notificationElement = document.getElementById(sender+"_notif");
    if (notificationElement) {
        notificationElement.textContent = parseInt(count);
    }
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

    ///
   
    ws.send(JSON.stringify(data));
}

function AskforNotification(){
    const data2 = {
        sender: currentUser,
        type: "notification"
    };
    ws.send(JSON.stringify(data2));
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
        updateNotification(msg.sender, "0")
        displayMessage(msg.message, msg.date, msg.destinataire);
    });
}

function closeChat() {
    document.getElementById('chat-container').style.display = 'none';
    // Ischatopened = false
    document.getElementById("chat-messages").innerHTML = '';
}
