package server

import (
	"fmt"
	"strconv"
	"strings"
)

type gameMode struct {
	str string
	num int
}

const (
	SurvivalNum = iota + 1
	CreativeNum
	AdventureNum
	SpectatorNum
)

const (
	SurvivalStr  = "Survival"
	CreativeStr  = "Creative"
	AdventureStr = "Adventure"
	SpectatorStr = "Spectator"
)

type Motd struct {
	edition         string
	motd            string
	protocolVersion int
	versionName     string
	playerCount     int
	maxPlayerCount  int
	serverUniqueId  string
	levelName       string
	gameMode        gameMode
	port4           int
	port6           int
}

func NewMotd() *Motd {
	return &Motd{
		edition:         "MCPE",
		motd:            "Dedicated Server",
		protocolVersion: 1,
		versionName:     "1.0.0",
		playerCount:     -1,
		maxPlayerCount:  -1,
		serverUniqueId:  "123",
		levelName:       "Bedrock Level",
		gameMode: gameMode{
			str: SurvivalStr,
			num: SurvivalNum,
		},
		port4: 19132,
		port6: 19132,
	}
}

func NewMotdFromString(str string) *Motd {
	splitMotd := strings.Split(str, ";")

	edition := func() string {
		if len(splitMotd) < 1 {
			return "MCPE"
		}
		return splitMotd[0]
	}()

	motd := func() string {
		if len(splitMotd) < 2 {
			return "Dedicated Server"
		}
		return splitMotd[1]
	}()

	protocolVersion := func() int {
		if len(splitMotd) < 3 {
			return 1
		}

		ver, err := strconv.Atoi(splitMotd[2])
		if err != nil {
			return 1
		}

		return ver
	}()

	versionName := func() string {
		if len(splitMotd) < 4 {
			return "1.0.0"
		}
		return splitMotd[3]
	}()

	playerCount := func() int {
		if len(splitMotd) < 5 {
			return 0
		}

		count, err := strconv.Atoi(splitMotd[4])
		if err != nil {
			return 0
		}

		return count
	}()

	maxPlayerCount := func() int {
		if len(splitMotd) < 6 {
			return 1
		}

		count, err := strconv.Atoi(splitMotd[5])
		if err != nil {
			return 1
		}

		return count
	}()

	serverUniqueId := func() string {
		if len(splitMotd) < 7 {
			return ""
		}
		return splitMotd[6]
	}()

	levelName := func() string {
		if len(splitMotd) < 8 {
			return ""
		}
		return splitMotd[7]
	}()

	gameModeStr := func() string {
		if len(splitMotd) < 9 {
			return string(CreativeStr)
		}
		return splitMotd[8]
	}()

	gameModeNum := func() int {
		if len(splitMotd) < 10 {
			return int(SurvivalNum)
		}

		mode, err := strconv.Atoi(splitMotd[9])
		if err != nil {
			return int(SurvivalNum)
		}

		return mode
	}()

	port4 := func() int {
		if len(splitMotd) < 11 {
			return 19132
		}

		port, err := strconv.Atoi(splitMotd[10])
		if err != nil {
			return 19132
		}

		return port
	}()

	port6 := func() int {
		if len(splitMotd) < 12 {
			return 19132
		}

		port, err := strconv.Atoi(splitMotd[11])
		if err != nil {
			return 19132
		}

		return port
	}()

	return &Motd{
		edition:         edition,
		motd:            motd,
		protocolVersion: protocolVersion,
		versionName:     versionName,
		playerCount:     playerCount,
		maxPlayerCount:  maxPlayerCount,
		serverUniqueId:  serverUniqueId,
		levelName:       levelName,
		gameMode: gameMode{
			str: gameModeStr,
			num: gameModeNum,
		},
		port4: port4,
		port6: port6,
	}
}

func (m Motd) String() string {
	return fmt.Sprintf(
		"%s;%s;%d;%s;%d;%d;%s;%s;%s;%d;%d;%d",
		m.edition,
		m.motd,
		m.protocolVersion,
		m.versionName,
		m.playerCount,
		m.maxPlayerCount,
		m.serverUniqueId,
		m.levelName,
		m.gameMode.str,
		m.gameMode.num,
		m.port4,
		m.port6,
	)
}
