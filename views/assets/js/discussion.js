const currentUser = document.querySelector('#user-full-name').textContent;
const messageInput = document.getElementById('chat-input');
let currentdiscussion;

let ws;
let currentDestinataire;

document.addEventListener('DOMContentLoaded', () => {
    loadUsers()
    openWs();

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

        discussionsArray.sort((a, b) => {
            const dateA = a.getAttribute('date');
            const dateB = b.getAttribute('date');
            return compareDiscussionOrder(dateA, dateB);
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

        if (message.type === 'history') {
            var temp = message.messages.reverse()
            loadHistory(temp);
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
            setDivDate(message.date, message.sender)

            if (currentdiscussion == message.sender) {
                displayMessage(message.message, message.date, message.destinataire, message.sender, false);
            }
            // console.log(message);
            var div = document.getElementById(message.sender + "_discussion")
            var ischatopened = div.getAttribute("ischatopened")
            // console.log(ischatopened);
            if (ischatopened === "false" && message.message != "") {
                updateNotification(message.sender, message.Notif.notification)

                var discussions = document.querySelectorAll('.discussion');
                var discussionsArray = Array.from(discussions);

                // console.log("discussionsArray",discussionsArray);

                discussionsArray.sort((a, b) => {
                    const dateA = a.getAttribute('date');
                    const dateB = b.getAttribute('date');
                    return compareDiscussionOrder(dateA, dateB);
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

function parseDate(dateStr) {
    // console.log("date", dateStr);
    if (dateStr) {
        const [datePart, timePart] = dateStr.split(' ');

        const [day, month, year] = datePart.split('/').map(Number);
        const [hours, minutes, seconds] = timePart.split(':').map(Number);
        return new Date(year, month - 1, day, hours, minutes, seconds);
    }
}

function compareDiscussionOrder(dateStr1, dateStr2) {
    const date1 = parseDate(dateStr1);
    const date2 = parseDate(dateStr2);

    return date2 - date1;
}

function updateNotification(sender, count) {
    if (sender.includes(" ")) {
        sender = sender.replace(" ", "_")
    }

    let notificationElement = document.getElementById(sender + "_notif");
    if (notificationElement) {
        if (count !== 0) {
            notificationElement.textContent = parseInt(count);
        } else {
            notificationElement.textContent = "";
        }
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

        ws.send(JSON.stringify(data));
        displayMessage(message, dateTime, currentDestinataire, currentUser, false);
    }
}

function displayMessage(message, date, destinataire, sender, prepend) {
    if (message != "") {
        const chatMessages = document.getElementById("chat-messages");
        const timeDiv = document.createElement("div");
        const smsContainer = document.createElement("div");
        const smsDiv = document.createElement("div");
        const smsUsername = document.createElement("div");

        smsContainer.classList.add("smsContainer")
        timeDiv.classList.add("smsDate");
        
        timeDiv.textContent = date;
        if (currentUser !== destinataire) {
            smsUsername.classList.add("smsUsername-sms")
            console.log(sender, message);
            smsUsername.textContent = sender
            smsUsername.style.justifyContent = "flex-end"
            // smsUsername.style.marginRight = "auto"
            smsDiv.classList.add("sms");
        } else {
            smsUsername.classList.add("smsUsername-sms-received")
            smsUsername.textContent = sender
            smsDiv.classList.add("sms-received");
        }

        smsDiv.textContent = message;
        smsContainer.appendChild(smsUsername)
        smsContainer.appendChild(smsDiv)
        if (prepend) {
            chatMessages.insertBefore(smsContainer, chatMessages.firstChild);
            chatMessages.insertBefore(timeDiv, chatMessages.firstChild);
        } else {
            chatMessages.appendChild(timeDiv);
            chatMessages.appendChild(smsContainer);
        }

    }
}

let messageChunks = [];
let currentChunkIndex = 0;

function chunks(messages) {
    messageChunks = [];
    // messageChunks.slice().reverse()
    for (let i = 0; i < messages.length; i += 10) {
        messageChunks.push(messages.slice(i, i + 10));
    }
    currentChunkIndex = 0;
}

function loadChunk() {
    if (currentChunkIndex < messageChunks.length) {
        const chunk = messageChunks[currentChunkIndex];



        chunk.forEach(msg => {
            updateNotification(msg.sender, 0);
            displayMessage(msg.message, msg.date, msg.destinataire, msg.sender, true);
        });



        currentChunkIndex++;
    }
}

let isChunkLoading = false;
function handleScroll() {
    // console.log(isChunkLoading);
    const chatMessages = document.getElementById("chat-messages");
    if (chatMessages.scrollTop === 0 && currentChunkIndex < messageChunks.length && !isChunkLoading) {
        isChunkLoading = true;

        const loadingIndicator = document.createElement("div");
        loadingIndicator.classList.add("loading-indicator");
        loadingIndicator.textContent = "Chargement en cours...";
        chatMessages.prepend(loadingIndicator);

        setTimeout(() => {
            chatMessages.removeChild(loadingIndicator);
            loadChunk();
            isChunkLoading = false;
        }, 500);
    }
}

function debounce(func, wait) {
    let timer;
    return (...args) => {
        clearTimeout(timer);
        timer = setTimeout(() => func(...args), wait);
    };
}


function loadHistory(messages) {
    const chatMessages = document.getElementById("chat-messages");

    chatMessages.innerHTML = '';
    chunks(messages)
    loadChunk();

    chatMessages.addEventListener('scroll', debounce(handleScroll, 1000));
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
        loadUsers()
    }
}

function setDivDate(date, sender) {
    document.querySelectorAll('.discussion').forEach(user => {
        if (user.getAttribute('id').split("_discussion")[0] === sender) {
            user.setAttribute("date", date)
        }

    });
}

let firstCall = true;

function loadUsers() {
    fetch('/api/AllUsers')
        .then(response => response.json())
        .then(users => {
            const usersContainer = document.getElementById('most-active-users');
            usersContainer.innerHTML = '';

            users.forEach((user, index) => {
                const userDiv = document.createElement('div');
                userDiv.className = 'discussion';
                userDiv.id = `${user.Username}_discussion`;
                userDiv.setAttribute('ischatopened', 'false');
                userDiv.setAttribute('date', user.LastMsgDate);
                // console.log("this is user", user.NbrNotifications, typeof (user.NbrNotifications));
                const userContainer = document.createElement('div');
                userContainer.className = 'usercontainer';

                const userStatus = document.createElement(user.Status === 'Online' ? 'strong' : 'em')

                if (!firstCall && index === users.length - 1) {
                    // console.log('dernier user')
                    userStatus.textContent = `${user.Username} - Online`;
                    // userStatus.tagName = 'strong';
                } else {
                    // console.log('C est pas dernier user')
                    userStatus.textContent = `${user.Username} - ${user.Status !== 'Online' ? 'Offline' : user.Status}`;
                }
                const notificationDiv = document.createElement('div');
                notificationDiv.className = 'notification';
                notificationDiv.id = `${user.Username}_notif`;
                if (user.NbrNotifications == 0) {
                    notificationDiv.textContent = "";
                } else {
                    notificationDiv.textContent = user.NbrNotifications;
                }

                userContainer.appendChild(userStatus);
                userContainer.appendChild(notificationDiv);
                userDiv.appendChild(userContainer);
                usersContainer.appendChild(userDiv);
            });

            if (firstCall) {
                firstCall = false;
            }

            document.querySelectorAll('.discussion').forEach(user => {
                const username = user.textContent.split(' - ')[0].trim();
                user.addEventListener('click', () => {
                    user.setAttribute("ischatopened", true);
                    currentdiscussion = username;
                    startChat(username);
                });
            });
        })
        .catch(error => console.error('Error loading users:', error));
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
