document.addEventListener('DOMContentLoaded', () => {
	// Обработка переключения вкладок
	const tabButtons = document.querySelectorAll('.tab-btn')
	const tabContents = document.querySelectorAll('.tab-content')

	tabButtons.forEach(button => {
		button.addEventListener('click', () => {
			const tabId = button.getAttribute('data-tab')

			// Удаляем активный класс у всех кнопок и контента
			tabButtons.forEach(btn => btn.classList.remove('active'))
			tabContents.forEach(content => content.classList.remove('active'))

			// Добавляем активный класс выбранной кнопке и контенту
			button.classList.add('active')
			document.getElementById(tabId).classList.add('active')
		})
	})

	// Обработка формы регистрации
	const registerForm = document.getElementById('registerForm')
	registerForm.addEventListener('submit', async e => {
		e.preventDefault()

		const email = document.getElementById('registerEmail').value
		const password = document.getElementById('registerPassword').value

		try {
			const response = await fetch(
				'http://localhost:8080/api/v1/auth/register',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						username: email.split('@')[0], // Генерируем username из email
						email: email,
						password: password,
						firstName: 'User', // Можно добавить поля для ввода имени
						lastName: 'User', // Можно добавить поля для ввода фамилии
					}),
				}
			)

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || 'Ошибка при регистрации')
			}

			const data = await response.json()
			showMessage(
				'success',
				'Регистрация успешна! Проверьте ваш email для подтверждения.'
			)
			registerForm.reset()
		} catch (error) {
			console.error('Ошибка при отправке запроса:', error)
			showMessage('error', error.message || 'Ошибка при отправке запроса')
		}
	})

	// Обработка формы входа
	const loginForm = document.getElementById('loginForm')
	loginForm.addEventListener('submit', async e => {
		e.preventDefault()

		const email = document.getElementById('loginEmail').value
		const password = document.getElementById('loginPassword').value

		try {
			const response = await fetch('http://localhost:8080/api/v1/auth/login', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					email: email,
					password: password,
				}),
			})

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || 'Ошибка при входе')
			}

			const data = await response.json()
			// Сохраняем токен в localStorage
			localStorage.setItem('token', data.token)
			showMessage('success', 'Вход выполнен успешно!')
			loginForm.reset()
			// Здесь можно добавить перенаправление на защищенную страницу
		} catch (error) {
			console.error('Ошибка при отправке запроса:', error)
			showMessage('error', error.message || 'Ошибка при отправке запроса')
		}
	})

	// Функция для отображения сообщений
	function showMessage(type, text) {
		const messageElement = document.createElement('div')
		messageElement.className = `${type}-message`
		messageElement.textContent = text
		messageElement.style.display = 'block'

		// Удаляем предыдущие сообщения
		const existingMessages = document.querySelectorAll(
			'.error-message, .success-message'
		)
		existingMessages.forEach(msg => msg.remove())

		// Добавляем новое сообщение
		document.querySelector('.auth-container').appendChild(messageElement)

		// Автоматически скрываем сообщение через 5 секунд
		setTimeout(() => {
			messageElement.remove()
		}, 5000)
	}

	// Проверяем наличие токена верификации в URL
	const urlParams = new URLSearchParams(window.location.search)
	const verificationToken = urlParams.get('token')

	if (verificationToken) {
		verifyEmail(verificationToken)
	}

	// Функция для верификации email
	async function verifyEmail(token) {
		try {
			const response = await fetch(
				'http://localhost:8080/api/v1/auth/confirm',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						token: token,
					}),
				}
			)

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || 'Ошибка при подтверждении email')
			}

			// Перенаправляем на страницу подтверждения
			window.location.href = '/static/email-confirmed.html'
		} catch (error) {
			console.error('Ошибка при отправке запроса:', error)
			showMessage('error', error.message || 'Ошибка при отправке запроса')
		}
	}
})
