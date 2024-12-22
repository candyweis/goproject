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

function gameLoop() {
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
      snakeLength = 4;
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
    foodX = Math.floor(Math.random() * tileCount);
    foodY = Math.floor(Math.random() * tileCount);
  }
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
setInterval(gameLoop, 100);
