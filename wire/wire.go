//+build wireinject

package wire


import "github.com/google/wire"

func InitMission(name string) Mission {
	wire.Build(NewMonster, NewPlayer, NewMission)
	return Mission{}
}