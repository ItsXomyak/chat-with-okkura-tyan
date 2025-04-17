let headers = {}

document.addEventListener('DOMContentLoaded', async () => {
	const token = localStorage.getItem('token')
	if (!token) {
		alert('❌ Вы не авторизованы.')
		window.location.href = 'auth.html'
		return
	}

	headers = {
		Authorization: 'Bearer ' + token,
		'Content-Type': 'application/json',
	}

	await loadProfile()
	await loadStats()
	await loadHistory()
})

async function loadProfile() {
	try {
		const res = await fetch('http://localhost:8082/me/profile', { headers })
		if (!res.ok) throw new Error('Ошибка при загрузке профиля')
		const profile = await res.json()

		document.getElementById(
			'user-name'
		).textContent = `${profile.first_name} ${profile.last_name}`
		document.getElementById('user-email').textContent = profile.email || '—'
		document.getElementById('user-id').textContent = profile.user_id || '—'
		document.getElementById('user-role').textContent = 'пользователь'

		document.getElementById('nickname').value = profile.nickname || ''
		document.getElementById('first-name').value = profile.first_name || ''
		document.getElementById('last-name').value = profile.last_name || ''
		document.getElementById('avatar-url').value = profile.avatar_url || ''
	} catch (err) {
		console.error('Ошибка загрузки профиля:', err)
	}
}

async function loadStats() {
	try {
		const res = await fetch('http://localhost:8082/me/summary', { headers })
		if (!res.ok) throw new Error('Ошибка при получении статистики')
		const stats = await res.json()

		document.getElementById('total-amount').textContent = `${
			stats.totalAmount || 0
		}₽`
		document.getElementById('donate-count').textContent = stats.donateCount || 0
		document.getElementById('subscription').textContent =
			stats.subscription || '—'
	} catch (err) {
		console.error('Ошибка загрузки статистики:', err)
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
	const nickname = document.getElementById('nickname').value
	const firstName = document.getElementById('first-name').value
	const lastName = document.getElementById('last-name').value
	const avatarURL = document.getElementById('avatar-url').value
	const status = document.getElementById('update-status')

	try {
		const res = await fetch('http://localhost:8082/me/profile', {
			method: 'PATCH',
			headers,
			body: JSON.stringify({
				nickname,
				first_name: firstName,
				last_name: lastName,
				avatar_url: avatarURL,
			}),
		})

		if (res.ok) {
			status.textContent = '✅ Профиль обновлён!'
			status.className = 'text-green-600 mt-2 text-center'
			await loadProfile()
		} else {
			status.textContent = '❌ Ошибка при обновлении'
			status.className = 'text-red-600 mt-2 text-center'
		}
	} catch (err) {
		console.error(err)
		status.textContent = '❌ Сервер недоступен'
		status.className = 'text-red-600 mt-2 text-center'
	}
}
