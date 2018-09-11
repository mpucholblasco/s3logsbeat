TO DO
=====

Improvements
------------
[?] Repeated events (what happen with already present events? are they acked)
[x] Keep tracking of ACKnowledged events to delete message once all events have been ACKed (present on event.Private)
[ ] Use `*beat.Event` instead of `beat.Event` to avoid copying objects.
[KO] Do not pull from SQS if output is down -> couldn't find a way to do it because of library.
[x] Add information to monitor
[x] Add failures to monitor
[ ] Test Once
[ ] Add key fields to events
[ ] Change fields (name and similar) based on config
