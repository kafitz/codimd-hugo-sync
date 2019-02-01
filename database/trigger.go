package database

import "database/sql"

// NotifyNoteEvent is the query for generating the new post INSERT or UPDATE
// event to notify this program of blog changes
// http://coussej.github.io/2015/09/15/Listening-to-generic-JSON-notifications-from-PostgreSQL-in-Go/
var NotifyNoteEvent = `
	CREATE OR REPLACE FUNCTION notify_note_event() RETURNS TRIGGER AS $$
		DECLARE
			data json;
			notification json;
		BEGIN
        	-- Convert the old or new row to JSON, based on the kind of action.
        	-- Action = DELETE?             -> OLD row
			-- Action = INSERT or UPDATE?   -> NEW row
			IF (TG_OP = 'DELETE') THEN
				data = row_to_json(OLD);
			ELSE
				data = row_to_json(NEW);
			END IF;

			notification = json_build_object(
				'table', TG_TABLE_NAME,
				'action', TG_OP,
				'data', data
			);

			PERFORM pg_notify('events', notification::text);

			RETURN NULL;
		END;
	$$ LANGUAGE plpgsql
`

// NotifyNoteTrigger is the query for listening to the Posts table and raising the event
var NotifyNoteTrigger = `
	CREATE TRIGGER posts_notify_event
	AFTER INSERT OR UPDATE OR DELETE ON "Notes"
		FOR EACH ROW EXECUTE PROCEDURE notify_note_event();
`

// AddTrigger adds the NotifyNoteEvent and NotifyNoteTrigger to the HackMD database
func AddTrigger(db *sql.DB) {
	db.Exec(NotifyNoteEvent)
	db.Exec(NotifyNoteTrigger)
}
