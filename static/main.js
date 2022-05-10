function staticUrl(url) {
    return `${window.pnpZone.static}plugins/pnp-chat/${url}`;
}

export default async function init() {
    const {host, protocol} = window.location;
    const socket = new WebSocket(`${protocol.replace("http", "ws")}//${host}/chat/socket`);

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
        socket.send(form[0].value);
        form[0].value = "";
    };

    const history = win.querySelector(".history");
    socket.onmessage = function (message) {
        const container = document.createElement("div");
        container.innerText = message.data;
        history.append(container);
    };
}

