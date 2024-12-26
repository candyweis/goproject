// static/js/shooter.js
// Инициализация Canvas
const canvas = document.getElementById('shooterCanvas');
const ctx = canvas.getContext('2d');

// Установка размеров Canvas
canvas.width = window.innerWidth;
canvas.height = window.innerHeight;

// Обработка изменения размеров окна
window.addEventListener('resize', () => {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
});

// Звуковые эффекты (опционально)
// const shootSound = new Audio('/static/sounds/shoot.wav');
// const enemyHitSound = new Audio('/static/sounds/enemy_hit.wav');
// const playerHitSound = new Audio('/static/sounds/player_hit.wav');

// Игрок
class Player {
    constructor() {
        this.width = 50;
        this.height = 50;
        this.x = canvas.width / 2 - this.width / 2;
        this.y = canvas.height - this.height - 30;
        this.speed = 7;
        this.color = 'blue';
        this.moveLeft = false;
        this.moveRight = false;
        this.moveUp = false;
        this.moveDown = false;
    }

    draw() {
        ctx.fillStyle = this.color;
        ctx.fillRect(this.x, this.y, this.width, this.height);
    }

    update() {
        if (this.moveLeft && this.x > 0) {
            this.x -= this.speed;
        }
        if (this.moveRight && this.x + this.width < canvas.width) {
            this.x += this.speed;
        }
        if (this.moveUp && this.y > 0) {
            this.y -= this.speed;
        }
        if (this.moveDown && this.y + this.height < canvas.height) {
            this.y += this.speed;
        }
    }
}

// Снаряд
class Bullet {
    constructor(x, y) {
        this.width = 5;
        this.height = 15;
        this.x = x + 22.5; // Центр относительно игрока
        this.y = y;
        this.speed = 10;
        this.color = 'yellow';
    }

    draw() {
        ctx.fillStyle = this.color;
        ctx.fillRect(this.x, this.y, this.width, this.height);
    }

    update() {
        this.y -= this.speed;
    }

    isOffScreen() {
        return this.y + this.height < 0;
    }
}

// Враг
class Enemy {
    constructor() {
        this.width = 40;
        this.height = 40;
        this.x = Math.random() * (canvas.width - this.width);
        this.y = -this.height;
        this.speed = 3 + Math.random() * 2;
        this.color = 'red';
    }

    draw() {
        ctx.fillStyle = this.color;
        ctx.fillRect(this.x, this.y, this.width, this.height);
    }

    update() {
        this.y += this.speed;
    }

    isOffScreen() {
        return this.y > canvas.height;
    }
}

// Основные переменные
const player = new Player();
const bullets = [];
const enemies = [];
let enemySpawnInterval = 1500; // Время между спавнами врагов (мс)
let lastEnemySpawn = Date.now();
let score = 0;
let highScore = 0;
let lives = 3;
let gameOver = false;

// Обработка нажатий клавиш
window.addEventListener('keydown', (e) => {
    switch (e.code) {
        case 'ArrowLeft':
            player.moveLeft = true;
            break;
        case 'ArrowRight':
            player.moveRight = true;
            break;
        case 'ArrowUp':
            player.moveUp = true;
            break;
        case 'ArrowDown':
            player.moveDown = true;
            break;
        case 'Space':
            if (!gameOver) {
                bullets.push(new Bullet(player.x, player.y));
                // shootSound.currentTime = 0;
                // shootSound.play();
            }
            break;
    }
});

window.addEventListener('keyup', (e) => {
    switch (e.code) {
        case 'ArrowLeft':
            player.moveLeft = false;
            break;
        case 'ArrowRight':
            player.moveRight = false;
            break;
        case 'ArrowUp':
            player.moveUp = false;
            break;
        case 'ArrowDown':
            player.moveDown = false;
            break;
    }
});

// Проверка коллизий между двумя прямоугольниками
function isColliding(a, b) {
    return (
        a.x < b.x + b.width &&
        a.x + a.width > b.x &&
        a.y < b.y + b.height &&
        a.y + a.height > b.y
    );
}

// Рендеринг счета и жизней
function renderHUD() {
    const scoreElement = document.getElementById('score');
    scoreElement.textContent = `Очки: ${score} | Лучший: ${highScore} | Жизни: ${lives}`;
}

// Обработка окончания игры
function endGame() {
    gameOver = true;
    const gameOverElement = document.getElementById('gameOver');
    gameOverElement.textContent = `Игра окончена! Ваш счет: ${score}. Нажмите F5 для перезапуска.`;
    gameOverElement.style.display = 'block';

    // Сохранение лучшего счета
    if (score > highScore) {
        highScore = score;
        localStorage.setItem('highScore', highScore);
    }

    // Отправка очков на сервер
    sendScore();
}

// Главная игровая функция
function gameLoop() {
    if (gameOver) {
        return;
    }

    ctx.clearRect(0, 0, canvas.width, canvas.height);

    // Обновление и рисование игрока
    player.update();
    player.draw();

    // Обновление и рисование снарядов
    for (let i = bullets.length - 1; i >= 0; i--) {
        bullets[i].update();
        bullets[i].draw();

        if (bullets[i].isOffScreen()) {
            bullets.splice(i, 1);
        }
    }

    // Спавн врагов
    const now = Date.now();
    if (now - lastEnemySpawn > enemySpawnInterval) {
        enemies.push(new Enemy());
        lastEnemySpawn = now;
    }

    // Обновление и рисование врагов
    for (let i = enemies.length - 1; i >= 0; i--) {
        enemies[i].update();
        enemies[i].draw();

        if (enemies[i].isOffScreen()) {
            enemies.splice(i, 1);
        }
    }

    // Проверка коллизий
    for (let i = enemies.length - 1; i >= 0; i--) {
        // Проверка столкновения с игроком
        if (isColliding(player, enemies[i])) {
            // playerHitSound.play();
            lives -= 1;
            renderHUD();

            if (lives <= 0) {
                endGame();
            }
            enemies.splice(i, 1);
            continue;
        }

        // Проверка столкновения с снарядами
        for (let j = bullets.length - 1; j >= 0; j--) {
            if (isColliding(bullets[j], enemies[i])) {
                // enemyHitSound.play();
                enemies.splice(i, 1);
                bullets.splice(j, 1);
                score += 10;
                renderHUD();
                break;
            }
        }
    }

    requestAnimationFrame(gameLoop);
}

// Функция отправки очков на сервер
function sendScore() {
    fetch('/score', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            game_type: 'shooter',
            score: score
        })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Score sent:', data);
    })
    .catch((error) => {
        console.error('Error:', error);
    });
}

// Инициализация счёта
highScore = localStorage.getItem('highScore') || 0;
renderHUD();

// Запуск игры
gameLoop();
