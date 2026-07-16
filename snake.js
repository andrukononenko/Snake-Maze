// snake.js
const readline = require('readline');
const { stdin, stdout } = process;

class SnakeGame {
    constructor(rows = 21, cols = 21) {
        this.rows = rows;
        this.cols = cols;
        this.maze = [];
        this.snake = [];
        this.dir = {r: 0, c: 1};
        this.food = null;
        this.score = 0;
        this.speed = 150; // ms
        this.gameOver = false;
        this.running = true;
        this.generateMaze();
        this.spawnSnake();
        this.spawnFood();
        this.setupInput();
        this.loop();
    }

    generateMaze() {
        this.maze = Array.from({length: this.rows}, () => Array(this.cols).fill('#'));
        // DFS carve
        const stack = [{r: 1, c: 1}];
        this.maze[1][1] = ' ';
        const dirs = [[-2,0], [2,0], [0,-2], [0,2]];
        while (stack.length) {
            const cur = stack[stack.length-1];
            const neighbours = [];
            dirs.forEach((d, idx) => {
                const nr = cur.r + d[0], nc = cur.c + d[1];
                if (nr > 0 && nr < this.rows-1 && nc > 0 && nc < this.cols-1 && this.maze[nr][nc] === '#') {
                    neighbours.push(idx);
                }
            });
            if (neighbours.length) {
                const idx = neighbours[Math.floor(Math.random() * neighbours.length)];
                const dr = dirs[idx][0], dc = dirs[idx][1];
                const nr = cur.r + dr, nc = cur.c + dc;
                this.maze[cur.r + dr/2][cur.c + dc/2] = ' ';
                this.maze[nr][nc] = ' ';
                stack.push({r: nr, c: nc});
            } else {
                stack.pop();
            }
        }
        // borders
        for (let i = 0; i < this.rows; i++) {
            this.maze[i][0] = '#';
            this.maze[i][this.cols-1] = '#';
        }
        for (let j = 0; j < this.cols; j++) {
            this.maze[0][j] = '#';
            this.maze[this.rows-1][j] = '#';
        }
    }

    spawnSnake() {
        let center = {r: Math.floor(this.rows/2), c: Math.floor(this.cols/2)};
        if (this.maze[center.r][center.c] !== ' ') {
            for (let r = 1; r < this.rows-1; r++) {
                for (let c = 1; c < this.cols-1; c++) {
                    if (this.maze[r][c] === ' ') {
                        center = {r, c};
                        break;
                    }
                }
                if (center.r !== Math.floor(this.rows/2)) break;
            }
        }
        this.snake = [center];
        this.dir = {r: 0, c: 1};
    }

    spawnFood() {
        const empty = [];
        for (let r = 1; r < this.rows-1; r++) {
            for (let c = 1; c < this.cols-1; c++) {
                if (this.maze[r][c] === ' ') {
                    const occupied = this.snake.some(p => p.r === r && p.c === c);
                    if (!occupied) empty.push({r, c});
                }
            }
        }
        if (empty.length) {
            this.food = empty[Math.floor(Math.random() * empty.length)];
        } else {
            this.gameOver = true;
        }
    }

    update() {
        if (this.gameOver) return;
        const head = this.snake[0];
        const newHead = {r: head.r + this.dir.r, c: head.c + this.dir.c};
        // wall
        if (newHead.r <= 0 || newHead.r >= this.rows-1 || newHead.c <= 0 || newHead.c >= this.cols-1 || this.maze[newHead.r][newHead.c] === '#') {
            this.gameOver = true;
            return;
        }
        // self
        if (this.snake.some(p => p.r === newHead.r && p.c === newHead.c)) {
            this.gameOver = true;
            return;
        }
        this.snake.unshift(newHead);
        if (newHead.r === this.food.r && newHead.c === this.food.c) {
            this.score++;
            this.speed = Math.max(50, this.speed * 0.95);
            this.spawnFood();
            if (this.gameOver) return;
        } else {
            this.snake.pop();
        }
    }

    draw() {
        console.clear();
        console.log(`Score: ${this.score}  Speed: ${Math.floor(1000/this.speed)}`);
        const display = this.maze.map(row => [...row]);
        for (const p of this.snake) {
            display[p.r][p.c] = 'O';
        }
        if (this.food) {
            display[this.food.r][this.food.c] = '*';
        }
        display[this.snake[0].r][this.snake[0].c] = 'X';
        for (const row of display) {
            console.log(row.join(''));
        }
        if (this.gameOver) {
            console.log('GAME OVER! Press R to restart, Q to quit.');
        }
    }

    setupInput() {
        readline.emitKeypressEvents(stdin);
        stdin.setRawMode(true);
        stdin.on('keypress', (str, key) => {
            if (key.ctrl && key.name === 'c') process.exit();
            const name = key.name;
            if (name === 'q') { this.running = false; process.exit(); }
            if (name === 'r' && this.gameOver) {
                this.__init__(this.rows, this.cols);
                return;
            }
            if (!this.gameOver) {
                if (name === 'up' || name === 'w') {
                    if (this.dir.r !== 1) this.dir = {r: -1, c: 0};
                } else if (name === 'down' || name === 's') {
                    if (this.dir.r !== -1) this.dir = {r: 1, c: 0};
                } else if (name === 'left' || name === 'a') {
                    if (this.dir.c !== 1) this.dir = {r: 0, c: -1};
                } else if (name === 'right' || name === 'd') {
                    if (this.dir.c !== -1) this.dir = {r: 0, c: 1};
                }
            }
        });
    }

    __init__(rows, cols) {
        Object.assign(this, new SnakeGame(rows, cols));
    }

    loop() {
        if (!this.running) return;
        this.update();
        this.draw();
        setTimeout(() => this.loop(), this.speed);
    }
}

new SnakeGame(21, 21);
