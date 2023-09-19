package biz

import (
	"github.com/kdimtricp/aical/pkg/openai"
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
				"google_id": map[string]interface{}{
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
			"required": []string{"google_id", "start_time", "end_time"},
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
