package apnian

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/electronicpanopticon/gobrick"
	"github.com/mitchellh/go-homedir"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/token"
	"github.com/spf13/viper"
	"log"
)

type ApnianConfigurer struct {
	ConfigName string
	Root string
}

type Apnian struct {
	P8KeyName  string
	Topic      string
	APNSKeyID  string
	TeamID     string
	Configurer *ApnianConfigurer
	Client     *apns2.Client

}

// New returns an Apnian filed with the values in its config file.
// Locations it looks for are:
//		.
//		..
//		$GOPATH/config
//		$HOME
func New(configName string) (*Apnian, error) {
	ac := ApnianConfigurer{configName, gobrick.GetGOPATH()}
	return ac.getApnian()
}

// AuthKeyPath returns the path to the ECDSA private key specified in the Apnian file.
func (ac Apnian) AuthKeyPath() string {
	rel := fmt.Sprintf("keys/%s", ac.P8KeyName)
	return fmt.Sprintf("%s/%s", ac.Configurer.Root, rel)
}

// AuthKey returns the ECDSA private key specified in the Apnian file.
func (ac Apnian) AuthKey() (*ecdsa.PrivateKey, error) {
	return token.AuthKeyFromFile(ac.AuthKeyPath())
}

// Token represents an Apple Provider Authentication Token (JSON Web Token) configured
// with the values from the Apnian file.
func (ac Apnian) Token() (*token.Token, error) {
	authKey, err := ac.AuthKey()
	if err != nil {
		return &token.Token{}, err
	}
	return &token.Token{
		AuthKey:  authKey,
		KeyID:    ac.APNSKeyID,
		TeamID:   ac.TeamID,
	}, nil
}

func (ac Apnian) Notification(deviceID string, payload *APS) *apns2.Notification {
	notification := &apns2.Notification{}
	notification.DeviceToken = deviceID
	notification.Topic = ac.Topic
	notification.Payload = payload.ToJsonBytes()
	return notification
}

// getApnian returns an Apnian from the configured Viper instance.
func (ac ApnianConfigurer) getApnian() (*Apnian, error) {
	ac.configureViper()

	var c Apnian
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&c); err != nil {
		return nil, err
	}
	c.Configurer = &ac
	return &c, nil
}

// configureViper
func (ac ApnianConfigurer) configureViper() {
	viper.SetConfigName(ac.ConfigName)
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	home, err := homedir.Dir()
	if err == nil {
		viper.AddConfigPath(home)
	} else {
		log.Println("unable to get homedir")
	}
	viper.AddConfigPath(ac.Root + "/config")
}