package sender

import (
	"github.com/7amiro0/home_work_golang/hw12_13_14_15_calendar/internal/app"
	"github.com/streadway/amqp"
)

type Sender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	Logger  Logger
}

type Logger interface {
	app.Logger
}

type OptionsReadQueue struct {
	AutoAck, Exclusive, NoLocal, NoWait bool
}

func New(logger Logger) *Sender {
	return &Sender{
		Logger: logger,
	}
}

func (s *Sender) Start(name, url string, optRead OptionsReadQueue) (delivery <-chan amqp.Delivery, err error) {
	if s.conn, err = amqp.Dial(url); err != nil {
		s.Logger.Error("[ERR] Error connect to amqp: ", err)
		return delivery, err
	}

	if s.channel, err = s.conn.Channel(); err != nil {
		s.Logger.Error("[ERR] Error connect to channel: ", err)
		return delivery, err
	}

	return s.channel.Consume(name, "", optRead.AutoAck, optRead.Exclusive, optRead.NoLocal, optRead.NoWait, nil)
}

func (s *Sender) Stop() error {
	if err := s.channel.Close(); err != nil {
		s.Logger.Error("[ERR] error close to channel: ", err)
		return err
	}

	return s.conn.Close()
}
