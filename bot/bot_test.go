package bot

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"gitlab.com/etke.cc/buscarron/mocks"
)

type BotSuite struct {
	suite.Suite
	lp  *mocks.Linkpearl
	bot *Bot
}

var ctxMatcher = mock.MatchedBy(func(_ context.Context) bool { return true })

func (s *BotSuite) SetupTest() {
	s.T().Helper()
	s.lp = &mocks.Linkpearl{}
	s.bot = New(s.lp)
}

func (s *BotSuite) TearDownTest() {
	s.lp.AssertExpectations(s.T())
}

func (s *BotSuite) TestNew() {
	bot := New(s.lp)

	s.IsType(&Bot{}, bot)
}

func (s *BotSuite) TestError_NoLinkpearl() {
	bot := New(nil)

	bot.Error(context.TODO(), id.RoomID("!doesnt:matt.er"), "msg %s", "arg")
}

func (s *BotSuite) TestError() {
	roomID := id.RoomID("!doesnt:matt.er")
	s.lp.On("Send", ctxMatcher, roomID, &event.Content{
		Parsed: &event.MessageEventContent{
			MsgType:  event.MsgNotice,
			Body:     "ERROR: msg arg",
			Mentions: &event.Mentions{},
		},
	}).Return(id.EventID("$doesnt:matt.er"), nil).Once()

	s.bot.Error(context.TODO(), id.RoomID("!doesnt:matt.er"), "msg %s", "arg")
}

func (s *BotSuite) TestSend() {
	roomID := id.RoomID("!doesnt:matt.er")
	s.lp.On("Send", ctxMatcher, roomID, &event.Content{
		Parsed: &event.MessageEventContent{
			MsgType:  event.MsgNotice,
			Body:     "msg",
			Mentions: &event.Mentions{},
		},
	}).Return(id.EventID("$doesnt:matt.er"), nil).Once()

	s.bot.Send(context.TODO(), id.RoomID("!doesnt:matt.er"), "msg", nil)
}

func (s *BotSuite) TestSend_Error() {
	roomID := id.RoomID("!doesnt:matt.er")
	s.lp.On("Send", ctxMatcher, roomID, &event.Content{
		Parsed: &event.MessageEventContent{
			MsgType:  event.MsgNotice,
			Body:     "msg",
			Mentions: &event.Mentions{},
		},
	}).Return(id.EventID("$doesnt:matt.er"), errors.New("test")).Once()

	s.bot.Send(context.TODO(), id.RoomID("!doesnt:matt.er"), "msg", nil)
}

func (s *BotSuite) TestSendFile() {
	roomID := id.RoomID("!doesnt:matt.er")
	req := &mautrix.ReqUploadMedia{
		FileName:      "test.txt",
		ContentBytes:  []byte("test"),
		ContentLength: int64(len([]byte("test"))),
		ContentType:   "text/plain",
	}
	s.lp.On("SendFile", ctxMatcher, roomID, req, event.MsgFile).Return(nil).Once()

	s.bot.SendFile(context.TODO(), id.RoomID("!doesnt:matt.er"), req)
}

func (s *BotSuite) TestSendFile_Error() {
	roomID := id.RoomID("!doesnt:matt.er")
	req := &mautrix.ReqUploadMedia{
		FileName:      "test.txt",
		ContentBytes:  []byte("test"),
		ContentLength: int64(len([]byte("test"))),
		ContentType:   "text/plain",
	}
	s.lp.On("SendFile", ctxMatcher, roomID, req, event.MsgFile).Return(errors.New("test")).Once()
	s.lp.On("Send", ctxMatcher, roomID, &event.Content{
		Parsed: &event.MessageEventContent{
			MsgType:  "m.notice",
			Body:     "ERROR: cannot upload file: test",
			Mentions: &event.Mentions{},
		},
	}).Return(id.EventID("$doesnt:matt.er"), nil).Once()
	s.bot.SendFile(context.TODO(), id.RoomID("!doesnt:matt.er"), req)
}

func (s *BotSuite) TestStart() {
	s.lp.On("Start", ctxMatcher).Return(nil).Once()

	s.bot.Start()
}

func (s *BotSuite) TestStart_Error() {
	s.lp.On("Start", ctxMatcher).Return(errors.New("test")).Once()

	s.Error(s.bot.Start())
}

func (s *BotSuite) TestStop() {
	s.lp.On("Stop", ctxMatcher).Once()

	s.bot.Stop()
}

func TestBotSuite(t *testing.T) {
	suite.Run(t, new(BotSuite))
}
