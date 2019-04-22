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
		"spell.wav",
		"doorOpen_1.ogg",
		"doorClose_1.ogg",
		"interface1.wav",
		"interface2.wav",
		"explodemini.wav",
		"explode.wav",
		"bookFlip2.ogg",
		"coin3.wav",
		"piano.wav",
		"orchestra.wav",
		"voices/male_standard_1.ogg",
		"voices/female_standard_1.ogg",
		"monsters/rat.wav",
		"monsters/spider.wav",
	}
	for _, name := range sounds {
		sound, err := mix.LoadWAV("ui2d/assets/sounds/" + name)
		if err != nil {
			panic(err)
		}
		sound.Volume(48)
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

func (mp *MusicPlayer) PlaySound(name string, volume int) {
	mp.Sounds[name].Volume(volume)
	mp.Sounds[name].Play(-1, 0)
}

func (mp *MusicPlayer) On(e *game.Event) {
	switch e.Action {
	case game.ActionChangeLevel:
		mp.PlayMusicForLevel(e.Payload["levelType"])
	// case game.ActionWalk:
	// 	mp.PlaySound("footstep00.ogg")
	case game.ActionOpenDoor:
		mp.PlaySound("doorOpen_1.ogg", 32)
	case game.ActionCloseDoor:
		mp.PlaySound("doorClose_1.ogg", 32)
	case game.ActionMenuOpen:
		mp.PlaySound("interface2.wav", 64)
	case game.ActionMenuClose:
		mp.PlaySound("interface2.wav", 64)
	case game.ActionMenuSelect:
		mp.PlaySound("interface1.wav", 48)
	case game.ActionMenuConfirm:
		mp.PlaySound("interface2.wav", 64)
	case game.ActionAttack:
		mp.PlaySound("footstep08.ogg", 48)
	case game.ActionPower:
		mp.PlayPowerSound(e)
	case game.ActionHurt:
		mp.PlaySound("footstep01.ogg", 48)
	case game.ActionExplode:
		mp.PlayExplosion(e)
	case game.ActionTalk:
		mp.PlayVoice(e)
	case game.ActionTake:
		mp.PlaySound("interface1.wav", 48)
	case game.ActionReadBook:
		mp.PlaySound("bookFlip2.ogg", 64)
	case game.ActionRoar:
		mp.PlayMonsterRoar(e)
	case game.ActionTakeGold:
		mp.PlaySound("coin3.wav", 48)
	case game.ActionQuestFinished:
		mp.PlaySound("orchestra.wav", 64)
	case game.ActionCharacteristicUp:
		mp.PlaySound("piano.wav", 64)
	}
}

func (mp *MusicPlayer) PlayMusicForLevel(levelType string) {
	switch levelType {
	case game.LevelTypeGrotto:
		mp.ChangeMusic("dark_fallout.ogg")
	case game.LevelTypeOutdoor:
		mp.ChangeMusic("forest.mp3")
	case game.LevelTypeCity:
		mp.ChangeMusic("warped.mp3")
	}
}

func (mp *MusicPlayer) PlayPowerSound(e *game.Event) {
	volume := 48
	sound := "magic1.wav"
	typ, exists := e.Payload["type"]
	if exists {
		switch typ {
		case game.PowerEnergyBall:
			sound = "magic1.wav"
			volume = 32
		case game.PowerInvocation:
			sound = "spell.wav"
			volume = 64
		}
	}
	mp.PlaySound(sound, volume)
}

func (mp *MusicPlayer) PlayExplosion(e *game.Event) {
	volume := 48
	sound := "explodemini.wav"
	size, exists := e.Payload["size"]
	if exists {
		switch size {
		case game.ExplosionSizeSmall:
			sound = "footstep08.ogg"
		case game.ExplosionSizeMedium:
			sound = "explodemini.wav"
		case game.ExplosionSizeLarge:
			sound = "explode.wav"
			volume = 64
		}
	}
	mp.PlaySound(sound, volume)
}

func (mp *MusicPlayer) PlayVoice(e *game.Event) {
	sound := "interface1.wav"
	voice, exists := e.Payload["voice"]
	if exists {
		switch voice {
		case game.VoiceMaleStandard:
			sound = "voices/male_standard_1.ogg"
		case game.VoiceFemaleStandard:
			sound = "voices/female_standard_1.ogg"
		}
	}
	mp.PlaySound(sound, 48)
}

func (mp *MusicPlayer) PlayMonsterRoar(e *game.Event) {
	sound := "interface1.wav"
	monsterName, exists := e.Payload["monster"]
	if exists {
		monsterRune := game.Tile(monsterName[0])
		switch monsterRune {
		case game.Spider:
			sound = "monsters/spider.wav"
		case game.Rat:
			sound = "monsters/rat.wav"
		}
	}
	mp.PlaySound(sound, 64)
}
