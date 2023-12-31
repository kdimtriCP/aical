package biz

import (
	"context"
	"encoding/json"
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
				"google_calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The Google ID of the calendar for the event creation.",
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
			"required": []string{"google_calendar_id", "title", "start_time", "end_time"},
		},
	}
}

// UpdateEventFunctionDescription is a function that returns description of a function that updates an event
func updateEventFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "update_event",
		Description: "Updates an event in the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"google_calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the Google calendar for the event update.",
				},
				"google_event_id": map[string]interface{}{
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
			"required": []string{"google_calendar_id", "google_event_id"},
		},
	}
}

// deleteEventFunctionDescription is a function that returns description of a function that deletes an event
func deleteEventFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "delete_event",
		Description: "Deletes an event from the calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"google_calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The ID of the Google calendar for the event deletion.",
				},
				"google_event_id": map[string]interface{}{
					"type":        "string",
					"description": "The Google ID of the event.",
				},
			},
			"required": []string{"google_calendar_id", "google_event_id"},
		},
	}
}

func currentTimeFunction(_ context.Context, _ string) string {
	return time.Now().Format(time.RFC3339)
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

// listEventsFunctionDescription is a function that returns description of a function that lists events
func listEventsFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "list_events",
		Description: "Lists events in the google calendar",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"google_calendar_id": map[string]interface{}{
					"type":        "string",
					"description": "The Google provided ID of the calendar where the events should be listed.",
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
			"required": []string{"google_calendar_id", "start_time", "end_time"},
		},
	}
}

// listUserCalendarsFunctionDescription is a function that returns description of a function that lists user calendars
func listUserCalendarsFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "list_user_calendars",
		Description: "Lists user calendars",
		Parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}
}

func adjustDateFunction(_ context.Context, arguments string) string {
	args := &struct {
		Date string `json:"date"`
		Days int    `json:"days"`
	}{}

	err := json.Unmarshal([]byte(arguments), args)
	if err != nil {
		return err.Error()
	}

	currentDate, err := time.Parse(time.RFC3339, args.Date)
	if err != nil {
		return err.Error()
	}

	newDate := currentDate.AddDate(0, 0, args.Days)
	return newDate.Format(time.RFC3339)
}

// adjustDateFunctionDescription is a function that returns description of a function that adjusts date
func adjustDateFunctionDescription() openai.FunctionDescription {
	return openai.FunctionDescription{
		Name:        "adjust_date",
		Description: "Adjusts date by adding or subtracting days",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"date": map[string]interface{}{
					"type":        "string",
					"description": "The date in RFC3339 format.",
				},
				"days": map[string]interface{}{
					"type":        "integer",
					"description": "The number of days to add or subtract.",
				},
			},
			"required": []string{"date", "days"},
		},
	}
}
