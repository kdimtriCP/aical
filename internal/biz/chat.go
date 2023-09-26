package biz

import (
	"context"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/log"
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
	er     EventRepo
}

// NewChatUseCase .
func NewChatUseCase(cfg *conf.OpenAI, logger log.Logger, gr GoogleRepo, cr CalendarRepo, er EventRepo) *ChatUseCase {
	return &ChatUseCase{
		log:    log.NewHelper(logger),
		client: openai.NewClient(cfg.Api.Key, cfg.Api.Model),
		fr:     openai.NewRegistry(),
		gr:     gr,
		cr:     cr,
		er:     er,
	}
}

// systemMessage returns a system message for assistant
func systemMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{
		Role: "system",
		Content: "You are an AI assistant that helps the user manage his calendar with smart event scheduling. " +
			"If a user asks to create an event, first use list_events to analyze the user's existing events for the specified day. " +
			"If there are no events or there are free slots, suggest the best times for the new event. If the day is fully booked, notify the user. " +
			"Use create_event to finalize the creation of the event." +
			"Use current_time to get the current time." +
			"Use adjust_date to adjust the current date by a number of days. " +
			"For example to get tomorrow's date use current_time to get today's date and use adjust_date(1) to get tomorrow.",
	}
}

func (uc *ChatUseCase) UserChat(ctx context.Context, question string) (string, error) {
	messageContext := make([]openai.ChatCompletionMessage, 0)
	messageContext = append(messageContext, systemMessage())
	messageContext = append(messageContext, openai.ChatCompletionMessage{
		Role:    "user",
		Content: question,
	})

	uc.fr.Register(currentTimeFunctionDescription().Name, currentTimeFunctionDescription(), currentTimeFunction)
	uc.fr.Register(adjustDateFunctionDescription().Name, adjustDateFunctionDescription(), adjustDateFunction)
	uc.fr.Register(createEventFunctionDescription().Name, createEventFunctionDescription(), uc.createEventFunction)
	uc.fr.Register(updateEventFunctionDescription().Name, updateEventFunctionDescription(), uc.updateEventFunction)
	uc.fr.Register(deleteEventFunctionDescription().Name, deleteEventFunctionDescription(), uc.deleteEventFunction)
	uc.fr.Register(listEventsFunctionDescription().Name, listEventsFunctionDescription(), uc.listEventsFunction)
	uc.fr.Register(listUserCalendarsFunctionDescription().Name, listUserCalendarsFunctionDescription(), uc.listUserCalendarsFunction)

	request := &openai.ChatCompletionRequest{
		Messages:  messageContext,
		Functions: uc.fr.Descriptions(),
	}
	uc.log.Debugf("Chat request: \n%v", request)
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
		GoogleCalendarID string    `json:"google_calendar_id"`
		Title            string    `json:"title"`
		Location         string    `json:"location,omitempty"`
		StartTime        time.Time `json:"start_time"`
		EndTime          time.Time `json:"end_time"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}

	event := &Event{
		Summary:   args.Title,
		Location:  args.Location,
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	}
	// Create event in ggogle calendar
	token := GetToken(ctx)
	if token == nil {
		return "token not found in context"
	}
	if args.GoogleCalendarID == "" {
		args.GoogleCalendarID = "primary"
	}
	e, err := uc.gr.CreateCalendarEvent(ctx, token, event, args.GoogleCalendarID)
	if err != nil {
		return err.Error()
	}
	return e.String()
}

func (uc *ChatUseCase) updateEventFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("updateEventFunction: %s", arguments)
	args := &struct {
		GoogleCalendarID string    `json:"google_calendar_id"`
		GoogleID         string    `json:"google_event_id"`
		Title            string    `json:"title,omitempty"`
		Location         string    `json:"location,omitempty"`
		StartTime        time.Time `json:"start_time,omitempty"`
		EndTime          time.Time `json:"end_time,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}
	event := &Event{
		GoogleID:  args.GoogleID,
		Summary:   args.Title,
		Location:  args.Location,
		StartTime: args.StartTime,
		EndTime:   args.EndTime,
	}

	token := GetToken(ctx)
	if token == nil {
		return "token not found in context"
	}

	if args.GoogleCalendarID == "" {
		args.GoogleCalendarID = "primary"
	}
	e, err := uc.gr.UpdateCalendarEvent(ctx, token, event, args.GoogleCalendarID)
	if err != nil {
		return err.Error()
	}
	return e.String()
}

func (uc *ChatUseCase) deleteEventFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("deleteEventFunction: %s", arguments)
	token := GetToken(ctx)
	if token == nil {
		return "token not found in context"
	}
	args := &struct {
		GoogleCalendarID string `json:"google_calendar_id"`
		GoogleEventID    string `json:"google_event_id"`
	}{}
	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}

	event := &Event{
		GoogleID: args.GoogleEventID,
	}
	if args.GoogleCalendarID == "" {
		args.GoogleCalendarID = "primary"
	}
	err = uc.gr.DeleteCalendarEvent(ctx, token, event, args.GoogleCalendarID)
	if err != nil {
		return err.Error()
	}
	return "Event deleted"
}

func (uc *ChatUseCase) listEventsFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("listEventsFunction: %s", arguments)
	args := &struct {
		GoogleCalendarID string `json:"google_calendar_id"`
		StartTime        string `json:"start_time,omitempty"`
		EndTime          string `json:"end_time,omitempty"`
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
	if token == nil {
		return "error: token not found in context"
	}

	if args.GoogleCalendarID == "" {
		args.GoogleCalendarID = "primary"
	}

	events, err := uc.gr.ListCalendarEvents(ctx, token, args.GoogleCalendarID, &GoogleListEventsOption{
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
	return "[" + strings.Join(eventsString, ",") + "]"
}

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
