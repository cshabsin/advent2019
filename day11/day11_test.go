package main

import "testing"

func TestRobot(t *testing.T) {
	robot := paintBot{posX: 1, posY: 1}
	if got, _ := robot.Read(); got != 0 {
		t.Errorf("Expected 0")
	}
	if robot.phase {
		t.Errorf("Expected phase false")
	}
	robot.Write(1)
	if !robot.phase {
		t.Errorf("Expected phase true")
	}
	if got, _ := robot.Read(); got != 1 {
		t.Errorf("Expected 1")
	}
	robot.Write(0)
	if robot.phase {
		t.Errorf("Expected phase false")
	}
	if got, _ := robot.Read(); got != 0 {
		t.Errorf("Expected 0")
	}
	if robot.posX != 0 {
		t.Errorf("posX expected 0, got %d (posY = %d)", robot.posX, robot.posY)
	}
}
