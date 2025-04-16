const themeToggle = document.getElementById('theme-toggle')
themeToggle.addEventListener('click', () => {
	document.body.classList.toggle('light')
	const isLight = document.body.classList.contains('light')
	themeToggle.innerHTML = `<i class="fas fa-${isLight ? 'sun' : 'moon'}"></i>`
	localStorage.setItem('theme', isLight ? 'light' : 'dark')
	console.log('123')
})

// Инициализация темы
if (localStorage.getItem('theme') === 'light') {
	document.body.classList.add('light')
	themeToggle.innerHTML = '<i class="fas fa-sun"></i>'
}
