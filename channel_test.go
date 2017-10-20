package main

import (
	"testing"
)

/*
func TestCA(t *testing.T) {

	dp, err := NewCounterDataProvider()
	if err != nil {
		t.Error(err)
	}

	b := []byte{}

	fmt.Println(string(dp.GetCSVHeader(b)))
	fmt.Println(string(dp.GetData(b)))
	fmt.Println(string(dp.GetData(b)))
	fmt.Println(string(dp.GetData(b)))

}

//
//
//
func TestCB(t *testing.T) {

	dp, err := NewCounterDataProvider()
	if err != nil {
		t.Error(err)
	}

	var c *Channel
	c, err = NewChannel(1, 1, dp)

	fmt.Print(string(c.PrepareRecord()))
	c.Success()
	fmt.Print(string(c.PrepareRecord()))
	c.Success()
	fmt.Print(string(c.PrepareRecord()))

}

func TestCC(t *testing.T) {

	dp, err := NewCounterDataProvider()
	if err != nil {
		t.Error(err)
	}

	var c *Channel
	c, err = NewChannel(1, 1, dp)

	fmt.Print(string(c.GetCSVHeader()))

}
*/

func TestCD(t *testing.T) {

	s := NewChannelService()
	err := s.InitChannels() //s.PrepareChannels()
	if err != nil {
		t.Error(err)
	}

}
