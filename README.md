 - Password vault holds all mailserver password encrypted with an admin password
 - Needs 2fa on all important commands (not on innocent read commands)
   - If recieved 2fa on last 5 minutes, not ask again
 - client - server model
   - client is cli
   - gRPC comms to server
   - server can interact with docker daemon, database, rspamd controller, etc
