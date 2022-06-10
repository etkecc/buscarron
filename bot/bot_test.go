package bot

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/etke.cc/go/logger"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/mocks"
)

type BotSuite struct {
	suite.Suite
	lp  *mocks.Linkpearl
	bot *Bot
}

func (s *BotSuite) SetupTest() {
	s.T().Helper()
	log := logger.New("matrix.", "TRACE")
	s.lp = &mocks.Linkpearl{}
	s.bot = New(s.lp, log)
}

func (s *BotSuite) TearDownTest() {
	s.lp.AssertExpectations(s.T())
}

func (s *BotSuite) TestNew() {
	bot := New(s.lp, logger.New("matrix.", "TRACE"))

	s.IsType(&Bot{}, bot)
}

func (s *BotSuite) TestError_NoLinkpearl() {
	bot := New(nil, logger.New("matrix.", "TRACE"))

	bot.Error(id.RoomID("!doesnt:matt.er"), "msg %s", "arg")
}

func (s *BotSuite) TestError() {
	roomID := id.RoomID("!doesnt:matt.er")
	s.lp.On("Send", roomID, &event.MessageEventContent{
		MsgType: event.MsgNotice,
		Body:    "ERROR: msg arg",
	}).Return(id.EventID("$doesnt:matt.er"), nil).Once()

	s.bot.Error(id.RoomID("!doesnt:matt.er"), "msg %s", "arg")
}

func (s *BotSuite) TestSend() {
	roomID := id.RoomID("!doesnt:matt.er")
	s.lp.On("Send", roomID, &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    "msg",
	}).Return(id.EventID("$doesnt:matt.er"), nil).Once()

	s.bot.Send(id.RoomID("!doesnt:matt.er"), "msg")
}

func (s *BotSuite) TestSend_Error() {
	roomID := id.RoomID("!doesnt:matt.er")
	s.lp.On("Send", roomID, &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    "msg",
	}).Return(id.EventID("$doesnt:matt.er"), errors.New("test")).Once()
	s.lp.On("Send", roomID, &event.MessageEventContent{
		MsgType: event.MsgNotice,
		Body:    "ERROR: cannot send message: test",
	}).Return(id.EventID("$doesnt:matt.er"), nil).Once()

	s.bot.Send(id.RoomID("!doesnt:matt.er"), "msg")
}

func (s *BotSuite) TestStart() {
	s.lp.On("Start").Return(nil).Once()

	s.bot.Start()
}

func (s *BotSuite) TestStart_Error() {
	s.lp.On("Start").Return(errors.New("test")).Once()

	fn := func() {
		s.bot.Start()
	}

	s.Panics(fn)
}

func (s *BotSuite) TestStop() {
	s.lp.On("Stop").Once()

	s.bot.Stop()
}

func TestBotSuite(t *testing.T) {
	suite.Run(t, new(BotSuite))
}
