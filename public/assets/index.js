var myEditor = document.getElementById('myEditor');
var editor = CodeMirror.fromTextArea(myEditor, {
  lineNumbers: true,
  mode: "ruby"
});
editor.save();

const connectButton = document.getElementById('connect');
let port;
if ('serial' in navigator) {
  connectButton.addEventListener('click', function () {
    console.log("click");

    if (port) {
      port.close();
      port = undefined;
      connectButton.innerText = 'Connect';
    }
    else {
      getReader();
    }
  });

  connectButton.disabled = false;
}

const runButton = document.getElementById('run');
const outputArea = document.getElementById('output');

// refs: https://wicg.github.io/serial/#dom-serialportinfo
async function getReader() {
  port = await navigator.serial.requestPort({});
  await port.open({ baudRate: 115200 });
  connectButton.innerText = '🔌 Disconnect';

  runButton.addEventListener('click', (event) => {
    const code = editor.getValue();
    outputArea.innerText = 'Compiling...'
    $.ajax({
      type:     "POST",
      url:      "http://localhost:3000/compile/",
      data: JSON.stringify({
        "source_code": code
      }),
      dataType: "json",
      xhrFields: {
          withCredentials: true
      },

      }).done(function(data) {
        const byte_code = data["byte_code"]
        if (port && port.writable) {
          bytes = new Uint8Array(byte_code.match(/.{1,2}/g).map(v => parseInt(v, 16)))
          const writer = port.writable.getWriter();
          console.log(byte_code)
          const encoder = new TextEncoder();
          writer.write(encoder.encode(byte_code));
          writer.releaseLock();
        }

      }).fail(function(data) {
          console.log('Ajax fail (communication error)');
      });
  });
}
