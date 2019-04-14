package ui2d

import (
	"github.com/veandco/go-sdl2/mix"
	"thelema/game"
)

type MusicPlayer struct {
	Musics           map[string]*mix.Music
	Sounds           map[string]*mix.Chunk
	CurrentMusicName string
}

func NewMusicPlayer() *MusicPlayer {

	mix.Init(mix.INIT_MP3 | mix.INIT_OGG)

	mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)

	mp := &MusicPlayer{}
	mp.Musics = make(map[string]*mix.Music)
	mp.Sounds = make(map[string]*mix.Chunk)
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

func (mp *MusicPlayer) LoadSounds() {
	var sounds []string = []string{
		"footstep00.ogg",
		"footstep01.ogg",
		"footstep08.ogg",
		"magic1.wav",
		"doorOpen_1.ogg",
		"doorClose_1.ogg",
		"interface1.wav",
		"interface2.wav",
		"explodemini.wav",
		// TODO
	}
	for _, name := range sounds {
		sound, err := mix.LoadWAV("ui2d/assets/sounds/" + name)
		if err != nil {
			panic(err)
		}
		sound.Volume(64)
		mp.Sounds[name] = sound
	}
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

func (mp *MusicPlayer) PlaySound(name string) {
	mp.Sounds[name].Play(-1, 0)
}

func (mp *MusicPlayer) On(e *game.Event) {
	switch e.Action {
	case game.ActionChangeLevel:
		mp.PlayMusicForLevel(e.Payload["levelType"])
	// case game.ActionWalk:
	// 	mp.PlaySound("footstep00.ogg")
	case game.ActionOpenDoor:
		mp.PlaySound("doorOpen_1.ogg")
	case game.ActionCloseDoor:
		mp.PlaySound("doorClose_1.ogg")
	case game.ActionMenuOpen:
		mp.PlaySound("interface2.wav")
	case game.ActionMenuClose:
		mp.PlaySound("interface2.wav")
	case game.ActionMenuSelect:
		mp.PlaySound("interface1.wav")
	case game.ActionMenuConfirm:
		mp.PlaySound("interface2.wav")
	case game.ActionAttack:
		mp.PlaySound("footstep08.ogg")
	case game.ActionPower:
		mp.PlaySound("magic1.wav")
	case game.ActionHurt:
		mp.PlaySound("footstep01.ogg")
	case game.ActionExplode:
		mp.PlaySound("explodemini.wav")
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
