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
    };



    ws.onmessage = (e) => {
        // console.log("Received:", e.data);

        const data = e.data.split('],[')
        var mess
        for (let i = data.length-1; i >= 0; i--) {
            const element = data[i];

            var time = element.replace("[", "").replace("]", "").slice(-19)
            if (i == data.length - 1) {
                mess = element.replace("]", "").slice(0, element.length - 20)

            } else if (i == 0) {
                mess = element.replace("[", "").slice(0, element.length - 20)
            } else {
                mess = element.slice(0, element.length - 20)
            }
            console.log("time", time);
            console.log("mess", mess);

            if (mess != "") {
                setupSendMessage(ws);

                displayMessage(mess, time);
            }
            
        }
        data.forEach((e, i) => {
            
        });

        // const Jsonstr = e.data.replace("\\","")
        // console.log("this is js",Jsonstr.replace("\\",""));
        // console.log(JSON.parse(e.data));
        var res = e.data.split("],");
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
        if (message != "") {
            const data = {
                mess: message,
                // uname: document.getElementById('chat-user-name').textContent,
                // date: dateTime,
            };
            var date = new Date()
            const dateTime = `${date.getDate()}/${date.getMonth() + 1}/${date.getFullYear()} ${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`;
            ws.send(JSON.stringify(data))
            displayMessage(message, dateTime);
        }


    });
}

function displayMessage(mess, time) {
    // console.log("data", data);
    if (mess != "") {
        const endroit = document.getElementById("chat-messages");
        // for (let i = 0; i < data.length; i++) {
        // const element = data[i];
        // console.log("element0", data[0]);
        // console.log("element1", data[1]);


        const mns = document.createElement("div");
        const dayeu = document.createElement("div");

        // const dateTime = `${data.date.getDate()}/${data.date.getMonth() + 1}/${data.date.getFullYear()} ${data.date.getHours()}:${data.date.getMinutes()}:${data.date.getSeconds()}`;
        // const dateTime = "time"
        dayeu.classList.add("dayeu");
        dayeu.textContent = time;

        mns.classList.add("sms");
        mns.textContent = mess;

        endroit.appendChild(dayeu);
        endroit.appendChild(mns);
    }

}

document.querySelector("#close-chat").addEventListener('click', () => {
    document.getElementById('chat-container').style.display = 'none';
    document.getElementById("chat-messages").innerHTML = ''
});
