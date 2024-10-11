//go:build ignore_unused
// +build ignore_unused

package main

import (
	"fmt"
	_ "log"
	"os"

	_ "github.com/davecgh/go-spew/spew"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/yaml.v2"
	"strings"
	"text/template"
)

type Config struct {
	Server ServerConfig   `yaml:"server"`
	Client []ClientConfig `yaml:"client"`
}

type ServerConfig struct {
	Address   string   `yaml:"address"`
	Key       string   `yaml:"key"`
	PublicKey string   `yaml:"publickey"`
	Host      string   `yaml:"host"`
	Port      string   `yaml:"port"`
	DNS       []string `yaml:"DNS"`
	AllowedIP []string `yaml:"AllowedIP"`
}

type ClientConfig struct {
	Name      string   `yaml:"name"`
	Address   string   `yaml:"address"`
	Key       string   `yaml:"key"`
	PublicKey string   `yaml:"publickey"`
	AllowedIP []string `yaml:"AllowedIP"`
}

// Define the template for the configuration file
const serverConfigTemplate = `
# created on ${DATE}

[Interface]
Address = {{.Server.Address}}/24
ListenPort = {{ .Server.Port }}
PrivateKey = {{.Server.Key}}
PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE

{{range .Client}}
# {{.Name}}
[Peer]
# private key: {{.Key}}
PublicKey = {{.PublicKey}}
AllowedIPs = {{.Address}}/32
{{end}}


`

// Define the template for the configuration file
const clientConfigTemplate = `
# created on ${DATE}

{{- with index .Client 0 }}
[Interface]

Address = {{ .Address }}/24
ListenPort = {{ $.Server.Port }}
PrivateKey = {{ .Key }}
DNS = {{ index $.Server.DNS 1 }}, {{ index $.Server.DNS 1 }}
{{- end }}

[Peer]

# client public key
PublicKey = {{ .Server.PublicKey }}
Endpoint = {{ .Server.Host }}:{{.Server.Port}}

{{- with index .Client 0 }}
  {{- if .AllowedIP }}
AllowedIPs = {{ join .AllowedIP ", " }}
  {{- else }}
AllowedIPs = 0.0.0.0/0, ::/0
  {{- end }}
{{- end }}



`

// AllowedIPs = 0.0.0.0/0, ::/0

func createServerConfig(data Config) error {

	//spew.Dump(data)

	tmpl, err := template.New("config").Parse(serverConfigTemplate)
	if err != nil {
		return err
	}

	path := configdir + "server.wg0.conf"
	fmt.Printf("creating %s\n", path)
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}

func createClientConfig(data Config) error {

	//spew.Dump(data)
	//fmt.Printf("file: %s\n", data.Client[0].Name)

	// Create a new template and register the custom function
	tmpl := template.New("config").Funcs(template.FuncMap{
		"join": Join,
	})

	tmpl, err := tmpl.Parse(clientConfigTemplate)
	if err != nil {
		fmt.Printf("failed to parse the template\n")
		return err
	}

	//fmt.Printf("Server Key: %s\n", cfg.Server.Key)
	path := configdir + data.Client[0].Name + ".conf"
	fmt.Printf("creating %s\n", path)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return err
	}

	return nil
}

func generatePublicKey(privateKey string) (string, error) {
	privKey, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return "", err
	}

	pubKey := privKey.PublicKey()
	return pubKey.String(), nil
}

// Function to add PublicKey to every ClientConfig
func addPublicKeys(config *Config) error {

	publicKey, err := generatePublicKey(config.Server.Key)
	config.Server.PublicKey = publicKey

	for i := range config.Client {
		publicKey, err = generatePublicKey(config.Client[i].Key)
		if err != nil {
			return fmt.Errorf("failed to generate public key for client %s: %v", config.Client[i].Name, err)
		}
		config.Client[i].PublicKey = publicKey
	}
	return nil
}

// Join is a custom function that joins a slice of strings with a separator.
func Join(elements []string, separator string) string {
	// fmt.Printf("using %s\n", separator)
	return strings.Join(elements, separator)
}

var configdir string

func main() {

	configdir = os.Args[1] + "/"
	//fmt.Printf("using %s\n", configdir)

	var cfg Config

	{
		f, err := os.ReadFile("wg.yaml")
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(f, &cfg)
		if err != nil {
			panic(err)
		}

		addPublicKeys(&cfg)

	}

	//spew.Dump(cfg)
	createServerConfig(cfg)

	for i := range cfg.Client {
		//fmt.Printf(" creating conf for peer %d %s\n", i, cfg.Client[i].Name)
		c := Config{Server: cfg.Server, Client: []ClientConfig{cfg.Client[i]}}
		createClientConfig(c)
	}

}
