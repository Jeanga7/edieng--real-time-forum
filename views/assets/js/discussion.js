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
        const username = user.textContent.split(' - ')[0];
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

// function sendMessage(event) {
//     event.preventDefault();
//     const messageInput = document.getElementById('chat-input');
//     const message = messageInput.value;
//     const username = document.getElementById('chat-user-name').textContent;

//     fetch('/send-message', {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json'
//         },
//         body: JSON.stringify({ username, message })
//     }).then(() => {
//         messageInput.value = '';
//         const messageElement = document.createElement('div');
//         messageElement.textContent = message;
//         document.getElementById('chat-messages').appendChild(messageElement);
//     });
// }


document.querySelector("#close-chat").addEventListener('click', () => {
    document.getElementById('chat-container').style.display = 'none';

})