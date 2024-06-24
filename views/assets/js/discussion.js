const currentUser = document.querySelector('#user-full-name').textContent;
const messageInput = document.getElementById('chat-input');
let currentdiscussion;
let ws;
let currentDestinataire;
let messageChunks = [];
let currentChunkIndex = 0;
let isChunkLoading = false;
let firstCall = true;
document.addEventListener('DOMContentLoaded', () => {
    loadUsers()
    openWs();
    document.querySelector("#send-message").addEventListener("click", sendMessage);
    document.querySelector("#close-chat").addEventListener('click', () => {
        document.querySelectorAll('.discussion').forEach(user => {
            user.setAttribute("Ischatopened", false)
        });
        closeChat();
    });
    messageInput.addEventListener('keydown', () => {
        const data3 = {
            sender: currentUser,
            destinataire: currentdiscussion,
            type: "typing"
        };
        ws.send(JSON.stringify(data3));
    });
});
function openWs() {
    ws = new WebSocket("ws://localhost:8989/messages");
    ws.onopen = () => {
        console.log("Connected to WebSocket...");
        sortDiscussion();
        AskforNotification()
    };
    ws.onmessage = (e) => {
        const message = JSON.parse(e.data);
        if (message.type === 'history') {
            loadHistory(message.messages.reverse());
        } else if (message.type == "typing") {
            if (currentdiscussion === message.sender) {
                showTyping();
            }
        } else if (message.type == "new_user") {
            handleNewUser(message.sender)
        } else {
            handleMessage(message);
        }
    }
}
function handleMessage(message) {
    console.log("NOOOOOOOOOOOOOOON");
    setDivDate(message.date, message.sender);
    if (currentdiscussion == message.sender) {
        displayMessage(message.message, message.date, message.destinataire, message.sender, false);
    }
    var div = document.getElementById(message.sender + "_discussion");
    var ischatopened = div.getAttribute("ischatopened");
    if (ischatopened === "false" && message.message != "") {
        updateNotification(message.sender, message.Notif.notification);
        sortDiscussion();
    }
}
function showTyping() {
    document.getElementById('typingIndicator').style.opacity = 1;
    setTimeout(() => {
        document.getElementById('typingIndicator').style.opacity = 0;
    }, 1000);
}
function sortDiscussion() {
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
}
function parseDate(dateStr) {
    if (dateStr) {
        const [datePart, timePart] = dateStr.split(' ');
        const [day, month, year] = datePart.split('/').map(Number);
        const [hours, minutes, seconds] = timePart.split(':').map(Number);
        return new Date(year, month - 1, day, hours, minutes, seconds);
    }
}
function compareDiscussionOrder(dateStr1, dateStr2) {
    return parseDate(dateStr2) - parseDate(dateStr1)
}
function updateNotification(sender, count) {
    let notificationElement = document.getElementById(sender.replace(" ", "_") + "_notif");
    if (notificationElement) {
        notificationElement.textContent = count !== 0 ? parseInt(count) : "";
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
        scrollToLastMessage();
    }
}
function displayMessage(message, date, destinataire, sender, prepend) {
    if (!message) return
    const chatMessages = document.getElementById("chat-messages");
    const timeDiv = document.createElement("div");
    const smsContainer = document.createElement("div");
    const smsDiv = document.createElement("div");
    const smsUsername = document.createElement("div");
    smsContainer.classList.add("smsContainer")
    timeDiv.classList.add("smsDate");
    timeDiv.textContent = date;
    smsUsername.textContent = sender
    if (currentUser !== destinataire) {
        smsUsername.classList.add("smsUsername-sms")
        smsDiv.classList.add("sms");
    } else {
        smsUsername.classList.add("smsUsername-sms-received")
        smsDiv.classList.add("sms-received");
    }
    smsDiv.textContent = message;
    smsContainer.appendChild(smsUsername);
    smsContainer.appendChild(smsDiv);
    if (prepend) {
        chatMessages.prepend(timeDiv, smsContainer);
    } else {
        chatMessages.append(timeDiv, smsContainer);
        scrollToLastMessage();
    }
}
function chunks(messages) {
    messageChunks = [];
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
function handleScroll() {
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
                const userContainer = document.createElement('div');
                userContainer.className = 'usercontainer';
                const userStatus = document.createElement(user.Status === 'Online' ? 'strong' : 'em')
                if (!firstCall && index === users.length - 1) {
                    userStatus.textContent = `${user.Username} - Online`;
                } else {
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
function scrollToLastMessage() {
    const chatMessages = document.getElementById("chat-messages");
    chatMessages.scrollTop = chatMessages.scrollHeight;
}