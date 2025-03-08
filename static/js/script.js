// Відкриття/закриття модального вікна реєстрації
const modal = document.getElementById('registrationModal');
const openBtn = document.getElementById('openModal');
const closeBtn = document.querySelector('.close');

if (openBtn && closeBtn) {
    openBtn.addEventListener('click', () => {
        modal.style.display = 'block';
    });

    closeBtn.addEventListener('click', () => {
        modal.style.display = 'none';
    });
}

// Реєстрація
const registrationForm = document.getElementById('registrationForm');
if (registrationForm) {
    registrationForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            username: e.target.username.value,
            email: e.target.email.value,
            password: e.target.password.value
        };

        try {
            const response = await fetch('/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            const result = await response.json();
            if (result.success) {
                alert('Реєстрація успішна!');
                modal.style.display = 'none';
                window.location.href = '/login'; // Перенаправлення на сторінку входу
            } else {
                alert('Помилка: ' + result.error);
            }
        } catch (error) {
            alert('Помилка мережі');
        }
    });
}

// Вхід
const loginForm = document.getElementById('loginForm');
if (loginForm) {
    loginForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const formData = {
            username: e.target.username.value,
            password: e.target.password.value
        };

        try {
            const response = await fetch('/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData)
            });

            const result = await response.json();
            if (result.success) {
                window.location.href = '/dashboard'; // Перенаправлення на панель управління
            } else {
                alert('Помилка: ' + result.error);
            }
        } catch (error) {
            alert('Помилка мережі');
        }
    });
}

// Вийти
function logout() {
    fetch('/logout', { method: 'POST' })
        .then(() => window.location.href = '/');
}