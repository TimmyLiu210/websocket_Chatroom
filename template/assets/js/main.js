const LEFT = "left";
const RIGHT = "right";

const EVENT_MESSAGE = "message"
const EVENT_OTHER = "other"

var PERSON_IMG = "https://www.flaticon.com/svg/static/icons/svg/3408/3408584.svg";
var PERSON_NAME = "Guest" + Math.floor(Math.random() * 1000);

var url = "ws://" + window.location.host + "/ws?id=" + PERSON_NAME;
var ws = new WebSocket(url);

var chatroom = document.getElementsByClassName("msger-chat")
var text = document.getElementById("msg");
var send = document.getElementById("send")

send.onclick = function (e) {
    handleMessageEvent()
}

text.onkeydown = function (e) {
    if (e.keyCode === 20 && text.value !== "") {
        handleMessageEvent()
    }
};

ws.onmessage = function (e) {
    var m = JSON.parse(e.data)
    var msg = ""
    switch (m.event) {
        case EVENT_MESSAGE:
            if (m.name == PERSON_NAME) {
                msg = getMessage(m.name, m.photo, RIGHT, m.content);
            } else {
                msg = getMessage(m.name, m.photo, LEFT, m.content);
            }
            break;
        case EVENT_OTHER:
            if (m.name != PERSON_NAME) {
                msg = getEventMessage(m.name + " " + m.content)
            } else {
                msg = getEventMessage("您已" + m.content)
            }
            break;
    }
    insertMsg(msg, chatroom[0]);
};

function handleMessageEvent() {
    ws.send(JSON.stringify({
        "event": "message",
        "photo": PERSON_IMG,
        "name": PERSON_NAME,
        "content": text.value,
    }));
    text.value = "";
}

function getEventMessage(msg) {
    var msg = `<div class="msg-left">${msg}</div>`
    return msg
}

function getMessage(name, img, side, text) {
    const d = new Date()

    var msg = `
    <div class="msg ${side}-msg">
      <div class="msg-img" style="background-image: url(${img})"></div>

      <div class="msg-bubble">
        <div class="msg-info">
          <div class="msg-info-name">${name}</div>
          <div class="msg-info-time">${d.getFullYear}/${d.getMonth()}/${d.getDay()}${d.getHours()}:${d.getMinutes()}</div>
        </div>

        <div class="msg-text">${text}</div>
      </div>
    </div>
  `
    return msg;
}

function insertMsg(msg, domObj) {
    domObj.insertAdjacentHTML("beforeend", msg);
    domObj.scrollTop += 500;
}