package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Channel struct {
	Num       byte `yaml:"num"`
	id        []byte
	PK        int64 `yaml:"pk"`
	csvHeader []byte
	b         []byte
	B         *[]byte
	dp        IDataProvider
	t         *IntToByteArr
}

const CH_CSV_HEADER = "S#H;ID;Timestamp;PK;"

//
//
//
func NewChannel(num byte, pk int64, dp IDataProvider) (c *Channel, err error) {

	if num < 1 || num > 96 {
		s := fmt.Sprintf("channel num should between '0..96' but was %d", num)
		return nil, errors.New(s)
	}

	return &Channel{
		csvHeader: []byte(string(CH_CSV_HEADER)),
		Num:       num,
		id:        []byte{'S', '#', (num / 10) + 48, (num % 10) + 48, ';'},
		b:         make([]byte, 0, 256),
		dp:        dp,
		PK:        pk,
		t:         NewIntToByteArr(),
	}, nil

}

//
// GetCSVHeader
//
func (c *Channel) GetCSVHeader() []byte {
	c.b = c.b[:0]
	c.b = append(c.b, c.csvHeader...)
	return c.dp.GetCSVHeader(c.b)
}

//
//
//
func (c *Channel) header() {

	c.b = c.b[:0]
	var t = c.t

	// id
	c.b = append(c.b, c.id...)

	// time stamp
	t.in = time.Now().UnixNano()
	c.b = append(c.b, intToByteArr(t)...)
	c.b = append(c.b, ';')

	// primary key
	t.in = c.PK
	c.b = append(c.b, intToByteArr(t)...)
	c.b = append(c.b, ';')

}

//
// Succes increments the value of the primary key
//
func (c *Channel) Success() {
	c.PK++
}

//
//
//
func (c *Channel) Bytes() []byte {
	return c.b
}

//
//
//
func (c *Channel) GetRecord() []byte {
	c.header()
	c.b = c.dp.GetData(c.b)
	c.b = append(c.b, 10)
	return c.b
}

//
// ChannelService
//
type ChannelService struct {
	HowMuch  int
	Channels []*Channel
}

//
//
//
func NewChannelService() *ChannelService {
	return &ChannelService{
		Channels: []*Channel{},
	}
}

//
//
//
func (s *ChannelService) InitChannels() error {

	type Ch struct {
		Num byte  `yaml:"num"`
		PK  int64 `yaml:"pk"`
	}

	type Chler struct {
		HowMuch  int
		Channels []*Ch
	}

	chler := &Chler{}

	err := ReadYML(chler, "app.yml")
	if err != nil {
		return nil
	}

	s.HowMuch = chler.HowMuch

	var dp IDataProvider
	var l = len(chler.Channels)
	for i, ch := range chler.Channels {

		if s.HowMuch < l {
			if i == s.HowMuch {
				break
			}
		}

		dp, err = NewCounterDataProvider()
		if err != nil {
			return err
		}

		c, err := NewChannel(byte(ch.Num), ch.PK, dp)
		if err != nil {
			return err
		}

		s.Channels = append(s.Channels, c)

	}

	for i := len(s.Channels) + 1; i <= s.HowMuch; i++ {

		dp, err = NewCounterDataProvider()
		if err != nil {
			return err
		}

		c, err := NewChannel(byte(i), 1, dp)
		if err != nil {
			return err
		}

		s.Channels = append(s.Channels, c)

	}

	return nil
}

//
//
//
func (s *ChannelService) PrepareChannels() error {

	for i := 1; i < 5; i++ {
		c, err := NewChannel(byte(i), 1, nil)
		if err != nil {
			return err
		}
		fmt.Println(string(c.Num))

		s.Channels = append(s.Channels, c)
	}

	return writeYML(s, "app.yml")

}

//
//
//
func writeYML(v interface{}, filename string) error {

	var (
		b   []byte
		err error
	)

	if b, err = yaml.Marshal(v); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, b, 0644)
}

//
//
//
func ReadYML(v interface{}, filename string) error {
	in, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(in, v)
}
