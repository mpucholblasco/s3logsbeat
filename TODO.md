TO DO
=====

Improvements
------------
[x] Repeated events
[ ] Keep tracking of ACKnowledged events to delete message once all events have been ACKed (present on event.Private)
[ ] Use `*beat.Event` instead of `beat.Event` to avoid copying objects.
[ ] Do not pull from SQS if output is down.
[ ] Add information to monitor
