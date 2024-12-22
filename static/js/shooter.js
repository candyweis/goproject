const c = document.getElementById("shooterCanvas");
const cx = c.getContext("2d");

let playerX = c.width / 2;
let playerY = c.height / 2;

let bullets = [];
let enemies = [];

// Создаём несколько врагов
for (let i = 0; i < 5; i++) {
  enemies.push({
    x: Math.random() * c.width,
    y: Math.random() * c.height,
    size: 15
  });
}

function draw() {
  // Заливка фона
  cx.fillStyle = "black";
  cx.fillRect(0, 0, c.width, c.height);

  // Игрок
  cx.fillStyle = "white";
  cx.fillRect(playerX - 10, playerY - 10, 20, 20);

  // Пули
  cx.fillStyle = "yellow";
  bullets.forEach((b) => {
    b.x += b.vx;
    b.y += b.vy;
    cx.fillRect(b.x - 2, b.y - 2, 4, 4);
  });

  // Враги
  cx.fillStyle = "red";
  enemies.forEach((e) => {
    cx.fillRect(e.x - e.size / 2, e.y - e.size / 2, e.size, e.size);
  });

  // Удаляем пули, которые вышли за границы
  bullets = bullets.filter((b) => b.x >= 0 && b.x <= c.width && b.y >= 0 && b.y <= c.height);

  requestAnimationFrame(draw);
}

document.addEventListener("keydown", (e) => {
  const speed = 5;
  switch (e.key) {
    case "ArrowLeft":
      playerX -= speed;
      break;
    case "ArrowRight":
      playerX += speed;
      break;
    case "ArrowUp":
      playerY -= speed;
      break;
    case "ArrowDown":
      playerY += speed;
      break;
    case " ":
      // Стреляем пулей вверх
      bullets.push({ x: playerX, y: playerY, vx: 0, vy: -5 });
      break;
  }
});

draw();
