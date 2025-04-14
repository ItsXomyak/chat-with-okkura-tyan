const toggleChatListBtn = document.getElementById('toggle-chat-list')
const closeChatListBtn = document.getElementById('close-chat-list')
const chatListSidebar = document.getElementById('chat-list-sidebar')
const chatWindow = document.getElementById('chat-window')
const chatInput = document.getElementById('chat-input')
const sendMessageBtn = document.getElementById('send-message')
const createChatBtn = document.getElementById('create-chat')
const newChatNameInput = document.getElementById('new-chat-name')
const chatList = document.getElementById('chat-list')
const okkuraSvg = document.getElementById('okkura-svg')
const okkuraMouth = document.getElementById('okkura-mouth')
const searchToggle = document.getElementById('toggle-search')
const searchPanel = document.getElementById('search-panel')
const searchInput = document.getElementById('search-input')
const searchResults = document.getElementById('search-results')
const customizeBtn = document.getElementById('okkura-customize')
const customizePanel = document.getElementById('customize-panel')
const hairColor = document.getElementById('hair-color')
const dressColor = document.getElementById('dress-color')

let chats = []
let currentChatId = null
const userId = 'user123'
let ws = null

customizeBtn.addEventListener('click', () => {
	customizePanel.classList.toggle('show')
})

hairColor.addEventListener('input', () => {
	document
		.querySelectorAll('.okkura-hair')
		.forEach(el => el.setAttribute('fill', hairColor.value))
})

dressColor.addEventListener('input', () => {
	document.querySelector('.okkura-dress').setAttribute('fill', dressColor.value)
})

searchToggle.addEventListener('click', () => {
	searchPanel.classList.toggle('show')
	if (searchPanel.classList.contains('show')) searchInput.focus()
})

searchInput.addEventListener('input', () => {
	const query = searchInput.value.toLowerCase()
	searchResults.innerHTML = ''
	chats
		.filter(chat => chat.name.toLowerCase().includes(query))
		.forEach(chat => {
			const li = document.createElement('li')
			li.className = 'text-cyan-200 cursor-pointer'
			li.textContent = chat.name
			li.addEventListener('click', () => {
				currentChatId = chat.id
				connectWebSocket()
				renderChatList()
				renderMessages()
				searchPanel.classList.remove('show')
				searchInput.value = ''
			})
			searchResults.appendChild(li)
		})
})

async function updateLastMessage(chatId, content) {
	const chat = chats.find(c => c.id === chatId)
	if (chat) chat.lastMessage = content
	renderChatList()
}

function setOkkuraEmotion(emotion) {
	okkuraSvg.classList.remove('okkura-happy', 'okkura-thinking')
	if (emotion === 'happy') {
		okkuraSvg.classList.add('okkura-happy')
		okkuraMouth.setAttribute('d', 'M90 120 Q100 110 110 120')
	} else if (emotion === 'thinking') {
		okkuraSvg.classList.add('okkura-thinking')
		okkuraMouth.setAttribute('d', 'M90 120 Q100 120 110 120')
	} else {
		okkuraMouth.setAttribute('d', 'M90 120 Q100 130 110 120')
	}
}

async function loadChats() {
	try {
		const res = await fetch(`http://localhost:8080/chats?user_id=${userId}`)
		if (!res.ok) throw new Error('Failed to load chats')
		chats = await res.json()
		if (chats.length > 0 && !currentChatId) {
			currentChatId = chats[0].id
			connectWebSocket()
			await renderMessages()
		}
		renderChatList()
	} catch (err) {
		console.error('Load chats error:', err)
		chatWindow.innerHTML = '<p class="text-red-400">Ошибка загрузки чатов</p>'
	}
}

function connectWebSocket() {
	if (ws) ws.close()
	if (!currentChatId) return
	ws = new WebSocket(`ws://localhost:8080/ws/chats/${currentChatId}`)
	ws.onmessage = () => {
		renderMessages()
		setOkkuraEmotion('happy')
		setTimeout(() => setOkkuraEmotion('normal'), 1000)
	}
	ws.onclose = () => setTimeout(connectWebSocket, 1000)
}

async function sendMessage(content, isFromAI = false) {
	if (!currentChatId || !ws) return
	const msg = { user_id: userId, content, is_from_ai: isFromAI }
	ws.send(JSON.stringify(msg))
	await updateLastMessage(currentChatId, content)
	chatInput.value = ''
	setOkkuraEmotion('thinking')
	setTimeout(() => setOkkuraEmotion('normal'), 1000)
}

async function renderMessages() {
	if (!currentChatId) {
		chatWindow.innerHTML = '<p class="text-gray-400">Выберите чат</p>'
		return
	}
	try {
		const res = await fetch(
			`http://localhost:8080/chats/${currentChatId}/messages`
		)
		if (!res.ok) throw new Error('Failed to load messages')
		const messages = await res.json()
		chatWindow.innerHTML = ''
		messages.forEach(msg => {
			msg.isEditing = msg.isEditing || false
			const messageElement = document.createElement('div')
			messageElement.className = `flex ${
				msg.user_id === userId ? 'justify-end' : 'justify-start'
			} mb-4 group chat-message`
			messageElement.innerHTML = `
        <div class="flex flex-col items-${
					msg.user_id === userId ? 'end' : 'start'
				}">
          <div class="chat-message ${
						msg.user_id === userId
							? 'bg-gradient-to-r from-cyan-500 to-blue-600 text-white'
							: 'bg-gray-700 text-cyan-200'
					} p-3 rounded-lg max-w-xs">
            ${
							msg.isEditing && msg.user_id === userId
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
								msg.user_id === userId ? 'text-cyan-200' : 'text-gray-400'
							}">${new Date(msg.created_at).toLocaleTimeString()}</span>
            `
						}
          </div>
          ${
						msg.user_id === userId && !msg.isEditing
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
			chatWindow.scrollTop = chatWindow.scrollHeight
		})
	} catch (err) {
		console.error('Render messages error:', err)
		chatWindow.innerHTML =
			'<p class="text-red-400">Ошибка загрузки сообщений</p>'
	}
}

function renderChatList() {
	chatList.innerHTML = ''
	document.getElementById('current-chat-name').textContent =
		chats.find(c => c.id === currentChatId)?.name || 'Чат с Okkura-тян'
	chats.forEach(chat => {
		const li = document.createElement('li')
		li.className = `p-2 hover:bg-gray-700 rounded-lg flex justify-between items-center ${
			chat.id === currentChatId ? 'bg-gray-600' : ''
		}`
		li.innerHTML = `
  <div class="flex items-center gap-3">
    <i class="fas fa-comment-alt text-pink-300"></i>
    <div>
      <span class="text-cyan-200 font-semibold">${chat.name}</span>
      <p class="text-sm text-gray-400 truncate w-32">${
				chat.lastMessage || 'Нет сообщений'
			}</p>
    </div>
  </div>
  <div class="relative group">
    <button class="text-cyan-200 hover:text-pink-300">
      <i class="fas fa-ellipsis-v"></i>
    </button>
    <div class="absolute right-0 top-6 hidden group-hover:block bg-gray-800 rounded-lg shadow-lg p-2">
      <button onclick="renameChat(${
				chat.id
			})" class="block text-cyan-200 hover:text-pink-300 w-full text-left">Переименовать</button>
      <button onclick="exportChat(${
				chat.id
			})" class="block text-cyan-200 hover:text-pink-300 w-full text-left">Экспорт</button>
      <button onclick="deleteChat(${
				chat.id
			})" class="block text-cyan-200 hover:text-pink-300 w-full text-left">Удалить</button>
    </div>
  </div>
	
`

		li.addEventListener('click', e => {
			if (!e.target.closest('button')) {
				currentChatId = chat.id
				connectWebSocket()
				renderChatList()
				renderMessages()
			}
		})
		chatList.appendChild(li)
	})
}

async function createChat() {
	const name = newChatNameInput.value.trim() || `Чат ${chats.length + 1}`
	try {
		const res = await fetch('http://localhost:8080/chats', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ user_id: userId, name }),
		})
		if (!res.ok) throw new Error('Failed to create chat')
		const chat = await res.json()
		chats.push(chat)
		currentChatId = chat.id
		newChatNameInput.value = ''
		renderChatList()
		connectWebSocket()
		renderMessages()
	} catch (err) {
		console.error('Create chat error:', err)
	}
}

async function editMessage(id) {
	try {
		const res = await fetch(
			`http://localhost:8080/chats/${currentChatId}/messages`
		)
		if (!res.ok) throw new Error('Failed to load messages')
		const messages = await res.json()
		const msg = messages.find(m => m.id === id)
		if (msg) msg.isEditing = true
		renderMessages()
	} catch (err) {
		console.error('Edit message error:', err)
	}
}

async function saveMessage(id) {
	const textarea = chatWindow.querySelector(`textarea`)
	if (!textarea) return
	const newContent = textarea.value.trim()
	if (!newContent) return
	try {
		const res = await fetch(`http://localhost:8080/messages/${id}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ content: newContent }),
		})
		if (!res.ok) throw new Error('Failed to update message')
		const updatedMsg = await res.json()
		if (!updatedMsg.is_from_ai) {
			await sendMessage(
				`О, ты поменял сообщение! Я Okkura-тян, и мне нравится твой новый вайб! 😺`,
				true
			)
		}
		renderMessages()
	} catch (err) {
		console.error('Save message error:', err)
	}
}

async function cancelEdit(id) {
	try {
		const res = await fetch(
			`http://localhost:8080/chats/${currentChatId}/messages`
		)
		if (!res.ok) throw new Error('Failed to load messages')
		const messages = await res.json()
		const msg = messages.find(m => m.id === id)
		if (msg) msg.isEditing = false
		renderMessages()
	} catch (err) {
		console.error('Cancel edit error:', err)
	}
}

async function deleteMessage(id) {
	try {
		const res = await fetch(`http://localhost:8080/messages/${id}`, {
			method: 'DELETE',
		})
		if (!res.ok) throw new Error('Failed to delete message')
		renderMessages()
	} catch (err) {
		console.error('Delete message error:', err)
	}
}

async function renameChat(id) {
	const newName = prompt(
		'Новое название чата:',
		chats.find(c => c.id === id).name
	)
	if (newName) {
		chats = chats.map(c => (c.id === id ? { ...c, name: newName } : c))
		renderChatList()
	}
}

async function exportChat(id) {
	try {
		const res = await fetch(`http://localhost:8080/chats/${id}/messages`)
		if (!res.ok) throw new Error('Failed to load messages')
		const chat = chats.find(c => c.id === id)
		const messages = await res.json()
		const data = JSON.stringify({ ...chat, messages }, null, 2)
		const blob = new Blob([data], { type: 'application/json' })
		const url = URL.createObjectURL(blob)
		const a = document.createElement('a')
		a.href = url
		a.download = `${chat.name}.json`
		a.click()
		URL.revokeObjectURL(url)
	} catch (err) {
		console.error('Export chat error:', err)
	}
}

const themeToggle = document.getElementById('theme-toggle')
themeToggle.addEventListener('click', () => {
	document.body.classList.toggle('light')
	const isLight = document.body.classList.contains('light')
	themeToggle.innerHTML = `<i class="fas fa-${isLight ? 'sun' : 'moon'}"></i>`
	localStorage.setItem('theme', isLight ? 'light' : 'dark')
})

// Инициализация темы
if (localStorage.getItem('theme') === 'light') {
	document.body.classList.add('light')
	themeToggle.innerHTML = '<i class="fas fa-sun"></i>'
}

async function deleteChat(id) {
	if (chats.length <= 1) {
		alert('Нельзя удалить последний чат!')
		return
	}
	try {
		chats = chats.filter(c => c.id !== id)
		if (currentChatId === id) {
			currentChatId = chats[0].id
			connectWebSocket()
		}
		renderChatList()
		renderMessages()
	} catch (err) {
		console.error('Delete chat error:', err)
	}
}

createChatBtn.addEventListener('click', createChat)

sendMessageBtn.addEventListener('click', () => {
	const content = chatInput.value.trim()
	if (content) {
		sendMessage(content)
	}
})

chatInput.addEventListener('keypress', e => {
	if (e.key === 'Enter') {
		sendMessageBtn.click()
	}
})

toggleChatListBtn.addEventListener('click', () => {
	chatListSidebar.classList.toggle('-translate-x-full')
})

closeChatListBtn.addEventListener('click', () => {
	chatListSidebar.classList.add('-translate-x-full')
})

loadChats()
