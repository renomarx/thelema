Thelema project
===============

- Coded in [GO](https://golang.org/)
- Using [SDL2](https://github.com/veandco/go-sdl2)


Install & Run (dev mode)
-------------

- Install GO: https://golang.org/
- Follow the steps to install sdl2 & sdl2-mixer: https://github.com/veandco/go-sdl2
- Clone the repo into your go/src folder
- In the repo folder:
   - `cp config.json.template config.json`
   - `go run main.go` (should be long the first time)


Graphics
--------

- Main tileset is `ui2d/assets/tiles.png`
- Tiles mapping for this tileset is `ui2d/assets/atlas-index.txt`
- Maps are auto-loaded from `game/maps/**` folders
- You can see the tiles names in `game/tiles.go`

- Player textures are in `ui2d/assets/player`

- Pnjs textures are in `ui2d/assets/pnjs`


Audio
-----

- Musics are auto-loaded from `ui2d/musics` folder
- Sounds are auto-loaded from `ui2d/sounds/**` folders


Dialogs
-----

- Pnj dialogs are auto-loaded from `game/pnjs/**/*.json` files


Books
-----

- Books are auto-loaded from `game/books` folder
