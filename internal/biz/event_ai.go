package biz

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kdimtricp/aical/pkg/openai"
	"time"
)

// CreateEventFunctionDescription is a function that returns description of a function that creates an event
func CreateEventFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
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
				"is_all_day": map[string]interface{}{
					"type":        "boolean",
					"description": "Indicates if the event is an all-day event. Optional parameter.",
				},
			},
			"required": []string{"calendar_id", "title", "start_time", "end_time"},
		},
	}
}
func CreateEventFunction(arguments string) (*Event, error) {
	args := &struct {
		CalendarID string    `json:"calendar_id"`
		GoogleID   string    `json:"google_id,omitempty"`
		Title      string    `json:"title"`
		Location   string    `json:"location,omitempty"`
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
		IsAllDay   bool      `json:"is_all_day,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return nil, err
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
		IsAllDay:   args.IsAllDay,
	}

	// Далее, вы можете добавить код для создания события в вашей системе
	// и вернуть созданное событие.

	return event, nil
}

// UpdateEventFunctionDescription is a function that returns description of a function that updates an event
func UpdateEventFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
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
				"is_all_day": map[string]interface{}{
					"type":        "boolean",
					"description": "Indicates if the event is an all-day event. Optional parameter.",
				},
			},
			"required": []string{"calendar_id", "google_id"},
		},
	}
}
func UpdateEventFunction(arguments string) (*Event, error) {
	args := &struct {
		CalendarID string    `json:"calendar_id"`
		GoogleID   string    `json:"google_id"`
		Title      string    `json:"title,omitempty"`
		Location   string    `json:"location,omitempty"`
		StartTime  time.Time `json:"start_time,omitempty"`
		EndTime    time.Time `json:"end_time,omitempty"`
		IsAllDay   bool      `json:"is_all_day,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return nil, err
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
		IsAllDay:   args.IsAllDay,
	}

	// Далее, вы можете добавить код для обновления события в вашей системе
	// и вернуть обновленное событие.

	return event, nil
}

// DeleteEventFunctionDescription is a function that returns description of a function that deletes an event
func DeleteEventFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
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
func DeleteEventFunction(arguments string) error {
	args := &struct {
		CalendarID string `json:"calendar_id"`
		GoogleID   string `json:"google_id"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err
	}

	// Далее, вы можете добавить код для удаления события в вашей системе

	return nil
}

// GetEventFunctionDescription is a function that returns description of a function that gets an event
func GetEventFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
		Name:        "get_event",
		Description: "Gets an event from the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique ID of the calendar where the event should be retrieved.",
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
func GetEventFunction(arguments string) (*Event, error) {
	args := &struct {
		CalendarID string `json:"calendar_id"`
		GoogleID   string `json:"google_id"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return nil, err
	}

	// Далее, вы можете добавить код для получения события из вашей системы
	// и вернуть полученное событие.

	return nil, nil
}

// ListEventsFunctionDescription is a function that returns description of a function that lists events
func ListEventsFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
		Name:        "list_events",
		Description: "Lists events from the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The unique ID of the calendar where the events should be retrieved.",
				},
				"start_time": map[string]interface{}{
					"type":        "string",
					"description": "The start time of the events in RFC3339 format. Optional parameter.",
				},
				"end_time": map[string]interface{}{
					"type":        "string",
					"description": "The end time of the events in RFC3339 format. Optional parameter.",
				},
			},
			"required": []string{"calendar_id"},
		},
	}
}
func ListEventsFunction(arguments string) ([]*Event, error) {
	args := &struct {
		CalendarID string    `json:"calendar_id"`
		StartTime  time.Time `json:"start_time,omitempty"`
		EndTime    time.Time `json:"end_time,omitempty"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return nil, err
	}

	// Далее, вы можете добавить код для получения списка событий из вашей системы
	// и вернуть полученный список событий.

	return nil, nil
}

// ListCalendarsFunctionDescription is a function that returns description of a function that lists calendars
func ListCalendarsFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
		Name:        "list_calendars",
		Description: "Lists calendars",
		Parameters:  map[string]interface{}{},
	}
}
func ListCalendarsFunction(arguments string) ([]*Calendar, error) {
	// Далее, вы можете добавить код для получения списка календарей из вашей системы
	// и вернуть полученный список календарей.

	return nil, nil
}

// CurrentTimeFunctionDescription is a function that returns description of a function that returns current time
func CurrentTimeFunctionDescription() openai.ChatCompletionFunction {
	return openai.ChatCompletionFunction{
		Name:        "current_time",
		Description: "Returns current time in RFC3339 format",
		Parameters:  map[string]interface{}{},
	}
}
func CurrentTimeFunction(arguments string) string {
	return time.Now().Format(time.RFC3339)
}
