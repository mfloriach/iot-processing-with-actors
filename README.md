# TODO
Priority — alarmas antes que telemetría.
Fairness — un device con 10.000 tareas no monopoliza el worker.
Work stealing — workers idle pueden robar trabajo.
Backpressure — limitar tareas pendientes.
Per-device ordering — mantener orden dentro de cada device.
Coalescing — si llegan 100 actualizaciones de temperatura, procesar solo la última.
Rate limiting — limitar cuánto trabajo genera un device.
Deadlines — tareas con deadline tienen prioridad.

# Implemented
* Actor mailbox → autoridad para modificar el estado.
* Snapshot/atomic → lectura rápida del último estado.
* No mutex
* Group devices in shards
* Worker per shard, no per actor
* SSE for event publish
