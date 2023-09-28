package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/internal/conf"
	"github.com/kdimtricp/aical/pkg/openai"
	"strings"
	"time"
)

type OpenAIUseCase struct {
	log    *log.Helper
	client *openai.Client
	fr     *openai.Registry
	gr     GoogleRepo
}

// NewOpenAIUseCase .
func NewOpenAIUseCase(cfg *conf.OpenAI, logger log.Logger, gr GoogleRepo) *OpenAIUseCase {
	return &OpenAIUseCase{
		log:    log.NewHelper(logger),
		client: openai.NewClient(cfg.Api.Key, cfg.Api.Model),
		fr:     openai.NewRegistry(),
		gr:     gr,
	}
}

func (uc *OpenAIUseCase) buildOpenAIChangeEventQuery(changes []*EventHistory) (string, error) {
	if len(changes) == 0 {
		return "", fmt.Errorf("no changes provided")
	}

	// Convert each change into a string representation
	changeDescriptions := make([]string, len(changes))
	for i, change := range changes {
		// This is a simple representation; you can expand it based on the structure of your EventHistory
		changeDescriptions[i] = change.changeDescription()
	}
	combinedDescription := strings.Join(changeDescriptions, "; ")
	query := fmt.Sprintf("Adjust the planning based on the recent calendar changes: %s", combinedDescription)
	return query, nil
}

func (uc *OpenAIUseCase) buildOpenAIEventsQuery(events []*Event) (string, error) {
	if len(events) == 0 {
		return "", fmt.Errorf("no events provided")
	}

	// Convert each change into a string representation
	eventDescriptions := make([]string, len(events))
	for i, event := range events {
		// This is a simple representation; you can expand it based on the structure of your EventHistory
		eventDescriptions[i] = event.String()
	}
	combinedDescription := strings.Join(eventDescriptions, "; ")
	query := fmt.Sprintf("Adjust the planning based on this recent calendar events: %s", combinedDescription)
	return query, nil
}

func (uc *OpenAIUseCase) openAISystemQuery() string {
	return fmt.Sprint("You are my planing assistant. Your job is to help me plan my day. I will give you a list of events and changes to my calendar and you will help me plan my day.")
}

func (uc *OpenAIUseCase) GenerateCalendarEvents(ctx context.Context, calendar *Calendar, events []*Event) error {
	uc.log.Debugf("generate calendar events for calendar %s", calendar.ID)
	// Build the query
	messageContext := make([]openai.ChatCompletionMessage, 0)
	messageContext = append(messageContext, openai.ChatCompletionMessage{
		Role:    "system",
		Content: uc.openAISystemQuery(),
	})
	eventsQuery, err := uc.buildOpenAIEventsQuery(events)
	if err != nil {
		return err
	}
	messageContext = append(messageContext, openai.ChatCompletionMessage{
		Role:    "user",
		Content: eventsQuery,
	})

	uc.fr.Register(currentTimeFunctionDescription().Name, currentTimeFunctionDescription(), uc.currentTimeFunction)
	uc.fr.Register(createEventFunctionDescription().Name, createEventFunctionDescription(), uc.createEventFunction)

	request := &openai.ChatCompletionRequest{
		Messages:  messageContext,
		Functions: uc.fr.Descriptions(),
	}
	for {
		response, err := uc.client.DoRequest(ctx, request)
		if err != nil {
			return err
		}
		if response.Choices[0].FinishReason == "stop" {
			uc.log.Debugf("generate calendar events for calendar %s: %s", calendar.ID, response.Choices[0].Message.Content)
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
	return nil
}

func (uc *OpenAIUseCase) deleteEventFunction(ctx context.Context, arguments string) string {
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

func (uc *OpenAIUseCase) currentTimeFunction(_ context.Context, arguments string) string {
	uc.log.Debugf("currentTimeFunction: %s", arguments)
	return time.Now().Format(time.RFC3339)
}

func (uc *OpenAIUseCase) listEventsFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("listEventsFunction: %s", arguments)
	args := &struct {
		CalendarID string `json:"calendar_id"`
		StartTime  string `json:"start_time,omitempty"`
		EndTime    string `json:"end_time,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}

	token := GetToken(ctx)
	events, err := uc.gr.ListCalendarEvents(ctx, token, args.CalendarID, &GoogleListEventsOption{
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
func (uc *OpenAIUseCase) updateEventFunction(ctx context.Context, arguments string) string {
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
func (uc *OpenAIUseCase) createEventFunction(ctx context.Context, arguments string) string {
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
