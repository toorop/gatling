package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/smtp"
	//"time"
)

// A gatling bullet is... a smtp transaction ;)
type gatling struct {
	Target   string
	MailFrom string
	RcptTo   string
	BodySize uint64
	Bullets  uint64
}

// Fire creates a new SMTP transaction and "send" a mail
func (g *gatling) Fire(cDone *chan (bool)) {
	err := g.fire()
	if err != nil {
		log.Println(err)
	}
	*cDone <- true
}

func (g *gatling) fire() error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(g.Target)
	if err != nil {
		return err
	}
	for i := uint64(0); i < g.Bullets; i++ {
		log.Printf("On envoit une buillet %d/%d", i, g.Bullets)

		// Set the sender and recipient first
		if err := c.Mail(g.MailFrom); err != nil {
			c.Quit()
			return err
		}
		if err := c.Rcpt(g.RcptTo); err != nil {
			c.Quit()
			return err
		}

		// Send the email body.
		wc, err := c.Data()
		if err != nil {
			c.Quit()
			return err
		}

		_, err = fmt.Fprintf(wc, getRandBody(g.BodySize))
		if err != nil {
			c.Quit()
			return err
		}
		err = wc.Close()
		if err != nil {
			c.Quit()
			return err
		}

		err = c.Reset()
		if err != nil {
			c.Quit()
			return err
		}
	}
	c.Quit()
	return nil

}

func getRandBody(size uint64) string {
	var bytes = make([]byte, size)
	rand.Read(bytes)
	return string(bytes)
}
