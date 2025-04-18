let headers = {}

document.addEventListener('DOMContentLoaded', async () => {
	try {
		// Проверка авторизации
		const token = localStorage.getItem('token')
		if (!token) {
			showNotification('❌ Вы не авторизованы.', 'error')
			window.location.href = 'auth.html'
			return
		}

		headers = {
			Authorization: `Bearer ${token}`,
			'Content-Type': 'application/json',
		}

		// Показываем состояние загрузки
		showLoadingState()

		// Загружаем все данные параллельно
		await Promise.all([
			loadProfile().catch(err =>
				console.error('Ошибка загрузки профиля:', err)
			),
			loadStats().catch(err =>
				console.error('Ошибка загрузки статистики:', err)
			),
			loadHistory().catch(err =>
				console.error('Ошибка загрузки истории:', err)
			),
		])
	} catch (err) {
		console.error('Общая ошибка:', err)
		showNotification('Произошла ошибка при загрузке данных', 'error')
	}
})

function showLoadingState() {
	// Безопасное обновление элемента
	const updateElementLoading = id => {
		const element = document.getElementById(id)
		if (element) {
			element.innerHTML = `
				<div class="animate-pulse bg-purple-200 h-6 w-32 rounded"></div>
			`
		}
	}

	// Обновляем основную информацию
	updateElementLoading('user-name')
	updateElementLoading('user-email')
	updateElementLoading('user-id')
	updateElementLoading('user-role')

	// Обновляем статистику
	const statsElements = ['total-amount', 'donate-count', 'subscription']
	statsElements.forEach(id => updateElementLoading(id))

	// Обновляем аватар
	const avatarImg = document.getElementById('profile-avatar')
	if (avatarImg) {
		avatarImg.style.opacity = '0.5'
	}

	// Очищаем историю донатов
	const donateHistory = document.getElementById('donate-history')
	if (donateHistory) {
		donateHistory.innerHTML = `
			<tr>
				<td colspan="4" class="py-4 text-center">
					<div class="animate-pulse flex justify-center">
						<div class="h-4 bg-purple-200 rounded w-3/4"></div>
					</div>
				</td>
			</tr>
		`
	}
}

function showNotification(message, type = 'success') {
	const notification = document.createElement('div')
	notification.className = `fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg transform transition-all duration-500 z-50 ${
		type === 'error'
			? 'bg-red-100 border-red-400 text-red-700'
			: 'bg-green-100 border-green-400 text-green-700'
	}`
	notification.style.opacity = '0'
	notification.innerHTML = message
	document.body.appendChild(notification)

	// Анимация появления
	requestAnimationFrame(() => {
		notification.style.opacity = '1'
	})

	// Автоматическое скрытие
	setTimeout(() => {
		notification.style.opacity = '0'
		setTimeout(() => notification.remove(), 500)
	}, 5000)
}

// Добавим функцию для безопасного обновления элементов
function safeUpdateElement(id, content, defaultValue = '—') {
	const element = document.getElementById(id)
	if (element) {
		if (typeof content === 'undefined' || content === null) {
			element.textContent = defaultValue
		} else {
			element.textContent = content
		}
	}
}

// Обновим функцию loadProfile
async function loadProfile() {
	try {
		const res = await fetch('http://localhost:8082/me/profile', {
			method: 'GET',
			headers,
		})

		if (!res.ok) {
			const errorData = await res.json().catch(() => ({}))
			throw new Error(errorData.error || 'Ошибка при загрузке профиля')
		}

		const profile = await res.json()
		console.log('Данные профиля:', profile)

		// Безопасно обновляем все элементы
		const fullName =
			[profile.first_name, profile.last_name].filter(Boolean).join(' ') ||
			profile.nickname ||
			'Пользователь'

		safeUpdateElement('user-name', fullName)
		safeUpdateElement('user-email', profile.email)
		safeUpdateElement('user-id', profile.user_id)
		safeUpdateElement(
			'user-role',
			profile.role === 'admin' ? 'Администратор' : 'Пользователь'
		)

		// Обновляем поля формы
		const formFields = {
			nickname: profile.nickname,
			'first-name': profile.first_name,
			'last-name': profile.last_name,
			'avatar-url': profile.avatar_url,
		}

		Object.entries(formFields).forEach(([id, value]) => {
			const input = document.getElementById(id)
			if (input) {
				input.value = value || ''
			}
		})

		// Обновляем аватар
		const avatarImg = document.getElementById('profile-avatar')
		if (avatarImg) {
			avatarImg.src = profile.avatar_url || 'images/default-avatar.png'
			avatarImg.style.opacity = '1'
		}
	} catch (err) {
		console.error('Ошибка загрузки профиля:', err)
		showNotification('Ошибка при загрузке профиля: ' + err.message, 'error')
	}
}

async function loadStats() {
	try {
		const res = await fetch('http://localhost:8082/me/summary', { headers })
		if (!res.ok) throw new Error('Ошибка при получении статистики')
		const stats = await res.json()

		// Форматируем сумму с двумя десятичными знаками
		const formattedAmount = Number(stats.totalAmount).toFixed(2)
		document.getElementById('total-amount').textContent = `${formattedAmount}₽`
		document.getElementById('donate-count').textContent = stats.donateCount || 0
		document.getElementById('subscription').textContent =
			stats.subscription || '—'
	} catch (err) {
		console.error('Ошибка загрузки статистики:', err)
		// Показываем ошибку в интерфейсе
		document.getElementById('total-amount').textContent = 'Ошибка'
		document.getElementById('donate-count').textContent = 'Ошибка'
		document.getElementById('subscription').textContent = 'Ошибка'
	}
}

async function loadHistory() {
	try {
		const res = await fetch('http://localhost:8082/me/donates', { headers })
		if (!res.ok) throw new Error('Ошибка при получении донатов')
		const history = await res.json()

		const table = document.getElementById('donate-history')
		table.innerHTML = ''

		if (history.length === 0) {
			table.innerHTML =
				'<tr><td colspan="4" class="py-4 text-center text-purple-500">Нет донатов</td></tr>'
			return
		}

		history.forEach(d => {
			const tr = document.createElement('tr')
			tr.innerHTML = `
				<td class="py-2 pr-4">${d.amount}₽</td>
				<td class="py-2 pr-4">${d.method}</td>
				<td class="py-2 pr-4">${d.message || '—'}</td>
				<td class="py-2">${new Date(d.created_at).toLocaleString('ru-RU')}</td>
			`
			table.appendChild(tr)
		})
	} catch (err) {
		console.error('Ошибка загрузки истории донатов:', err)
	}
}

async function updateProfile() {
	const status = document.getElementById('update-status')
	const updateButton = event.target
	const originalButtonText = updateButton.textContent

	try {
		// Показываем состояние загрузки на кнопке
		updateButton.disabled = true
		updateButton.innerHTML = `
			<svg class="animate-spin h-5 w-5 mr-2 inline" viewBox="0 0 24 24">
				<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
				<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
			</svg>
			Сохранение...
		`

		const data = {
			nickname: document.getElementById('nickname').value,
			first_name: document.getElementById('first-name').value,
			last_name: document.getElementById('last-name').value,
			avatar_url: document.getElementById('avatar-url').value,
		}

		const res = await fetch('http://localhost:8082/me/profile', {
			method: 'PATCH',
			headers,
			body: JSON.stringify(data),
		})

		if (!res.ok) {
			throw new Error('Ошибка при обновлении профиля')
		}

		// Обновляем информацию на странице
		await loadProfile()

		showNotification('✅ Профиль успешно обновлен!')
	} catch (err) {
		console.error('Ошибка при обновлении профиля:', err)
		showNotification('❌ ' + (err.message || 'Ошибка при обновлении'), 'error')
	} finally {
		// Возвращаем кнопку в исходное состояние
		updateButton.disabled = false
		updateButton.textContent = originalButtonText
	}
}

function handleError(error) {
	if (error.message.includes('401') || error.message.includes('403')) {
		localStorage.removeItem('token')
		showNotification(
			'❌ Сессия истекла. Необходима повторная авторизация',
			'error'
		)
		setTimeout(() => {
			window.location.href = 'auth.html'
		}, 2000)
		return
	}

	showNotification('❌ ' + error.message, 'error')
}

// Добавляем CSS для анимаций
const style = document.createElement('style')
style.textContent = `
	.fade-in {
		animation: fadeIn 0.5s ease-in-out;
	}
	
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
`
document.head.appendChild(style)
