


# Implemented
* Actor mailbox → autoridad para modificar el estado.
* Snapshot/atomic → lectura rápida del último estado.
* No go routine for actor
* Work stealing — workers idle pueden robar trabajo.
* SSE for event publish
* Fairness — un device con 10.000 tareas no monopoliza el worker.
* quantum - You don't necessarily need exactly one task per turn. Use a quantum
* Backpressure — limitar tareas pendientes.