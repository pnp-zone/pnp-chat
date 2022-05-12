function staticUrl(url) {
    return `${window.pnpZone.static}plugins/pnp-chat/${url}`;
}

const JOIN = "join";
const MESSAGE = "message";
const LEAVE = "leave";

class ChatConnection {
    constructor(username) {
        const {host, protocol} = window.location;
        this.socket = new WebSocket(`${protocol.replace("http", "ws")}//${host}/chat/socket`);
        this.username = username;
    }

    join() {
        this.socket.send(JSON.stringify({
            type: JOIN,
            user: this.username,
            msg: null,
        }));
    }

    send(msg) {
        this.socket.send(JSON.stringify({
            type: MESSAGE,
            user: this.username,
            msg,
        }));
    }

    leave() {
        this.socket.send(JSON.stringify({
            type: LEAVE,
            user: this.username,
            msg: null,
        }));
    }
}

export default async function init() {
    const chat = new ChatConnection("Anonymous");

    const screen = await window.pnpZone.waw;
    const win = await screen.newWindow({dock: "left", title: "Chat", icon: staticUrl("icon.svg")});
    win.innerHTML = "" +
        "<div class='pnp-chat'>" +
        "   <div class='history'></div>" +
        "   <form>" +
        "       <input>" +
        "   </form>" +
        "</divc>" +
        "";

    const form = win.querySelector("form");
    form.onsubmit = function (event) {
        event.preventDefault();
        chat.send(form[0].value);
        form[0].value = "";
    };

    const history = win.querySelector(".history");
    chat.socket.onmessage = function (message) {
        const {type, user, msg} = JSON.parse(message.data);
        const container = document.createElement("div");
        switch (type) {
            case JOIN:
                container.innerText = `${user} joined the chat`; break;
            case MESSAGE:
                container.innerText = `<${user}> ${msg}`; break;
            case LEAVE:
                container.innerText = `${user} left the chat`; break;
        }
        history.append(container);
    };

    chat.join();
    window.addEventListener("beforeunload", function (event) {
        chat.leave();
    });
}

