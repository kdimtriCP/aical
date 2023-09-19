package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/pkg/openai"
	"strings"
	"time"
)

type ChatUseCase struct {
	log    *log.Helper
	client *openai.Client
	fr     *openai.Registry
	gr     GoogleRepo
	cr     CalendarRepo
}

// NewChatUseCase .
func NewChatUseCase(cfg *conf.OpenAI, logger log.Logger, gr GoogleRepo, cr CalendarRepo) *ChatUseCase {
	return &ChatUseCase{
		log:    log.NewHelper(logger),
		client: openai.NewClient(cfg.Api.Key, cfg.Api.Model),
		fr:     openai.NewRegistry(),
		gr:     gr,
		cr:     cr,
	}
}

func (uc *ChatUseCase) UserChat(ctx context.Context, question string) (string, error) {
	messageContext := make([]openai.ChatCompletionMessage, 0)
	messageContext = append(messageContext, openai.ChatCompletionMessage{
		Role: "system",
		Content: "You are AI assistant that helps user to manage his calendar and answer to user questions about his calendar. " +
			//			"You can use the following functions to create, update and delete events: create_event, update_event, delete_event." +
			"You can use the following functions to get information about user events and calendars: list_events, list_user_calendars. " +
			"You can also use the following functions to get information about current time: current_time.",
	})
	messageContext = append(messageContext, openai.ChatCompletionMessage{
		Role:    "user",
		Content: question,
	})

	uc.fr.Register(currentTimeFunctionDescription().Name, currentTimeFunctionDescription(), uc.currentTimeFunction)
	//	uc.fr.Register(createEventFunctionDescription().Name, createEventFunctionDescription(), uc.createEventFunction)
	//	uc.fr.Register(updateEventFunctionDescription().Name, updateEventFunctionDescription(), uc.updateEventFunction)
	//	uc.fr.Register(deleteEventFunctionDescription().Name, deleteEventFunctionDescription(), uc.deleteEventFunction)
	uc.fr.Register(listEventsFunctionDescription().Name, listEventsFunctionDescription(), uc.listEventsFunction)
	uc.fr.Register(listUserCalendarsFunctionDescription().Name, listUserCalendarsFunctionDescription(), uc.listUserCalendarsFunction)

	request := &openai.ChatCompletionRequest{
		Messages:  messageContext,
		Functions: uc.fr.Descriptions(),
	}
	var answer string
	for {
		response, err := uc.client.DoRequest(ctx, request)
		if err != nil {
			return err.Error(), err
		}
		if response.Choices[0].FinishReason == "stop" {
			uc.log.Debugf("\nQuestion: %s\nAnswer: %s",
				question,
				response.Choices[0].Message.Content,
			)
			answer = response.Choices[0].Message.Content
			break
		}
		if response.Choices[0].FinishReason == "function_call" {
			request.AddFunctionCall(
				response.Choices[0].Message.FunctionCall.Name,
				response.Choices[0].Message.FunctionCall.Arguments,
				uc.fr.Execute(
					ctx,
					response.Choices[0].Message.FunctionCall.Name,
					response.Choices[0].Message.FunctionCall.Arguments,
				),
			)
		}
	}
	return answer, nil
}

func (uc *ChatUseCase) createEventFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("createEventFunction: %s", arguments)
	args := &struct {
		CalendarID string    `json:"calendar_id"`
		GoogleID   string    `json:"google_id,omitempty"`
		Title      string    `json:"title"`
		Location   string    `json:"location,omitempty"`
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}
	calId, err := uuid.Parse(args.CalendarID)
	if err != nil {
		calId = uuid.Nil
	}
	event := &Event{
		CalendarID: calId,
		GoogleID:   args.GoogleID,
		Summary:    args.Title,
		Location:   args.Location,
		StartTime:  args.StartTime,
		EndTime:    args.EndTime,
	}
	// Create event in ggogle calendar
	token := GetToken(ctx)
	if token == nil {
		return "token not found in context"
	}
	e, err := uc.gr.CreateCalendarEvent(ctx, token, event, args.CalendarID)
	if err != nil {
		return err.Error()
	}
	return e.String()
}

func (uc *ChatUseCase) updateEventFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("updateEventFunction: %s", arguments)
	args := &struct {
		CalendarID string    `json:"calendar_id"`
		GoogleID   string    `json:"google_id"`
		Title      string    `json:"title,omitempty"`
		Location   string    `json:"location,omitempty"`
		StartTime  time.Time `json:"start_time,omitempty"`
		EndTime    time.Time `json:"end_time,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}
	calId, err := uuid.Parse(args.CalendarID)
	if err != nil {
		calId = uuid.Nil
	}
	event := &Event{
		CalendarID: calId,
		GoogleID:   args.GoogleID,
		Summary:    args.Title,
		Location:   args.Location,
		StartTime:  args.StartTime,
		EndTime:    args.EndTime,
	}

	token := GetToken(ctx)
	if token == nil {
		return "token not found in context"
	}
	e, err := uc.gr.UpdateCalendarEvent(ctx, token, event, args.CalendarID)
	if err != nil {
		return err.Error()
	}
	return e.String()
}

func (uc *ChatUseCase) deleteEventFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("deleteEventFunction: %s", arguments)
	args := &struct {
		CalendarID string `json:"calendar_id"`
		GoogleID   string `json:"google_id"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}

	calId, err := uuid.Parse(args.CalendarID)
	if err != nil {
		calId = uuid.Nil
	}

	event := &Event{
		CalendarID: calId,
		GoogleID:   args.GoogleID,
	}
	token := GetToken(ctx)
	if token == nil {
		return "token not found in context"
	}
	err = uc.gr.DeleteCalendarEvent(ctx, token, event, args.CalendarID)
	if err != nil {
		return err.Error()
	}
	return "Event deleted"
}

func (uc *ChatUseCase) currentTimeFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("currentTimeFunction: %s", arguments)
	return time.Now().Format(time.RFC3339)
}

func (uc *ChatUseCase) listEventsFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("listEventsFunction: %s", arguments)
	args := &struct {
		GoogleID  string `json:"google_id"`
		StartTime string `json:"start_time,omitempty"`
		EndTime   string `json:"end_time,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}
	if args.StartTime == "" {
		args.StartTime = time.Now().AddDate(0, 0, -int(time.Now().Weekday())+1).Format(time.RFC3339) // this week
	}

	if args.EndTime == "" {
		args.EndTime = time.Now().AddDate(0, 0, 14-int(time.Now().Weekday())).Format(time.RFC3339) // next week
	}
	token := GetToken(ctx)
	events, err := uc.gr.ListCalendarEvents(ctx, token, args.GoogleID, &GoogleListEventsOption{
		TimeMin: args.StartTime,
		TimeMax: args.EndTime,
	})
	if err != nil {
		return err.Error()
	}
	eventsString := make([]string, len(events))
	for i, event := range events {
		eventsString[i] = event.String()
	}
	uc.log.Debugf("listEventsFunction: %s", strings.Join(eventsString, ","))
	return "[" + strings.Join(eventsString, ",") + "]"
}

// listUserCalendarsFuntion is a function that returns a list of user calendars
func (uc *ChatUseCase) listUserCalendarsFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("listUserCalendarsFunction: %s", arguments)
	token := GetToken(ctx)
	calendars, err := uc.gr.ListUserCalendars(ctx, token)
	if err != nil {
		return err.Error()
	}
	calendarsString := make([]string, len(calendars))
	for i, calendar := range calendars {
		c, err := uc.cr.Get(ctx, calendar)
		if err != nil {
			return err.Error()
		}
		calendarsString[i] = c.String()
	}
	return "[" + strings.Join(calendarsString, ",") + "]"
}
