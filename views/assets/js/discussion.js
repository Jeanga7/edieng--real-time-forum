document.querySelectorAll('.discussion').forEach(user => {
    user.addEventListener('click', () => {
        const username = user.textContent.split(' - ')[0].trim();
        startChat(username);
    });
});

function startChat(username) {
    document.getElementById('chat-user-name').textContent = username;
    document.getElementById('chat-container').style.display = 'block';
    openWs(username);
}

function openWs(username) {
    const ws = new WebSocket("ws://localhost:8989/messages");

    ws.onopen = () => {
        console.log("Connected to WebSocket...");
        sendInitialMessage(ws, username);
        setupSendMessage(ws);
    };

    ws.onmessage = (e) => {
        console.log("Received:", e.data);
        // const Jsonstr = e.data.replace("\\","")
        // console.log("this is js",Jsonstr.replace("\\",""));
        console.log(JSON.parse(e.data));
        displayMessage(JSON.parse(e.data), "incoming");
    };
}

function sendInitialMessage(ws, username) {
    ws.send(username);
}

function setupSendMessage(ws) {
    document.querySelector("#send-message").addEventListener("click", () => {
        const messageInput = document.getElementById('chat-input');
        const message = messageInput.value;
        messageInput.value = "";

        const data = {
            mess: message,
            uname: document.getElementById('chat-user-name').textContent,
            date: new Date(),
        };

        ws.send(JSON.stringify(data));
        displayMessage(data, "outgoing");
    });
}

function displayMessage(data, type) {
    console.log("data",data);
    const endroit = document.getElementById("chat-messages");
    for (let i = 0; i < data.length; i++) {
        const element = data[i];
        console.log("element",element);

        const mns = document.createElement("div");
        const dayeu = document.createElement("div");

        // const dateTime = `${data.date.getDate()}/${data.date.getMonth() + 1}/${data.date.getFullYear()} ${data.date.getHours()}:${data.date.getMinutes()}:${data.date.getSeconds()}`;
        const dateTime = "time"
        dayeu.classList.add("dayeu");
        dayeu.textContent = dateTime;

        mns.classList.add("sms", type);
        mns.textContent = element.mess;

        endroit.appendChild(dayeu);
        endroit.appendChild(mns);
    }
}

document.querySelector("#close-chat").addEventListener('click', () => {
    document.getElementById('chat-container').style.display = 'none';
});
