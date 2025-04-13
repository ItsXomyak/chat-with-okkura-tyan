const toggleChatListBtn = document.getElementById('toggle-chat-list')
const closeChatListBtn = document.getElementById('close-chat-list')
const chatListSidebar = document.getElementById('chat-list-sidebar')
const chatWindow = document.getElementById('chat-window')
const chatInput = document.getElementById('chat-input')
const sendMessageBtn = document.getElementById('send-message')
const createChatBtn = document.getElementById('create-chat')
const newChatNameInput = document.getElementById('new-chat-name')
const chatList = document.getElementById('chat-list')

let chats = [{ id: Date.now(), name: 'Чат 1', messages: [] }]
let currentChatId = chats[0].id

toggleChatListBtn.addEventListener('click', () => {
	chatListSidebar.classList.toggle('-translate-x-full')
})

closeChatListBtn.addEventListener('click', () => {
	chatListSidebar.classList.add('-translate-x-full')
})

function addMessage(content, isUser = true) {
	const message = {
		id: Date.now(),
		content,
		isUser,
		timestamp: new Date().toLocaleTimeString(),
		isEditing: false,
	}
	const chat = chats.find(c => c.id === currentChatId)
	chat.messages.push(message)
	renderMessages()
}

function renderMessages() {
	const chat = chats.find(c => c.id === currentChatId)
	chatWindow.innerHTML = ''
	chat.messages.forEach(msg => {
		const messageElement = document.createElement('div')
		messageElement.className = `flex ${
			msg.isUser ? 'justify-end' : 'justify-start'
		} mb-4 group`
		messageElement.innerHTML = `
      <div class="flex flex-col items-${msg.isUser ? 'end' : 'start'}">
        <div class="chat-message ${
					msg.isUser
						? 'bg-gradient-to-r from-cyan-500 to-blue-600 text-white'
						: 'bg-gray-700 text-cyan-200'
				} p-3 rounded-lg max-w-xs">
          ${
						msg.isEditing && msg.isUser
							? `
            <textarea class="w-full p-1 rounded text-gray-800">${msg.content}</textarea>
            <div class="flex justify-end gap-2 mt-1">
              <button onclick="saveMessage(${msg.id})" class="text-xs text-green-400 hover:text-green-500">Сохранить</button>
              <button onclick="cancelEdit(${msg.id})" class="text-xs text-red-400 hover:text-red-500">Отмена</button>
            </div>
          `
							: `
            <p>${msg.content}</p>
            <span class="text-xs ${
							msg.isUser ? 'text-cyan-200' : 'text-gray-400'
						}">${msg.timestamp}</span>
          `
					}
        </div>
        ${
					msg.isUser && !msg.isEditing
						? `
          <div class="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <button onclick="editMessage(${msg.id})" class="text-pink-400 hover:text-pink-500"><i class="fas fa-pencil-alt"></i></button>
            <button onclick="deleteMessage(${msg.id})" class="text-pink-400 hover:text-pink-500"><i class="fas fa-trash"></i></button>
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
	const chat = chats.find(c => c.id === currentChatId)
	chat.messages = chat.messages.map(msg =>
		msg.id === id ? { ...msg, isEditing: true } : { ...msg, isEditing: false }
	)
	renderMessages()
}

function saveMessage(id) {
	const textarea = chatWindow.querySelector(`textarea`)
	if (textarea) {
		const newContent = textarea.value.trim()
		if (newContent) {
			const chat = chats.find(c => c.id === currentChatId)
			const msgIndex = chat.messages.findIndex(msg => msg.id === id)
			if (msgIndex !== -1) {
				chat.messages[msgIndex] = {
					...chat.messages[msgIndex],
					content: newContent,
					isEditing: false,
				}
				if (
					msgIndex + 1 < chat.messages.length &&
					!chat.messages[msgIndex + 1].isUser
				) {
					chat.messages.splice(msgIndex + 1, 1)
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
	const chat = chats.find(c => c.id === currentChatId)
	chat.messages = chat.messages.map(msg =>
		msg.id === id ? { ...msg, isEditing: false } : msg
	)
	renderMessages()
}

function deleteMessage(id) {
	const chat = chats.find(c => c.id === currentChatId)
	const msgIndex = chat.messages.findIndex(msg => msg.id === id)
	if (msgIndex !== -1) {
		chat.messages.splice(msgIndex, 1)
		if (msgIndex < chat.messages.length && !chat.messages[msgIndex].isUser) {
			chat.messages.splice(msgIndex, 1)
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

// Chat List Functionality
function renderChatList() {
	chatList.innerHTML = ''
	chats.forEach(chat => {
		const li = document.createElement('li')
		li.className = `p-2 hover:bg-gray-700 rounded-lg flex justify-between items-center ${
			chat.id === currentChatId ? 'bg-gray-600' : ''
		}`
		li.innerHTML = `
      <span class="text-cyan-200">${chat.name}</span>
      <div class="relative group">
        <button class="text-cyan-200 hover:text-pink-300">
          <i class="fas fa-ellipsis-v"></i>
        </button>
        <div class="absolute right-0 top-6 hidden group-hover:block bg-gray-800 rounded-lg shadow-lg p-2">
          <button onclick="renameChat(${chat.id})" class="block text-cyan-200 hover:text-pink-300 w-full text-left">Переименовать</button>
          <button onclick="exportChat(${chat.id})" class="block text-cyan-200 hover:text-pink-300 w-full text-left">Экспорт</button>
          <button onclick="deleteChat(${chat.id})" class="block text-cyan-200 hover:text-pink-300 w-full text-left">Удалить</button>
        </div>
      </div>
    `
		li.addEventListener('click', e => {
			if (!e.target.closest('button')) {
				currentChatId = chat.id
				renderChatList()
				renderMessages()
			}
		})
		chatList.appendChild(li)
	})
}

createChatBtn.addEventListener('click', () => {
	const name = newChatNameInput.value.trim() || `Чат ${chats.length + 1}`
	const newChat = { id: Date.now(), name, messages: [] }
	chats.push(newChat)
	currentChatId = newChat.id
	newChatNameInput.value = ''
	renderChatList()
	renderMessages()
})

function renameChat(id) {
	const newName = prompt(
		'Новое название чата:',
		chats.find(c => c.id === id).name
	)
	if (newName) {
		chats = chats.map(c => (c.id === id ? { ...c, name: newName } : c))
		renderChatList()
	}
}

function exportChat(id) {
	const chat = chats.find(c => c.id === id)
	const data = JSON.stringify(chat, null, 2)
	const blob = new Blob([data], { type: 'application/json' })
	const url = URL.createObjectURL(blob)
	const a = document.createElement('a')
	a.href = url
	a.download = `${chat.name}.json`
	a.click()
	URL.revokeObjectURL(url)
}

function deleteChat(id) {
	if (chats.length > 1) {
		chats = chats.filter(c => c.id !== id)
		currentChatId = chats[0].id
		renderChatList()
		renderMessages()
	} else {
		alert('Нельзя удалить последний чат!')
	}
}

// Initial render
renderChatList()
renderMessages()
