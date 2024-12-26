// static/js/snake.js
const canvas = document.getElementById("snakeCanvas");
const ctx = canvas.getContext("2d");

const gridSize = 20;
const tileCount = canvas.width / gridSize;

let snakePosX = 10;
let snakePosY = 10;
let snakeBody = [];
let snakeLength = 4;

let velocityX = 0;
let velocityY = 0;

let foodX = 5;
let foodY = 5;

let score = 0;
let gameOver = false;

// Отрисовка HUD
function drawHUD() {
  ctx.fillStyle = "white";
  ctx.font = "20px Arial";
  ctx.fillText(`Очки: ${score}`, 10, 20);
}

function gameLoop() {
  if (gameOver) {
    // Отображение сообщения об окончании игры
    ctx.fillStyle = "rgba(0, 0, 0, 0.5)";
    ctx.fillRect(0, canvas.height / 2 - 50, canvas.width, 100);
    ctx.fillStyle = "white";
    ctx.font = "30px Arial";
    ctx.textAlign = "center";
    ctx.fillText("Игра окончена!", canvas.width / 2, canvas.height / 2);
    ctx.fillText(`Ваш счёт: ${score}`, canvas.width / 2, canvas.height / 2 + 40);
    // Отправка очков на сервер
    sendScore();
    return;
  }

  snakePosX += velocityX;
  snakePosY += velocityY;

  // Зацикливание по краям
  if (snakePosX < 0) snakePosX = tileCount - 1;
  if (snakePosX > tileCount - 1) snakePosX = 0;
  if (snakePosY < 0) snakePosY = tileCount - 1;
  if (snakePosY > tileCount - 1) snakePosY = 0;

  // Рисуем фон
  ctx.fillStyle = "black";
  ctx.fillRect(0, 0, canvas.width, canvas.height);

  // Рисуем змейку
  ctx.fillStyle = "lime";
  for (let i = 0; i < snakeBody.length; i++) {
    let part = snakeBody[i];
    ctx.fillRect(part.x * gridSize, part.y * gridSize, gridSize - 2, gridSize - 2);

    // Столкновение с собой
    if (part.x === snakePosX && part.y === snakePosY) {
      gameOver = true;
      return;
    }
  }

  snakeBody.push({ x: snakePosX, y: snakePosY });
  while (snakeBody.length > snakeLength) {
    snakeBody.shift();
  }

  // Рисуем еду
  ctx.fillStyle = "red";
  ctx.fillRect(foodX * gridSize, foodY * gridSize, gridSize - 2, gridSize - 2);

  // Если съели еду
  if (snakePosX === foodX && snakePosY === foodY) {
    snakeLength++;
    score += 10; // Увеличение счёта при поедании еды
    foodX = Math.floor(Math.random() * tileCount);
    foodY = Math.floor(Math.random() * tileCount);
  }

  // Отрисовка HUD
  drawHUD();

  setTimeout(gameLoop, 100);
}

function keyPress(evt) {
  switch (evt.keyCode) {
    case 37: // left
      if (velocityX !== 1) {
        velocityX = -1; velocityY = 0;
      }
      break;
    case 38: // up
      if (velocityY !== 1) {
        velocityX = 0; velocityY = -1;
      }
      break;
    case 39: // right
      if (velocityX !== -1) {
        velocityX = 1; velocityY = 0;
      }
      break;
    case 40: // down
      if (velocityY !== -1) {
        velocityX = 0; velocityY = 1;
      }
      break;
  }
}

document.addEventListener("keydown", keyPress);

// Функция отправки очков на сервер
function sendScore() {
  fetch('/score', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      game_type: 'snake',
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

gameLoop();
