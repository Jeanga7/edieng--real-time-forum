const currentUser = document.querySelector('#user-full-name').textContent;
const messageInput = document.getElementById('chat-input');
// let Ischatopened = false
let currentdiscussion;

let ws;
let currentDestinataire;

document.addEventListener('DOMContentLoaded', () => {
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

        var discussions = document.querySelectorAll('.discussion');

        var discussionsArray = Array.from(discussions);

        discussionsArray.sort(function (a, b) {
            return compareDiscussionOrder(a, b);
        });

        var parentContainer = document.getElementById('most-active-users');
        parentContainer.innerHTML = '';

        discussionsArray.forEach(function (discussion) {
            parentContainer.appendChild(discussion);
        });
        AskforNotification()
    };



    ws.onmessage = (e) => {
        const message = JSON.parse(e.data);
        // console.log("message", message.sender);
        setDivDate(message.date, message.sender)

        if (message.type === 'history') {
            loadHistory(message.messages);
        } else if (message.type == "typing") {
            if (currentdiscussion === message.sender) {
                document.getElementById('typingIndicator').style.opacity = 1;
                setTimeout(() => {
                    document.getElementById('typingIndicator').style.opacity = 0;
                }, 1000);
            }

        } else if (message.type == "new_user") {
            handleNewUser(message.sender)
        } else {
            console.log("NOOOOOOOOOOOOOOON");
            if (currentdiscussion == message.sender) {
                displayMessage(message.message, message.date, message.destinataire);
            }
            // console.log(message);
            var div = document.getElementById(message.sender + "_discussion")
            var t = div.getAttribute("ischatopened")
            console.log(t);
            if (t === "false" && message.message != "") {
                updateNotification(message.sender, message.Notif.notification)


                var discussions = document.querySelectorAll('.discussion');
                var discussionsArray = Array.from(discussions);

                // console.log("discussionsArray",discussionsArray);

                discussionsArray.sort(function (a, b) {
                    // console.log("a", a);
                    return compareDiscussionOrder(a, b);
                });

                var parentContainer = document.getElementById('most-active-users');
                parentContainer.innerHTML = '';

                discussionsArray.forEach(function (discussion) {
                    parentContainer.appendChild(discussion);
                });

            }
        }
    };

}

// function compareDiscussionOrder(a, b) {
//     let ischatopened = a.getAttribute('id')
//     let v = b.getAttribute('id')
//     var notifA = parseInt(document.getElementById(ischatopened.replace("_discussion", "_notif")).textContent);
//     var notifB = parseInt(document.getElementById(v.replace("_discussion", "_notif")).textContent);
//     // console.log("*", ischatopened.replace("_discussion", "_notif"));
//     // console.log(notifA, notifB);
//     return notifB - notifA;
// }

function updateNotification(sender, count) {
    if (sender.includes(" ")) {
        sender = sender.replace(" ", "_")
    }

    let notificationElement = document.getElementById(sender + "_notif");
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
    ws.send(JSON.stringify(data));
}

function AskforNotification() {
    const data2 = {
        sender: currentUser,
        type: "notification"
    };
    ws.send(JSON.stringify(data2));
}

function sendMessage() {
    const message = messageInput.value.trim();
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
        updateNotification(msg.sender, 0)
        displayMessage(msg.message, msg.date, msg.destinataire);
    });
}

function closeChat() {
    document.getElementById('chat-container').style.display = 'none';
    // Ischatopened = false
    document.getElementById("chat-messages").innerHTML = '';
}

function handleNewUser(userConnected) {
    const userElement = document.querySelector(`#${userConnected}_discussion`);

    if (userElement) {
        const parts = userElement.textContent.split('-');
        if (parts[1].trim().includes("Offline")) {
            userElement.innerHTML = `<strong>${parts[0]} - ${parts[1].replace("Offline", "Online")}</strong>`;
        } else {
            userElement.innerHTML = `<em>${parts[0]} - ${parts[1].replace("Online", "Offline")}</em>`;
        }
    } else {
        // Ne marche pas cas ou l'utilisateur vient de se connecter

        // const newUserElement = `
        //     <div class="discussion" id="${userConnected}_discussion" ischatopened="false">
        //       <div class="usercontainer">
        //         <em>${userConnected} - Offline</em>
        //         <div class="notification" id="${userConnected}_notif">0</div>
        //       </div>
        //     </div>
        // `;
        // document.querySelector('#most-active-users').insertAdjacentHTML('beforeend', newUserElement);
    }
}

function setDivDate(date, sender) {
    document.querySelectorAll('.discussion').forEach(user => {
        console.log(user);
        // user.setAttribute("date", date)
    });
}

//====================================TYPING-PROGRESS========================================

messageInput.addEventListener('keydown', () => {
    const data3 = {
        sender: currentUser,
        destinataire: currentdiscussion,
        type: "typing"
    };
    ws.send(JSON.stringify(data3));
});

// order by the number of notifications
function compareDiscussionOrder(a, b) {
    let t = a.getAttribute('id')
    let v = b.getAttribute('id')
    var notifA = parseInt(document.getElementById(t.replace("_discussion", "_notif")).textContent);
    var notifB = parseInt(document.getElementById(v.replace("_discussion", "_notif")).textContent);
    console.log("*", t.replace("_discussion", "_notif"));
    console.log(notifA, notifB);
    return notifB - notifA;
}