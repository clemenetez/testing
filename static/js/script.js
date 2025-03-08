// Відкриття/закриття модального вікна
const modal = document.getElementById('registrationModal');
const openBtn = document.getElementById('openModal');
const closeBtn = document.querySelector('.close');

openBtn.addEventListener('click', () => {
    modal.style.display = 'block';
});

closeBtn.addEventListener('click', () => {
    modal.style.display = 'none';
});

// Відправка форми
document.getElementById('registrationForm').addEventListener('submit', async (e) => {
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
        } else {
            alert('Помилка: ' + result.error);
        }
    } catch (error) {
        alert('Помилка мережі');
    }
});