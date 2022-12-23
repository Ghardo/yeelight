package yeelight

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

type Yeelight struct {
	YLID 	   int32       		`json:"id"`
	Address    string      		`json:"address"`
	Persistent bool        		`json:"persistent",default0:"false"`
	Conn       net.Conn		 	`json:"-"`
	Timeout    time.Duration
}

type Command struct {
	ID     int32       `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Response struct {
	ID     int32       `json:"id"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

type Options struct {
	Smooth     int `default0:"200"`
}

type RGB struct {
	Value int
}

func (rgb *RGB) ToHex() string {
	return fmt.Sprintf("%06x", rgb.Value)
}

func (c *Command) GenerateID() {
	if c.ID == 0 {
		r := rand.NewSource(time.Now().UnixNano())
		c.ID = rand.New(r).Int31()
	}
}

func (c *Command) ToJson() ([]byte, error) {
	cmdJson, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	return cmdJson, nil
}

func (r *Response) FromJson(data []byte) error {
	return json.Unmarshal(data, &r)
}

func (yl *Yeelight) Connect() (err error) {
	if yl.Timeout == 0 {
		yl.Timeout = 3 * time.Second
	}

	yl.Conn, err = net.DialTimeout("tcp", yl.Address, yl.Timeout)
	if err != nil {
		return err
	}

	return nil
}

func (yl *Yeelight) SendCommand(c Command) (r Response, err error) {
	c.GenerateID()
	if err = yl.Connect(); err != nil {
		return
	}

	if !yl.Persistent {
		defer yl.Conn.Close()
	}

	cmdJSON, err := c.ToJson()
	if err != nil {
		return r, err
	}

	if _, err := fmt.Fprintf(yl.Conn, "%s\r\n", cmdJSON); err != nil {
		return r, err
	}

	responseBuf := bufio.NewReader(yl.Conn)
	response, err := responseBuf.ReadSlice('\n')
	if err != nil {
		return r, err
	}

	r.FromJson(response)
	return r, err
}

func (yl *Yeelight) GetProperties(names []string) (r Response, err error) {
	c := Command{
		Method: "get_prop",
		Params: names,
	}

	return yl.SendCommand(c)
}

func (yl *Yeelight) GetProperty(name string) (r Response, err error) {
	c := Command{
		Method: "get_prop",
		Params: []interface{}{name},
	}

	return yl.SendCommand(c)
}

// Wrapper Methods

func (yl *Yeelight) SetHexColor(color string, options Options) (err error) {
	color = strings.Replace(color, "#", "", -1)
	n, err := strconv.ParseUint(color, 16, 64)
	if err != nil {
		return
	}

	c := Command{
		Method: "set_rgb",
		Params: []interface{}{n, "smooth", options.Smooth},
	}

	_, err = yl.SendCommand(c)
	if err != nil {
		return
	}

	return nil
}

func (yl *Yeelight) GetHexColor() (h string, err error) {
	r, err := yl.GetProperty("rgb")
	if err != nil {
		return h, err
	}

	value, err := strconv.Atoi(r.Result.([]interface{})[0].(string))

	rgb := RGB{Value: value}
	if err != nil {
		return h, err
	}
	h = rgb.ToHex()
	return h, nil
}

func (yl *Yeelight) SetBright(value int8, options Options) (err error) {
	c := Command{
		Method: "set_bright",
		Params: []interface{}{value, "smooth", options.Smooth},
	}

	_, err = yl.SendCommand(c)
	if err != nil {
		return
	}

	return nil
}

func (yl *Yeelight) SetColorTemperature(value int16, options Options) (err error) {
	c := Command{
		Method: "set_ct_abx",
		Params: []interface{}{value, "smooth", options.Smooth},
	}

	_, err = yl.SendCommand(c)
	if err != nil {
		return
	}

	return nil
}

func (yl *Yeelight) GetBright() (value int8, err error) {
	r, err := yl.GetProperty("bright")
	if err != nil {
		return value, err
	}
	v, err := strconv.ParseInt(r.Result.([]interface{})[0].(string), 10, 8)
	if err != nil {
		return value, err
	}

	value = int8(v)

	return value, nil
}

func (yl *Yeelight) SetOn(options Options) (err error) {
	c := Command{
		Method: "set_power",
		Params: []interface{}{"on", "smooth", options.Smooth},
	}

	_, err = yl.SendCommand(c)
	if err != nil {
		return
	}

	return nil
}

func (yl *Yeelight) SetOff(options Options) (err error) {
	c := Command{
		Method: "set_power",
		Params: []interface{}{"off", "smooth", options.Smooth},
	}

	_, err = yl.SendCommand(c)
	if err != nil {
		return
	}

	return nil
}

func (yl *Yeelight) Toggle() (err error) {
	c := Command{
		Method: "toggle",
		Params: []interface{}{},
	}

	_, err = yl.SendCommand(c)
	if err != nil {
		return
	}

	return nil
}

func (yl *Yeelight) IsOn() (b bool, err error) {
	r, err := yl.GetProperty("power")
	if err != nil {
		return b, err
	}

	b = (r.Result.([]interface{})[0].(string) == "on")

	return b, err
}
