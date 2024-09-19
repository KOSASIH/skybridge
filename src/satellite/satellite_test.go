package satellite

import (
	"testing"
)

func TestSatellite(t *testing.T) {
	satellite := NewSatellite("sat1")
	groundStation1 := NewGroundStation("gs1")
	groundStation2 := NewGroundStation("gs2")
	transceiver1 := NewTransceiver("tx1")
	transceiver2 := NewTransceiver("tx2")

	satellite.AddGroundStation(groundStation1)
	satellite.AddGroundStation(groundStation2)
	satellite.AddTransceiver(transceiver1)
	satellite.AddTransceiver(transceiver2)

	err := satellite.CommunicateWithGroundStation(groundStation1, []byte("Hello, world!"))
	if err != nil {
		t.Errorf("Error communicating with ground station: %v", err)
	}

	err = satellite.CommunicateWithMultipleGroundStations([]*GroundStation{groundStation1, groundStation2}, []byte("Hello, world!"))
	if err != nil {
		t.Errorf("Error communicating with multiple ground stations: %v", err)
	}

	go satellite.StartOrbit()

	time.Sleep(30 * time.Second)
}
