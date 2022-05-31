package wire

import (
	"testing"
)


func TestWire(t *testing.T)  {
	mission := InitMission("test")
	mission.Start()
}