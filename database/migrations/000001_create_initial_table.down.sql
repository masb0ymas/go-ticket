DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS schedules;
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS ticket_types;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS transaction_details;

DROP EXTENSION IF EXISTS "uuid-ossp";

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_events_location;
DROP INDEX IF EXISTS idx_events_schedule;
DROP INDEX IF EXISTS idx_ticket_types_event;
DROP INDEX IF EXISTS idx_transaction_details_transaction;
DROP INDEX IF EXISTS idx_transaction_details_ticket_type;
DROP INDEX IF EXISTS idx_transaction_details_transaction;
DROP INDEX IF EXISTS idx_transaction_details_ticket_type;