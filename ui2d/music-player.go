package ui2d

import (
	"github.com/veandco/go-sdl2/mix"
	"thelema/game"
)

type MusicPlayer struct {
	Musics           map[string]*mix.Music
	CurrentMusicName string
}

func NewMusicPlayer() *MusicPlayer {

	mix.Init(mix.INIT_MP3 | mix.INIT_OGG)

	mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)

	mp := &MusicPlayer{}
	mp.Musics = make(map[string]*mix.Music)
	return mp
}

func (mp *MusicPlayer) LoadMusics() {
	var musics []string = []string{
		"doomed.mp3",
		"flags.mp3",
		"great_mission.mp3",
		"spacetime.mp3",
		"twists.mp3",
		"waking_the_devil.mp3",
		"warped.mp3",
		"dark_fallout.ogg",
		"forest.mp3",
	}
	for _, name := range musics {
		music, err := mix.LoadMUS("ui2d/assets/musics/" + name)
		if err != nil {
			panic(err)
		}
		mp.Musics[name] = music
	}
	mp.CurrentMusicName = "great_mission.mp3"
}

func (mp *MusicPlayer) PlayMusic() {
	mp.Musics[mp.CurrentMusicName].Play(-1)
}

func (mp *MusicPlayer) StopMusic() {
	mix.HaltMusic()
}

func (mp *MusicPlayer) ChangeMusic(musicName string) {
	mix.HaltMusic()
	mp.CurrentMusicName = musicName
	mp.PlayMusic()
}

func (mp *MusicPlayer) On(e *game.Event) {
	switch e.Type {
	case game.PlayerEventsType:
		switch e.Action {
		case game.ActionChangeLevel:
			mp.PlayMusicForLevel(e.Payload["levelType"])
		}
	}
}

func (mp *MusicPlayer) PlayMusicForLevel(levelType string) {
	switch levelType {
	case game.LevelTypeGrotto:
		mp.ChangeMusic("dark_fallout.ogg")
	case game.LevelTypeOutdoor:
		mp.ChangeMusic("forest.mp3")
	}
}
