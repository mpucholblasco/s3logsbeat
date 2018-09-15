TO DO
=====

Improvements
------------
[?] Repeated events (what happen with already present events? are they acked)
[x] Keep tracking of ACKnowledged events to delete message once all events have been ACKed (present on event.Private)
[x] Use `*beat.Event` instead of `beat.Event` to avoid copying objects.
[KO] Do not pull from SQS if output is down -> couldn't find a way to do it because of library.
[x] Add information to monitor
[x] Add failures to monitor
[x] Add key fields to events
[x] Test Once
[ ] Change fields (name and similar) based on config
[x] Flag to keep messages on queue (used on tests)
[x] Test CloudFront logs
[x] Include log type on @metadata (in format field)
[ ] make test & others

Features
[ ] Only one AWS credentials used and based on roles or .aws/config
