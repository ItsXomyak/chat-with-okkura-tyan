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
		const firstName = document.getElementById('registerFirstName').value
		const lastName = document.getElementById('registerLastName').value

		try {
			const response = await fetch(
				'http://localhost:8080/api/v1/auth/register',
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						username: email.split('@')[0],
						email: email,
						password: password,
						firstName: firstName || 'User',
						lastName: lastName || 'User',
					}),
				}
			)

			if (!response.ok) {
				const errorData = await response.json()
				throw new Error(errorData.error || 'Ошибка при регистрации')
			}

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
			// Перенаправляем на главную страницу
			window.location.href = '/static/index.html'
			console.log('переход')
		} catch (error) {
			console.error('Ошибка при отправке запроса:', error)
			showMessage('error', error.message || 'Ошибка при отправке запроса')
		}
	})

	// Обработка формы сброса пароля
	const resetForm = document.getElementById('resetForm')
	if (resetForm) {
		resetForm.addEventListener('submit', async e => {
			e.preventDefault()

			const email = document.getElementById('resetEmail').value

			try {
				const response = await fetch(
					'http://localhost:8080/api/v1/auth/password-reset-request',
					{
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({
							email: email,
						}),
					}
				)

				if (!response.ok) {
					const errorData = await response.json()
					throw new Error(errorData.error || 'Ошибка при запросе сброса пароля')
				}

				showMessage(
					'success',
					'Инструкции по сбросу пароля отправлены на ваш email.'
				)
				resetForm.reset()
			} catch (error) {
				console.error('Ошибка при отправке запроса:', error)
				showMessage('error', error.message || 'Ошибка при отправке запроса')
			}
		})
	}

	// Обработка формы подтверждения сброса пароля
	const resetConfirmForm = document.getElementById('resetConfirmForm')
	if (resetConfirmForm) {
		resetConfirmForm.addEventListener('submit', async e => {
			e.preventDefault()

			const token = document.getElementById('resetToken').value
			const password = document.getElementById('newPassword').value

			try {
				const response = await fetch(
					'http://localhost:8080/api/v1/auth/password-reset-confirm',
					{
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({
							token: token,
							password: password,
						}),
					}
				)

				if (!response.ok) {
					const errorData = await response.json()
					throw new Error(errorData.error || 'Ошибка при сбросе пароля')
				}

				showMessage(
					'success',
					'Пароль успешно изменен! Теперь вы можете войти с новым паролем.'
				)
				resetConfirmForm.reset()
				// Переключаем на вкладку входа
				document.querySelector('[data-tab="login"]').click()
			} catch (error) {
				console.error('Ошибка при отправке запроса:', error)
				showMessage('error', error.message || 'Ошибка при отправке запроса')
			}
		})
	}

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
