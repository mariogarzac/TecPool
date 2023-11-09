// Get elements
const messagesContainer = document.getElementById('messages-container');
const tripId = messagesContainer.getAttribute('data-trip-id');
const userId = messagesContainer.getAttribute('data-user-id');

const messageInput = document.getElementById('message-input');
const sendButton = document.getElementById('send-button');

// Create websocket connection
const socket = new WebSocket(`ws://localhost:3001/ws/chat/${tripId}/${userId}`);

// Event listener for sending a message
sendButton.addEventListener('click', function() {
    const message = messageInput.value;

    const messageData = {
        userId : userId,
        message: message,
    }

    const jsonMessage = JSON.stringify(messageData)

    socket.send(jsonMessage)

    messageInput.value = '';
});

// Send message with Enter key
messageInput.addEventListener('keypress', function(e) {
    if (e.key === 'Enter') {
        sendButton.click();
    }
});

// Handling incoming messages
socket.onmessage = function(event) {
    const jsonMessage = JSON.parse(event.data);

    // Create a div element for the message
    const messageElement = document.createElement('div');

    // Check if the message is from the user or another participant
    if (jsonMessage.self) {
        // Style the user's messages on the right side with a blue bubble
        messageElement.classList.add('flex', 'flex-row-reverse', 'items-end', 'mb-2', 'pt-2');
        messageElement.innerHTML = `
            <div class="bg-blue-500 text-white px-4 py-2 rounded-lg">
            ${jsonMessage.message}
            </div>
            `;
    } else {
        // Style other participants' messages on the left side with a grey bubble
        messageElement.classList.add('flex', 'items-start', 'mb-2', 'pt-2');
        messageElement.innerHTML = `
            <div class="bg-gray-300 text-black px-4 py-2 rounded-lg">
            ${jsonMessage.message}
            </div>
            `;
    }

    // Append the message element to the messagesContainer
    messagesContainer.appendChild(messageElement);

    // Scroll the messagesContainer to the bottom to show the latest message
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
};

// Handling WebSocket closure
socket.onclose = function(event) {
    if (event.wasClean) {
        alert(`Connection closed cleanly, code=${event.code}, reason=${event.reason}`);
    } else {
        // e.g., server process killed or network down
        alert('Connection died');
    }
};

// Handling WebSocket errors
socket.onerror = function(error) {
    alert(`[error] ${error.message}`);
};
