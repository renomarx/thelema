package ui2d

import (
	"thelema/game"

	"github.com/veandco/go-sdl2/mix"
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
	musics := game.LoadFilenames("ui2d/assets/musics/")
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
	sounds := game.LoadFilenames("ui2d/assets/sounds/")
	voices := game.LoadFilenames("ui2d/assets/sounds/voices/")
	for _, voice := range voices {
		sounds = append(sounds, "voices/"+voice)
	}
	monsters := game.LoadFilenames("ui2d/assets/sounds/monsters/")
	for _, monster := range monsters {
		sounds = append(sounds, "monsters/"+monster)
	}
	uis := game.LoadFilenames("ui2d/assets/sounds/ui/")
	for _, ui := range uis {
		sounds = append(sounds, "ui/"+ui)
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

func (mp *MusicPlayer) PlayMusic(volume int) {
	mix.VolumeMusic(volume)
	mp.Musics[mp.CurrentMusicName].Play(-1)
}

func (mp *MusicPlayer) StopMusic() {
	mix.HaltMusic()
}

func (mp *MusicPlayer) ChangeMusic(musicName string, volume int) {
	mix.HaltMusic()
	mp.CurrentMusicName = musicName
	mp.PlayMusic(volume)
}

func (mp *MusicPlayer) PlaySound(name string, volume int) {
	mp.Sounds[name].Volume(volume)
	mp.Sounds[name].Play(-1, 0)
}

func (mp *MusicPlayer) On(e *game.Event) {
	switch e.Action {
	case game.ActionChangeLevel:
		mp.PlayMusicForLevel(e.Payload["levelName"])
	// case game.ActionWalk:
	// 	mp.PlaySound("footstep00.ogg")
	case game.ActionOpenDoor:
		mp.PlaySound("doorOpen_1.ogg", 32)
	case game.ActionCloseDoor:
		mp.PlaySound("doorClose_1.ogg", 32)
	case game.ActionMenuOpen:
		mp.PlaySound("ui/1bip.wav", 48)
	case game.ActionMenuClose:
		mp.PlaySound("ui/1bip.wav", 48)
	case game.ActionMenuSelect:
		mp.PlaySound("ui/1bip.wav", 48)
	case game.ActionMenuConfirm:
		mp.PlaySound("ui/2bip.wav", 48)
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
		mp.PlaySound("ui/1bip.wav", 48)
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
	case game.ActionFight:
		mp.ChangeMusic("stress.wav", 64)
	case game.ActionStopFight:
		mp.PlayMusicForLevel(e.Payload["levelName"])
	}
}

func (mp *MusicPlayer) PlayMusicForLevel(levelName string) {
	switch levelName {
	case "abigail_crypt", "prison_underground":
		mp.ChangeMusic("dark_fallout.ogg", 32)
	case "world":
		mp.ChangeMusic("forest.mp3", 64)
	case "neoroma":
		mp.ChangeMusic("warped.mp3", 32)
	case "arcanea":
		mp.ChangeMusic("quick-melody.wav", 32)
	case "dresde":
		mp.ChangeMusic("ketamine.wav", 32)
	default:
		mp.StopMusic()
	}
}

func (mp *MusicPlayer) PlayPowerSound(e *game.Event) {
	volume := 48
	sound := "magic1.wav"
	typ, exists := e.Payload["type"]
	if exists {
		switch typ {
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
	sound := "voices/male_standard_1.ogg"
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
	sound := "monsters/rat.wav"
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
