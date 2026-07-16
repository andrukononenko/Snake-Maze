# snake.rb
require 'io/console'
require 'timeout'

class SnakeGame
  attr_reader :rows, :cols

  def initialize(rows=21, cols=21)
    @rows, @cols = rows, cols
    @maze = Array.new(rows) { Array.new(cols, '#') }
    @snake = []
    @dir = [0, 1]
    @food = nil
    @score = 0
    @speed = 0.15
    @game_over = false
    @running = true
    generate_maze
    spawn_snake
    spawn_food
    Thread.new { input_loop }
    game_loop
  end

  def generate_maze
    stack = [[1, 1]]
    @maze[1][1] = ' '
    dirs = [[-2,0], [2,0], [0,-2], [0,2]]
    until stack.empty?
      r, c = stack.last
      neigh = []
      dirs.each_with_index do |d, idx|
        nr, nc = r + d[0], c + d[1]
        if nr > 0 && nr < @rows-1 && nc > 0 && nc < @cols-1 && @maze[nr][nc] == '#'
          neigh << idx
        end
      end
      if neigh.any?
        idx = neigh.sample
        dr, dc = dirs[idx]
        nr, nc = r + dr, c + dc
        @maze[r + dr/2][c + dc/2] = ' '
        @maze[nr][nc] = ' '
        stack.push([nr, nc])
      else
        stack.pop
      end
    end
    (0...@rows).each { |i| @maze[i][0] = '#'; @maze[i][@cols-1] = '#' }
    (0...@cols).each { |j| @maze[0][j] = '#'; @maze[@rows-1][j] = '#' }
  end

  def spawn_snake
    cr, cc = @rows/2, @cols/2
    if @maze[cr][cc] != ' '
      (1...@rows-1).each do |r|
        (1...@cols-1).each do |c|
          if @maze[r][c] == ' '
            cr, cc = r, c
            break
          end
        end
        break if @maze[cr][cc] == ' '
      end
    end
    @snake = [[cr, cc]]
    @dir = [0, 1]
  end

  def spawn_food
    empty = []
    (1...@rows-1).each do |r|
      (1...@cols-1).each do |c|
        if @maze[r][c] == ' ' && !@snake.include?([r, c])
          empty << [r, c]
        end
      end
    end
    if empty.any?
      @food = empty.sample
    else
      @game_over = true
    end
  end

  def update
    return if @game_over
    head = @snake.first
    nr, nc = head[0] + @dir[0], head[1] + @dir[1]
    if nr <= 0 || nr >= @rows-1 || nc <= 0 || nc >= @cols-1 || @maze[nr][nc] == '#'
      @game_over = true
      return
    end
    if @snake.include?([nr, nc])
      @game_over = true
      return
    end
    @snake.unshift([nr, nc])
    if [nr, nc] == @food
      @score += 1
      @speed = [0.05, @speed - 0.005].max
      spawn_food
      return if @game_over
    else
      @snake.pop
    end
  end

  def draw
    system('clear') || system('cls')
    puts "Score: #{@score}  Speed: #{ (1/@speed).to_i }"
    display = @maze.map(&:dup)
    @snake.each { |r,c| display[r][c] = 'O' }
    display[@food[0]][@food[1]] = '*' if @food
    display[@snake.first[0]][@snake.first[1]] = 'X'
    display.each { |row| puts row.join }
    puts "GAME OVER! Press R to restart, Q to quit." if @game_over
  end

  def input_loop
    while @running
      char = STDIN.getch
      case char
      when 'q', 'Q' then @running = false; exit
      when 'r', 'R'
        if @game_over
          # restart: reinitialize
          initialize(@rows, @cols)
          return
        end
      when 'w' then @dir = [-1, 0] if @dir != [1, 0]
      when 's' then @dir = [1, 0] if @dir != [-1, 0]
      when 'a' then @dir = [0, -1] if @dir != [0, 1]
      when 'd' then @dir = [0, 1] if @dir != [0, -1]
      when "\e" # arrow keys
        c2 = STDIN.read_nonblock(2) rescue nil
        if c2 == '[A' then @dir = [-1, 0] if @dir != [1, 0]
        elsif c2 == '[B' then @dir = [1, 0] if @dir != [-1, 0]
        elsif c2 == '[C' then @dir = [0, 1] if @dir != [0, -1]
        elsif c2 == '[D' then @dir = [0, -1] if @dir != [0, 1]
        end
      end
    end
  end

  def game_loop
    while @running
      update
      draw
      sleep @speed
    end
  end
end

SnakeGame.new
