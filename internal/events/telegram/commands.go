package telegram

import (
	"errors"
	"log"
	"net/url"
	"read-advise-links-bot/internal/clients/telegram"
	e "read-advise-links-bot/internal/lib"
	"read-advise-links-bot/internal/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, userName)

	sendMsg := NewMessageSender(chatID, p.tg)

	// add page
	if isAddCmd(text) {
		// TODO: AddPage(
		return p.savePage(chatID, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendStart(chatID)
	default:
		return sendMsg(msgUnknownCommand)
	}
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) (err error) {
	defer func() {
		err = e.WrapIfErr("can not do command: sendMsg page", err)
	}()
	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	sendMsg := NewMessageSender(chatID, p.tg)

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can not do command: send random message", err) }()

	page, err := p.storage.PickRandom(userName)

	sendMsg := NewMessageSender(chatID, p.tg)

	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	sendMsg := NewMessageSender(chatID, p.tg)
	return sendMsg(msgHelp)
}

func (p *Processor) sendStart(chatID int) error {
	sendMsg := NewMessageSender(chatID, p.tg)
	return sendMsg(msgStart)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	// TODO make able for 'ya.ru' links, now it's only 'https://ya.ru'
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
