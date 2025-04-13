const toggleChatListBtn = document.getElementById('toggle-chat-list')
const closeChatListBtn = document.getElementById('close-chat-list')
const chatListSidebar = document.getElementById('chat-list-sidebar')
const chatWindow = document.getElementById('chat-window')
const chatInput = document.getElementById('chat-input')
const sendMessageBtn = document.getElementById('send-message')

toggleChatListBtn.addEventListener('click', () => {
	chatListSidebar.classList.toggle('-translate-x-full')
})

closeChatListBtn.addEventListener('click', () => {
	chatListSidebar.classList.add('-translate-x-full')
})

let messages = []

function addMessage(content, isUser = true) {
	const message = {
		id: Date.now(),
		content,
		isUser,
		timestamp: new Date().toLocaleTimeString(),
		isEditing: false,
	}
	messages.push(message)
	renderMessages()
}

function renderMessages() {
	chatWindow.innerHTML = ''
	messages.forEach(msg => {
		const messageElement = document.createElement('div')
		messageElement.className = `flex ${
			msg.isUser ? 'justify-end' : 'justify-start'
		} mb-4 group`
		messageElement.innerHTML = `
      <div class="flex flex-col items-${msg.isUser ? 'end' : 'start'}">
        <div class="${
					msg.isUser ? 'bg-blue-500 text-white' : 'bg-gray-200 text-gray-800'
				} p-3 rounded-lg max-w-xs">
          ${
						msg.isEditing && msg.isUser
							? `
            <textarea class="w-full p-1 rounded text-gray-800">${msg.content}</textarea>
            <div class="flex justify-end gap-2 mt-1">
              <button onclick="saveMessage(${msg.id})" class="text-xs text-green-500 hover:text-green-600">Сохранить</button>
              <button onclick="cancelEdit(${msg.id})" class="text-xs text-red-500 hover:text-red-600">Отмена</button>
            </div>
          `
							: `
            <p>${msg.content}</p>
            <span class="text-xs ${
							msg.isUser ? 'text-blue-200' : 'text-gray-500'
						}">${msg.timestamp}</span>
          `
					}
        </div>
        ${
					msg.isUser && !msg.isEditing
						? `
          <div class="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button onclick="editMessage(${msg.id})" class="text-blue-600 hover:text-blue-800"><i class="fas fa-pencil-alt"></i></button>
            <button onclick="deleteMessage(${msg.id})" class="text-blue-600 hover:text-blue-800"><i class="fas fa-trash"></i></button>
          </div>
        `
						: ''
				}
      </div>
    `
		chatWindow.appendChild(messageElement)
	})
	chatWindow.scrollTop = chatWindow.scrollHeight
}

function editMessage(id) {
	messages = messages.map(msg =>
		msg.id === id ? { ...msg, isEditing: true } : { ...msg, isEditing: false }
	)
	renderMessages()
}

function saveMessage(id) {
	const textarea = chatWindow.querySelector(`textarea`)
	if (textarea) {
		const newContent = textarea.value.trim()
		if (newContent) {
			const msgIndex = messages.findIndex(msg => msg.id === id)
			if (msgIndex !== -1) {
				messages[msgIndex] = {
					...messages[msgIndex],
					content: newContent,
					isEditing: false,
				}
				if (msgIndex + 1 < messages.length && !messages[msgIndex + 1].isUser) {
					messages.splice(msgIndex + 1, 1)
				}
				setTimeout(() => {
					addMessage(
						`О, ты поменял сообщение! Я Okkura-тян, и мне нравится твой новый вайб! 😺`,
						false
					)
				}, 1000)
			}
		}
	}
	renderMessages()
}

function cancelEdit(id) {
	messages = messages.map(msg =>
		msg.id === id ? { ...msg, isEditing: false } : msg
	)
	renderMessages()
}

function deleteMessage(id) {
	const msgIndex = messages.findIndex(msg => msg.id === id)
	if (msgIndex !== -1) {
		messages.splice(msgIndex, 1)
		if (msgIndex < messages.length && !messages[msgIndex].isUser) {
			messages.splice(msgIndex, 1)
		}
	}
	renderMessages()
}

sendMessageBtn.addEventListener('click', () => {
	const content = chatInput.value.trim()
	if (content) {
		addMessage(content)
		chatInput.value = ''
		setTimeout(() => {
			addMessage('Йо, я Okkura-тян! Классно болтать! 😎', false)
		}, 1000)
	}
})

chatInput.addEventListener('keypress', e => {
	if (e.key === 'Enter') {
		sendMessageBtn.click()
	}
})

const chatList = document.getElementById('chat-list')
;['Чат 1', 'Чат 2', 'Чат 3'].forEach(chat => {
	const li = document.createElement('li')
	li.className =
		'p-2 hover:bg-gray-100 rounded-lg flex justify-between items-center'
	li.innerHTML = `
    <span>${chat}</span>
    <button class="text-gray-500 hover:text-gray-700">
      <i class="fas fa-ellipsis-v"></i>
    </button>
  `
	chatList.appendChild(li)
})
