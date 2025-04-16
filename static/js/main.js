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
