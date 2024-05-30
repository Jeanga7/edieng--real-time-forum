/* document.querySelector(".discussion").addEventListener("click", function () {
    var name = document.querySelector(".discussion").textContent.split("-")[0];
    console.log(name);

    // var popup = document.createElement("div");
    // popup.classList.add("layer");
    // document.body.appendChild(popup);
    // popup.innerHTML = `
    //     <div class="popup">
    //         <span id="close_popup">X</span>
    //     </div>
    // `;
    
    let layer = document.querySelector(".layer");
    layer.style.display = "flex";
});

let closePopup = document.querySelector("#close_popup");
if (closePopup) {
    closePopup.addEventListener("click", () => {
        var layer = document.querySelector(".layer");
        layer.style.display = "none";
    });
} */

document.querySelectorAll('.discussion').forEach(user => {
    user.addEventListener('click', () => {
        const username = user.textContent.split(' - ')[0].trim();
        startChat(username);
    });
});

function startChat(username) {
    document.getElementById('chat-user-name').textContent = username;
    document.getElementById('chat-container').style.display = 'block';

    // fetchMessages(username);
}

// function fetchMessages(username) {
//     fetch(`/messages?username=${username}`)
//         .then(response => response.json())
//         .then(data => {
//             const messagesContainer = document.getElementById('chat-messages');
//             messagesContainer.innerHTML = '';
//             data.messages.forEach(msg => {
//                 const messageElement = document.createElement('div');
//                 messageElement.textContent = msg.content;
//                 messagesContainer.appendChild(messageElement);
//             });
//         });
// }

function sendMessage() {
    const ws = new WebSocket("ws://localhost:8989/messages")

    ws.onopen = () => {
        console.log("I'm Connected...");
        const messageInput = document.getElementById('chat-input');
        const message = messageInput.value;
        console.log(message);
        messageInput.value = ""
        const username = document.getElementById('chat-user-name').textContent;
        console.log(username);
        var date = new Date()

        console.log(date);
        var data = {
            mess: message,
            uname: username,
            date: date,
        }
    
        var jsonstring = JSON.stringify(data)
        const blob = new Blob([jsonstring], { type: 'application/json' });
    
        ws.send(blob)

        var endroit = document.getElementById("chat-messages")
        var mns = document.createElement("div")
        // var container = document.createElement("div")
        var dayeu = document.createElement("div")
        var dateTime = data.date.getDay()+"/"+data.date.getMonth()+"/"+data.date.getFullYear()+" "+data.date.getHours()+':'+data.date.getMinutes()+':'+data.date.getSeconds()
        dayeu.classList.add("dayeu")
        dayeu.textContent = dateTime
        mns.classList.add("sms")
        mns.textContent = data.mess

        // container.appendChild(dayeu)
        // container.appendChild(mns)


        endroit.appendChild(dayeu)
        endroit.appendChild(mns)
        // fetch('/send-message', {
        //     method: 'POST',
        //     headers: {
        //         'Content-Type': 'application/json'
        //     },
        //     body: JSON.stringify({ username, message })
        // }).then(() => {
        //     messageInput.value = '';
        //     const messageElement = document.createElement('div');
        //     messageElement.textContent = message;
        //     document.getElementById('chat-messages').appendChild(messageElement);
        // });
    }
    // event.preventDefault();

}

document.querySelector("#send-message").addEventListener("click",()=>{
    sendMessage()
})

document.querySelector("#close-chat").addEventListener('click', () => {
    document.getElementById('chat-container').style.display = 'none';

})