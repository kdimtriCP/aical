package biz

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/pkg/openai"
	"time"
)

// createEventFunctionDescription is a function that returns description of a function that creates an event
func createEventFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "create_event",
		Description: "Creates an event in the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique ID of the calendar where the event should be created.",
				},
				"google_id": map[string]interface{}{
					"type":        "string",
					"description": "The Google ID of the event. Optional parameter.",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "The summary or title of the event.",
				},
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The location of the event. Optional parameter.",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "The start time of the event in RFC3339 format.",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "The end time of the event in RFC3339 format.",
				},
			},
			"required": []string{"calendar_id", "title", "start_time", "end_time"},
		},
	}
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

// UpdateEventFunctionDescription is a function that returns description of a function that updates an event
func updateEventFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "update_event",
		Description: "Updates an event in the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique ID of the calendar where the event should be updated.",
				},
				"google_id": map[string]interface{}{
					"type":        "string",
					"description": "The Google ID of the event.",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "The summary or title of the event. Optional parameter.",
				},
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The location of the event. Optional parameter.",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "The start time of the event in RFC3339 format. Optional parameter.",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "The end time of the event in RFC3339 format. Optional parameter.",
				},
			},
			"required": []string{"calendar_id", "google_id"},
		},
	}
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

// deleteEventFunctionDescription is a function that returns description of a function that deletes an event
func deleteEventFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "delete_event",
		Description: "Deletes an event from the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique ID of the calendar where the event should be deleted.",
				},
				"google_id": map[string]interface{}{
					"type":        "string",
					"description": "The Google ID of the event.",
				},
			},
			"required": []string{"calendar_id", "google_id"},
		},
	}
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

// currentTimeFunctionDescription is a function that returns description of a function that returns current time
func currentTimeFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "current_time",
		Description: "Returns current time in RFC3339 format",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}
}
func (uc *OpenAIUseCase) currentTimeFunction(ctx context.Context, arguments string) string {
	uc.log.Debugf("currentTimeFunction: %s", arguments)
	return time.Now().Format(time.RFC3339)
}
