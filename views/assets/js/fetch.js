async function fetchUserData(userId) {
    return new Promise((resolve) => {
        let fd = new FormData();
        fd.append("id", userId);
        fetch(userDataRoute, {
                method: 'POST',
                body: fd
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            })
    });
}
async function fetchPosts(userId) {
    return new Promise((resolve) => {
        let fd = new FormData();
        fd.append("id", userId);
        fetch(postsRoute, {
                method: 'GET'
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            })
    });
}
async function fetchAllUsers() {
    return new Promise((resolve) => {
        let data = new FormData();
        data.append("sessionID", sessionID);
        fetch(allUsersRoute, {
                method: 'POST',
                body: data
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            });
    });
}
async function fetchChat(chatUser) {
    return new Promise((resolve) => {
        let fd = new FormData();
        fd.append("type", "fetch");
        fd.append("user", userNickname);
        fd.append("with", chatUser);
        fetch(discussionsRoute, {
                method: 'POST',
                body: fd
            })
            .then(r => r.json())
            .then(d => {
                resolve(d);
            });
    });
}