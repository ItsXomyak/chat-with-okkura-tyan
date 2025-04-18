// Mobile menu toggle
document.addEventListener('DOMContentLoaded', () => {
	const mobileMenuButton = document.querySelector('.md\\:hidden')
	const mobileMenu = document.querySelector('.md\\:flex')

	if (mobileMenuButton) {
		mobileMenuButton.addEventListener('click', () => {
			mobileMenu.classList.toggle('hidden')
		})
	}

	// Smooth scroll for anchor links
	document.querySelectorAll('a[href^="#"]').forEach(anchor => {
		anchor.addEventListener('click', function (e) {
			e.preventDefault()
			const target = document.querySelector(this.getAttribute('href'))
			if (target) {
				target.scrollIntoView({
					behavior: 'smooth',
				})
			}
		})
	})

	// Add hover effect to buttons
	const buttons = document.querySelectorAll('button')
	buttons.forEach(button => {
		button.classList.add('btn-hover-effect')
	})
})

document.addEventListener('mousemove', e => {
	const okkura = document.querySelector('.okkura-character')
	if (okkura) {
		const rect = okkura.getBoundingClientRect()
		const centerX = rect.left + rect.width / 2
		const centerY = rect.top + rect.height / 2

		const angle =
			(Math.atan2(e.clientY - centerY, e.clientX - centerX) * 180) / Math.PI
		okkura.style.transform = `rotate(${angle}deg)`
	}
})

document.addEventListener('DOMContentLoaded', () => {
	const token = localStorage.getItem('token')

	const loginLink = document.getElementById('login-link')
	const profileLink = document.getElementById('profile-link')
	const logoutLink = document.getElementById('logout-link')

	if (token) {
		loginLink?.classList.add('hidden')
		profileLink?.classList.remove('hidden')
		logoutLink?.classList.remove('hidden')
	} else {
		loginLink?.classList.remove('hidden')
		profileLink?.classList.add('hidden')
		logoutLink?.classList.add('hidden')
	}

	// Обработка выхода
	logoutLink?.addEventListener('click', e => {
		e.preventDefault()
		localStorage.removeItem('token')
		location.reload()
	})
})	
