package config

import (
    "encoding/json"
    "os"
)

type Forward struct{ Listen, Target string }

/* ---------- server ---------- */
type ServerConf struct {
    Listen   string    `json:"listen"`
    PSK      string    `json:"psk"`
    Forwards []Forward `json:"forwards"`
}

/* ---------- client ---------- */
type Peer struct {
    Name       string `json:"name"`
    Server     string `json:"server"`
    PSK        string `json:"psk"`
    SocksLocal string `json:"socks_local"`
    MaxStreams int    `json:"max_streams"`
}

type ClientConf struct {
    Peers []Peer `json:"peers"`
}

/* ---------- helpers ---------- */
func Load(path string, v any) error {
    b, err := os.ReadFile(path)
    if err != nil { return err }
    return json.Unmarshal(b, v)
}

func Save(path string, v any) error {
    b,_ := json.MarshalIndent(v,"","  ")
    return os.WriteFile(path,b,0600)
}
