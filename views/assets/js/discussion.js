document.querySelectorAll('.discussion').forEach(user => {
    user.addEventListener('click', () => {
        const username = user.textContent.split(' - ')[0].trim();
        startChat(username);
    });
});

function startChat(username) {
    document.getElementById('chat-user-name').textContent = username;
    document.getElementById('chat-container').style.display = 'block';
    openWs()
    // fetchMessages(username);
}

function openWs() {
    const ws = new WebSocket("ws://localhost:8989/messages")

    ws.onopen = () => {
        console.log("I'm Connected...");
        sendMessage(ws);
    }
    ws.onmessage = (e)=>{
        console.log("received",e.data);
    }
}

// document.querySelector("#send-message").addEventListener("click",()=>{
//     sendMessage()
// })

document.querySelector("#close-chat").addEventListener('click', () => {
    document.getElementById('chat-container').style.display = 'none';

})

function sendMessage(ws) {
    const username = document.getElementById('chat-user-name').textContent;

    ws.send(username)
    document.querySelector("#send-message").addEventListener("click", () => {


        const messageInput = document.getElementById('chat-input');
        const message = messageInput.value;
        console.log(message);
        messageInput.value = "";
        // console.log(username);
        var date = new Date();

        console.log(date);
        var data = {
            mess: message,
            uname: username,
            date: date,
        };

        // var jsonstring = JSON.stringify(data)
        // const blob = new Blob([jsonstring], { type: 'application/json' });
        // ws.send(blob)
        ws.send(JSON.stringify(data));


        var endroit = document.getElementById("chat-messages");
        var mns = document.createElement("div");
        // var container = document.createElement("div")
        var dayeu = document.createElement("div");
        var dateTime = `${data.date.getDate()}/${data.date.getMonth() + 1}/${data.date.getFullYear()}   ${data.date.getHours()}:${data.date.getMinutes()}:${data.date.getSeconds()}`;
        dayeu.classList.add("dayeu");
        dayeu.textContent = dateTime;
        mns.classList.add("sms");
        mns.textContent = data.mess;

        // container.appendChild(dayeu)
        // container.appendChild(mns)
        endroit.appendChild(dayeu);
        endroit.appendChild(mns);

    });
}
